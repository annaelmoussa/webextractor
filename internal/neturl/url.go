package neturl

import (
	"errors"
	"strings"
)

// URL represents a parsed URL
type URL struct {
	Scheme   string
	Host     string
	Path     string
	RawQuery string
	Fragment string
}

// Parse parses a raw URL string and returns a URL structure
func Parse(rawurl string) (*URL, error) {
	if rawurl == "" {
		return nil, errors.New("empty url")
	}
	
	u := &URL{}
	
	// Remove fragment first
	if idx := strings.Index(rawurl, "#"); idx >= 0 {
		u.Fragment = rawurl[idx+1:]
		rawurl = rawurl[:idx]
	}
	
	// Extract query string
	if idx := strings.Index(rawurl, "?"); idx >= 0 {
		u.RawQuery = rawurl[idx+1:]
		rawurl = rawurl[:idx]
	}
	
	// Extract scheme
	if idx := strings.Index(rawurl, "://"); idx >= 0 {
		u.Scheme = strings.ToLower(rawurl[:idx])
		rawurl = rawurl[idx+3:]
	} else {
		return nil, errors.New("missing protocol scheme")
	}
	
	// Split host and path
	if idx := strings.Index(rawurl, "/"); idx >= 0 {
		u.Host = rawurl[:idx]
		u.Path = rawurl[idx:]
	} else {
		u.Host = rawurl
		u.Path = "/"
	}
	
	// Validate scheme
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("unsupported protocol scheme")
	}
	
	// Validate host
	if u.Host == "" {
		return nil, errors.New("empty host")
	}
	
	return u, nil
}

// String reassembles the URL into a string
func (u *URL) String() string {
	var result strings.Builder
	
	if u.Scheme != "" {
		result.WriteString(u.Scheme)
		result.WriteString("://")
	}
	
	result.WriteString(u.Host)
	
	if u.Path == "" {
		result.WriteString("/")
	} else {
		result.WriteString(u.Path)
	}
	
	if u.RawQuery != "" {
		result.WriteString("?")
		result.WriteString(u.RawQuery)
	}
	
	if u.Fragment != "" {
		result.WriteString("#")
		result.WriteString(u.Fragment)
	}
	
	return result.String()
}

// ResolveReference resolves a URI reference to an absolute URI from an absolute base URI
func (u *URL) ResolveReference(ref *URL) *URL {
	if ref == nil {
		return u
	}
	
	// If ref has a scheme, it's absolute
	if ref.Scheme != "" {
		return ref
	}
	
	result := &URL{
		Scheme: u.Scheme,
		Host:   u.Host,
	}
	
	// Handle absolute path
	if strings.HasPrefix(ref.Path, "/") {
		result.Path = ref.Path
		result.RawQuery = ref.RawQuery
		result.Fragment = ref.Fragment
		return result
	}
	
	// Handle relative path
	if ref.Path == "" {
		result.Path = u.Path
		if ref.RawQuery != "" {
			result.RawQuery = ref.RawQuery
		} else {
			result.RawQuery = u.RawQuery
		}
		result.Fragment = ref.Fragment
		return result
	}
	
	// Resolve relative path
	basePath := u.Path
	if !strings.HasSuffix(basePath, "/") {
		// Remove filename from base path
		if idx := strings.LastIndex(basePath, "/"); idx >= 0 {
			basePath = basePath[:idx+1]
		} else {
			basePath = "/"
		}
	}
	
	result.Path = basePath + ref.Path
	result.RawQuery = ref.RawQuery
	result.Fragment = ref.Fragment
	
	// Clean up path (remove . and ..)
	result.Path = cleanPath(result.Path)
	
	return result
}

// cleanPath cleans up relative path components
func cleanPath(path string) string {
	if path == "" {
		return "/"
	}
	
	parts := strings.Split(path, "/")
	var clean []string
	
	for _, part := range parts {
		switch part {
		case "", ".":
			// Skip empty and current directory
			continue
		case "..":
			// Parent directory
			if len(clean) > 0 {
				clean = clean[:len(clean)-1]
			}
		default:
			clean = append(clean, part)
		}
	}
	
	result := "/" + strings.Join(clean, "/")
	if strings.HasSuffix(path, "/") && !strings.HasSuffix(result, "/") {
		result += "/"
	}
	
	return result
}