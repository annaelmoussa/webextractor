package io

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Result représente un résultat d'extraction.
type Result struct {
	Selector string   `json:"selector"`
	Matches  []string `json:"matches"`
}

// DocumentResult est la structure de niveau supérieur du format JSON.
type DocumentResult struct {
	URL     string   `json:"url"`
	Results []Result `json:"results"`
}

// StructuredResult représente le format de sortie structuré.
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

// validateOutputPath valide le chemin de sortie.
func validateOutputPath(path string) error {
	if path == "-" || path == "" {
		return nil
	}

	cleanPath := filepath.Clean(path)

	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid path: directory traversal not allowed")
	}

	if strings.HasPrefix(cleanPath, "/etc/") ||
		strings.HasPrefix(cleanPath, "/usr/") ||
		(strings.HasPrefix(cleanPath, "/var/") && !strings.HasPrefix(cleanPath, "/var/folders/")) {
		return fmt.Errorf("invalid path: cannot write to system directories")
	}

	return nil
}

// Write écrit le résultat dans le chemin de fichier donné ("-" signifie stdout).
// path est le chemin de sortie, doc est le résultat à écrire.
func Write(path string, doc DocumentResult) error {
	if err := validateOutputPath(path); err != nil {
		return fmt.Errorf("output path validation failed: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	if path != "-" && path != "" {
		file, err := os.Create(path) // #nosec G304 - path validated above - on valide le chemin ci-dessus
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

// WriteStructured écrit le résultat structuré dans le chemin de fichier donné ("-" signifie stdout).
// path est le chemin de sortie, doc est le résultat à écrire.
func WriteStructured(path string, doc StructuredResult) error {
	if err := validateOutputPath(path); err != nil {
		return fmt.Errorf("output path validation failed: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	if path != "-" && path != "" {
		file, err := os.Create(path) // #nosec G304 - path validated above - on valide le chemin ci-dessus
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
