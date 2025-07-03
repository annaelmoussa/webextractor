package io

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Result represents one extraction result.
type Result struct {
	Selector string   `json:"selector"`
	Matches  []string `json:"matches"`
}

// DocumentResult is the top-level structure of the JSON output.
type DocumentResult struct {
	URL     string   `json:"url"`
	Results []Result `json:"results"`
}

// StructuredResult represents the structured output format requested by the user
type StructuredResult struct {
	URL        string   `json:"url"`
	Title      string   `json:"title,omitempty"`
	H1         []string `json:"h1,omitempty"`
	H2         []string `json:"h2,omitempty"`
	H3         []string `json:"h3,omitempty"`
	Paragraphs []string `json:"paragraphs,omitempty"`
	Links      []string `json:"links,omitempty"`
	Images     []string `json:"images,omitempty"`
	Lists      []string `json:"lists,omitempty"`
}

// validateOutputPath checks if the output path is safe to write to
func validateOutputPath(path string) error {
	if path == "-" || path == "" {
		return nil
	}

	// Clean the path to prevent directory traversal
	cleanPath := filepath.Clean(path)

	// Check for directory traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid path: directory traversal not allowed")
	}

	// Ensure we're not trying to write to system directories
	if strings.HasPrefix(cleanPath, "/etc/") ||
		strings.HasPrefix(cleanPath, "/usr/") ||
		strings.HasPrefix(cleanPath, "/var/") {
		return fmt.Errorf("invalid path: cannot write to system directories")
	}

	return nil
}

// Write writes the result to the given file path ("-" means stdout).
func Write(path string, doc DocumentResult) error {
	if err := validateOutputPath(path); err != nil {
		return fmt.Errorf("output path validation failed: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	if path != "-" && path != "" {
		file, err := os.Create(path) // #nosec G304 - path validated above
		if err != nil {
			return err
		}
		defer file.Close()
		enc = json.NewEncoder(file)
	}
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	return nil
}

// WriteStructured writes the structured result to the given file path ("-" means stdout).
func WriteStructured(path string, doc StructuredResult) error {
	if err := validateOutputPath(path); err != nil {
		return fmt.Errorf("output path validation failed: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	if path != "-" && path != "" {
		file, err := os.Create(path) // #nosec G304 - path validated above
		if err != nil {
			return err
		}
		defer file.Close()
		enc = json.NewEncoder(file)
	}
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	return nil
}
