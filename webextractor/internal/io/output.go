package io

import (
	"encoding/json"
	"fmt"
	"os"
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

// Write writes the result to the given file path ("-" means stdout).
func Write(path string, doc DocumentResult) error {
	enc := json.NewEncoder(os.Stdout)
	if path != "-" && path != "" {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		enc = json.NewEncoder(file)
	}
	enc.SetIndent("", "  ") // 2 spaces indentation
	if err := enc.Encode(doc); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	return nil
}
