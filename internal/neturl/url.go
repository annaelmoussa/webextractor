package neturl

import (
	"errors"
	"strings"
)

// URL représente une URL analysée
type URL struct {
	Scheme   string
	Host     string
	Path     string
	RawQuery string
	Fragment string
}

// Parse analyse une chaîne d'URL brute et retourne une structure URL
func Parse(rawurl string) (*URL, error) {
	if rawurl == "" {
		return nil, errors.New("empty url")
	}
	
	u := &URL{}
	
	// On retire le fragment en premier
	if idx := strings.Index(rawurl, "#"); idx >= 0 {
		u.Fragment = rawurl[idx+1:]
		rawurl = rawurl[:idx]
	}
	
	// On extrait la chaîne de requête
	if idx := strings.Index(rawurl, "?"); idx >= 0 {
		u.RawQuery = rawurl[idx+1:]
		rawurl = rawurl[:idx]
	}
	
	// On extrait le schéma
	if idx := strings.Index(rawurl, "://"); idx >= 0 {
		u.Scheme = strings.ToLower(rawurl[:idx])
		rawurl = rawurl[idx+3:]
	} else {
		return nil, errors.New("missing protocol scheme")
	}
	
	// On sépare le host et le chemin
	if idx := strings.Index(rawurl, "/"); idx >= 0 {
		u.Host = rawurl[:idx]
		u.Path = rawurl[idx:]
	} else {
		u.Host = rawurl
		u.Path = "/"
	}
	
	// On valide le schéma
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("unsupported protocol scheme")
	}
	
	// On valide le host
	if u.Host == "" {
		return nil, errors.New("empty host")
	}
	
	return u, nil
}

// String reassemble l'URL en chaîne de caractères.
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

// ResolveReference résout une référence URI en URI absolue depuis une URI de base absolue.
func (u *URL) ResolveReference(ref *URL) *URL {
	if ref == nil {
		return u
	}
	
	// Si ref a un schéma, c'est absolu
	if ref.Scheme != "" {
		return ref
	}
	
	result := &URL{
		Scheme: u.Scheme,
		Host:   u.Host,
	}
	
	// On gère le chemin absolu
	if strings.HasPrefix(ref.Path, "/") {
		result.Path = ref.Path
		result.RawQuery = ref.RawQuery
		result.Fragment = ref.Fragment
		return result
	}
	
	// On gère le chemin relatif
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
	
	// On résout le chemin relatif
	basePath := u.Path
	if !strings.HasSuffix(basePath, "/") {
		// Supprime le nom de fichier du chemin de base
		if idx := strings.LastIndex(basePath, "/"); idx >= 0 {
			basePath = basePath[:idx+1]
		} else {
			basePath = "/"
		}
	}
	
	result.Path = basePath + ref.Path
	result.RawQuery = ref.RawQuery
	result.Fragment = ref.Fragment
	
	// On nettoie le chemin (on retire . et ..)
	result.Path = cleanPath(result.Path)
	
	return result
}

// cleanPath nettoie les composants de chemin relatifs
func cleanPath(path string) string {
	if path == "" {
		return "/"
	}
	
	parts := strings.Split(path, "/")
	var clean []string
	
	for _, part := range parts {
		switch part {
		case "", ".":
			// On ignore le chemin vide et le répertoire courant
			continue
		case "..":
			// On supprime le répertoire parent
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