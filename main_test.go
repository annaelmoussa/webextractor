package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	ioLib "webextractor/internal/io"
	"webextractor/internal/neturl"
)

func TestCLI(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><body><div>Hello</div></body></html>`))
	}))
	defer srv.Close()

	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	origArgs := os.Args
	os.Args = []string{"cmd", "-url", srv.URL, "-sel", "div"}

	main()

	w.Close()
	os.Stdout = origStdout
	os.Args = origArgs

	outBytes, _ := io.ReadAll(r)
	out := string(outBytes)
	if !strings.Contains(out, "Hello") {
		t.Fatalf("output not correct: %s", out)
	}
}

func TestURLParsing(t *testing.T) {
	_, err := parseURL("://invalid-url")
	if err == nil {
		t.Errorf("Expected error for invalid URL")
	}

	validURL, err := parseURL("https://example.com")
	if err != nil {
		t.Fatalf("Valid URL should not error: %v", err)
	}
	if validURL.String() != "https://example.com/" {
		t.Errorf("URL parsing failed, got: %s", validURL.String())
	}
}

func parseURL(urlStr string) (*neturl.URL, error) {
	return neturl.Parse(urlStr)
}

func TestSelectorParsing(t *testing.T) {
	selectors := strings.Split("div,.content,p", ",")

	if len(selectors) != 3 {
		t.Fatalf("Expected 3 selectors, got %d", len(selectors))
	}

	expected := []string{"div", ".content", "p"}
	for i, sel := range selectors {
		if strings.TrimSpace(sel) != expected[i] {
			t.Errorf("Expected selector %d to be '%s', got '%s'", i, expected[i], sel)
		}
	}
}

func TestEmptySelectorHandling(t *testing.T) {
	selectors := strings.Split("div, , p", ",")
	var cleaned []string

	for _, sel := range selectors {
		sel = strings.TrimSpace(sel)
		if sel != "" {
			cleaned = append(cleaned, sel)
		}
	}

	if len(cleaned) != 2 {
		t.Fatalf("Expected 2 non-empty selectors, got %d", len(cleaned))
	}
}

func TestStructuredDataProcessing(t *testing.T) {
	structuredData := map[string]interface{}{
		"title":      "Test Title",
		"h1":         []string{"H1 Title"},
		"h2":         []string{"H2 Title 1", "H2 Title 2"},
		"paragraphs": []string{"Para 1", "Para 2"},
		"links":      []string{"https://link1.com", "https://link2.com"},
		"images":     []string{"/image1.jpg"},
		"lists":      []string{"List item 1 | List item 2"},
	}

	result := ioLib.StructuredResult{URL: "https://test.com"}

	if title, ok := structuredData["title"].(string); ok {
		result.Title = title
	}
	if h1List, ok := structuredData["h1"].([]string); ok {
		result.H1 = h1List
	}
	if h2List, ok := structuredData["h2"].([]string); ok {
		result.H2 = h2List
	}
	if paragraphs, ok := structuredData["paragraphs"].([]string); ok {
		result.Paragraphs = paragraphs
	}
	if links, ok := structuredData["links"].([]string); ok {
		result.Links = links
	}
	if images, ok := structuredData["images"].([]string); ok {
		result.Images = images
	}
	if lists, ok := structuredData["lists"].([]string); ok {
		result.Lists = lists
	}

	if result.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", result.Title)
	}
	if len(result.H1) != 1 || result.H1[0] != "H1 Title" {
		t.Errorf("H1 processing failed")
	}
	if len(result.H2) != 2 {
		t.Errorf("Expected 2 H2 elements, got %d", len(result.H2))
	}
	if len(result.Paragraphs) != 2 {
		t.Errorf("Expected 2 paragraphs, got %d", len(result.Paragraphs))
	}
	if len(result.Links) != 2 {
		t.Errorf("Expected 2 links, got %d", len(result.Links))
	}
}

func TestSelectorDeduplication(t *testing.T) {
	collectedSelectors := []string{"div", "p", "div", ".class", "p", "span"}

	uniqueSelectors := make(map[string]struct{})
	for _, s := range collectedSelectors {
		uniqueSelectors[s] = struct{}{}
	}

	finalSelectors := make([]string, 0, len(uniqueSelectors))
	for s := range uniqueSelectors {
		finalSelectors = append(finalSelectors, s)
	}

	if len(finalSelectors) != 4 {
		t.Fatalf("Expected 4 unique selectors, got %d", len(finalSelectors))
	}
}
