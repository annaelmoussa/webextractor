package tui

import (
	"net/url"
	"strings"
	"testing"
	"webextractor/internal/parser"

	"golang.org/x/net/html"
)

func TestExtractPageInfo(t *testing.T) {
	htmlStr := `
		<html>
			<head><title>Test Page</title></head>
			<body>
				<h1>Main Title</h1>
				<h2>Subtitle</h2>
				<p>First paragraph with content.</p>
				<p>Second paragraph.</p>
				<a href="/link1">Link 1</a>
				<a href="https://example.com">External Link</a>
				<img src="/image1.jpg" alt="Test Image" />
				<ul>
					<li>Item 1</li>
					<li>Item 2</li>
				</ul>
			</body>
		</html>
	`

	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	testURL, _ := url.Parse("https://test.com")
	info := extractPageInfo(doc, testURL)

	// Test URL
	if info.URL != "https://test.com" {
		t.Errorf("Expected URL 'https://test.com', got '%s'", info.URL)
	}

	// Test title
	if info.Title != "Test Page" {
		t.Errorf("Expected title 'Test Page', got '%s'", info.Title)
	}

	// Test H1
	if len(info.H1) != 1 || info.H1[0] != "Main Title" {
		t.Errorf("Expected H1 ['Main Title'], got %v", info.H1)
	}

	// Test H2
	if len(info.H2) != 1 || info.H2[0] != "Subtitle" {
		t.Errorf("Expected H2 ['Subtitle'], got %v", info.H2)
	}

	// Test paragraphs
	if len(info.Paragraphs) != 2 {
		t.Errorf("Expected 2 paragraphs, got %d", len(info.Paragraphs))
	}

	// Test links (should be made absolute)
	if len(info.Links) != 2 {
		t.Errorf("Expected 2 links, got %d", len(info.Links))
	}

	// Test images
	if len(info.Images) != 1 {
		t.Errorf("Expected 1 image, got %d", len(info.Images))
	}
	if info.Images[0].Alt != "Test Image" {
		t.Errorf("Expected image alt 'Test Image', got '%s'", info.Images[0].Alt)
	}

	// Test lists
	if len(info.Lists) != 1 {
		t.Errorf("Expected 1 list, got %d", len(info.Lists))
	}
}

func TestHandleSelectAll(t *testing.T) {
	info := PageInfo{
		URL:        "https://test.com",
		Title:      "Test Title",
		H1:         []string{"H1 Title"},
		Paragraphs: []string{"Test paragraph"},
		Links:      []parser.Link{{Href: "https://example.com", Text: "Example"}},
	}

	result := handleSelectAll(info)

	if !result.Finished {
		t.Error("Expected result to be finished")
	}

	if len(result.Selectors) == 0 {
		t.Error("Expected selectors to be populated")
	}

	if result.SelectedData == nil {
		t.Error("Expected SelectedData to be populated")
	}

	// Check that title is in selected data
	if title, ok := result.SelectedData["title"].(string); !ok || title != "Test Title" {
		t.Errorf("Expected title 'Test Title' in SelectedData, got %v", result.SelectedData["title"])
	}
}

func TestHandleCategorySelection(t *testing.T) {
	info := PageInfo{
		URL:        "https://test.com",
		Title:      "Test Title",
		H1:         []string{"H1 Title"},
		Paragraphs: []string{"Test paragraph"},
	}

	// Test title selection
	result, err := handleCategorySelection("title", info)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result to be non-nil")
	}

	if !result.Finished {
		t.Error("Expected result to be finished")
	}

	if title, ok := result.SelectedData["title"].(string); !ok || title != "Test Title" {
		t.Errorf("Expected title 'Test Title' in SelectedData, got %v", result.SelectedData["title"])
	}

	// Test invalid selection
	result, err = handleCategorySelection("invalid", info)
	if err == nil {
		t.Error("Expected error for invalid selection")
	}

	// Test h1 selection
	result, err = handleCategorySelection("h1", info)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if h1List, ok := result.SelectedData["h1"].([]string); !ok || len(h1List) != 1 || h1List[0] != "H1 Title" {
		t.Errorf("Expected H1 ['H1 Title'] in SelectedData, got %v", result.SelectedData["h1"])
	}
}

func TestHandleLinkNavigation(t *testing.T) {
	links := []parser.Link{
		{Href: "https://example.com", Text: "Example"},
		{Href: "https://test.com", Text: "Test"},
	}

	// Test valid link navigation
	result, err := handleLinkNavigation("L0", links)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.NextURL != "https://example.com" {
		t.Errorf("Expected NextURL 'https://example.com', got '%s'", result.NextURL)
	}

	// Test invalid format
	_, err = handleLinkNavigation("L", links)
	if err == nil {
		t.Error("Expected error for invalid format")
	}

	// Test out of bounds
	_, err = handleLinkNavigation("L5", links)
	if err == nil {
		t.Error("Expected error for out of bounds index")
	}
}

// Helper function to find an element by tag name
func findElement(n *html.Node, tag string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tag {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if res := findElement(c, tag); res != nil {
			return res
		}
	}
	return nil
}
