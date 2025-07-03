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

func TestBuildSelectableElements(t *testing.T) {
	info := PageInfo{
		URL:        "https://test.com",
		Title:      "Test Title",
		H1:         []string{"H1 Title", "Another H1"},
		H2:         []string{"H2 Title"},
		Paragraphs: []string{"First paragraph", "Second paragraph"},
		Links:      []parser.Link{{Href: "https://example.com", Text: "Example"}},
		Images:     []ImageInfo{{Src: "/test.jpg", Alt: "Test Image"}},
		Lists:      []string{"Item1 | Item2"},
	}

	elements := buildSelectableElements(info)

	// Should have: 1 title + 2 h1 + 1 h2 + 2 p + 1 link + 1 image + 1 list = 9 elements
	expectedCount := 9
	if len(elements) != expectedCount {
		t.Errorf("Expected %d elements, got %d", expectedCount, len(elements))
	}

	// Check first element (title)
	if elements[0].Type != "title" || elements[0].Data.(string) != "Test Title" {
		t.Errorf("Expected first element to be title 'Test Title', got type=%s data=%v", elements[0].Type, elements[0].Data)
	}

	// Check indices are correct
	for i, elem := range elements {
		if elem.Index != i {
			t.Errorf("Expected element %d to have index %d, got %d", i, i, elem.Index)
		}
	}

	// Check that different types are present
	types := make(map[string]int)
	for _, elem := range elements {
		types[elem.Type]++
	}

	expectedTypes := map[string]int{
		"title": 1,
		"h1":    2,
		"h2":    1,
		"p":     2,
		"link":  1,
		"image": 1,
		"list":  1,
	}

	for expectedType, expectedCount := range expectedTypes {
		if types[expectedType] != expectedCount {
			t.Errorf("Expected %d elements of type '%s', got %d", expectedCount, expectedType, types[expectedType])
		}
	}
}

func TestParseIndices(t *testing.T) {
	tests := []struct {
		input     string
		maxIndex  int
		expected  []int
		shouldErr bool
	}{
		// Valid single indices
		{"0", 5, []int{0}, false},
		{"3", 5, []int{3}, false},
		{"4", 5, []int{4}, false},

		// Valid multiple indices
		{"0,2,4", 5, []int{0, 2, 4}, false},
		{"1,3", 5, []int{1, 3}, false},

		// Valid ranges
		{"0-2", 5, []int{0, 1, 2}, false},
		{"1-3", 5, []int{1, 2, 3}, false},

		// Valid combinations
		{"0,2-4", 5, []int{0, 2, 3, 4}, false},
		{"0,3,1-2", 5, []int{0, 3, 1, 2}, false},

		// Invalid cases
		{"", 5, nil, true},      // Empty input
		{"5", 5, nil, true},     // Out of bounds
		{"-1", 5, nil, true},    // Negative
		{"0-5", 5, nil, true},   // Range out of bounds
		{"3-1", 5, nil, true},   // Invalid range
		{"a", 5, nil, true},     // Non-numeric
		{"0,a", 5, nil, true},   // Mixed valid/invalid
		{"0-", 5, nil, true},    // Incomplete range
		{"0-2-4", 5, nil, true}, // Invalid range format
	}

	for _, test := range tests {
		result, err := parseIndices(test.input, test.maxIndex)

		if test.shouldErr {
			if err == nil {
				t.Errorf("Expected error for input '%s', but got none", test.input)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}

		// Check that all expected indices are present (order doesn't matter due to deduplication)
		if len(result) != len(test.expected) {
			t.Errorf("For input '%s': expected %d indices, got %d", test.input, len(test.expected), len(result))
			continue
		}

		// Convert to sets for comparison
		resultSet := make(map[int]bool)
		for _, idx := range result {
			resultSet[idx] = true
		}

		expectedSet := make(map[int]bool)
		for _, idx := range test.expected {
			expectedSet[idx] = true
		}

		// Check all expected indices are present
		for expectedIdx := range expectedSet {
			if !resultSet[expectedIdx] {
				t.Errorf("For input '%s': expected index %d not found in result %v", test.input, expectedIdx, result)
			}
		}

		// Check no unexpected indices are present
		for resultIdx := range resultSet {
			if !expectedSet[resultIdx] {
				t.Errorf("For input '%s': unexpected index %d found in result %v", test.input, resultIdx, result)
			}
		}
	}
}

func TestHandleIndexSelection(t *testing.T) {
	// Create test state
	info := PageInfo{
		Title:      "Test Title",
		H1:         []string{"H1 Title"},
		Paragraphs: []string{"Paragraph 1", "Paragraph 2"},
	}
	elements := buildSelectableElements(info)
	state := SelectionState{
		Elements: elements,
		Selected: make([]bool, len(elements)),
		PageInfo: info,
	}

	// Test valid selection
	err := handleIndexSelection("0,2", &state)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that indices 0 and 2 are selected
	if !state.Selected[0] {
		t.Error("Expected index 0 to be selected")
	}
	if !state.Selected[2] {
		t.Error("Expected index 2 to be selected")
	}
	if state.Selected[1] {
		t.Error("Expected index 1 to NOT be selected")
	}

	// Test selecting already selected indices (should be idempotent)
	err = handleIndexSelection("0", &state)
	if err != nil {
		t.Fatalf("Unexpected error for already selected index: %v", err)
	}

	// Test invalid selection
	err = handleIndexSelection("99", &state)
	if err == nil {
		t.Error("Expected error for out of bounds index")
	}

	// Test invalid format
	err = handleIndexSelection("invalid", &state)
	if err == nil {
		t.Error("Expected error for invalid format")
	}
}

func TestHandleFinishWithSelections(t *testing.T) {
	// Create test state with some selections
	info := PageInfo{
		Title:      "Test Title",
		H1:         []string{"H1 Title"},
		H2:         []string{"H2 Title"},
		Paragraphs: []string{"Paragraph 1", "Paragraph 2"},
		Links:      []parser.Link{{Href: "https://example.com", Text: "Example"}},
	}
	elements := buildSelectableElements(info)
	state := SelectionState{
		Elements: elements,
		Selected: make([]bool, len(elements)),
		PageInfo: info,
	}

	// Select some elements: title (0), first h1 (1), first paragraph (3)
	state.Selected[0] = true // title
	state.Selected[1] = true // h1
	state.Selected[3] = true // first paragraph

	result, err := handleFinishWithSelections(state)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.Finished {
		t.Error("Expected result to be finished")
	}

	if result.SelectedData == nil {
		t.Fatal("Expected SelectedData to be populated")
	}

	// Check title
	if title, ok := result.SelectedData["title"].(string); !ok || title != "Test Title" {
		t.Errorf("Expected title 'Test Title' in SelectedData, got %v", result.SelectedData["title"])
	}

	// Check H1
	if h1List, ok := result.SelectedData["h1"].([]string); !ok || len(h1List) != 1 || h1List[0] != "H1 Title" {
		t.Errorf("Expected H1 ['H1 Title'] in SelectedData, got %v", result.SelectedData["h1"])
	}

	// Check paragraphs
	if paragraphs, ok := result.SelectedData["paragraphs"].([]string); !ok || len(paragraphs) != 1 || paragraphs[0] != "Paragraph 1" {
		t.Errorf("Expected paragraphs ['Paragraph 1'] in SelectedData, got %v", result.SelectedData["paragraphs"])
	}

	// Check that H2 is NOT in the results (wasn't selected)
	if _, exists := result.SelectedData["h2"]; exists {
		t.Error("Expected H2 to NOT be in SelectedData")
	}

	// Test with no selections
	stateEmpty := SelectionState{
		Elements: elements,
		Selected: make([]bool, len(elements)),
		PageInfo: info,
	}

	result, err = handleFinishWithSelections(stateEmpty)
	if err != nil {
		t.Fatalf("Unexpected error for empty selection: %v", err)
	}

	if result.Finished {
		t.Error("Expected result to NOT be finished when no elements selected")
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

func TestTruncateText(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly ten chars", 17, "exactly ten chars"},
		{"this is a very long text that should be truncated", 20, "this is a very lo..."},
		{"", 5, ""},
		{"abc", 3, "abc"},
		{"abcd", 3, "..."}, // 4 chars with max 3 should become "..."
		{"a very long text", 10, "a very ..."},
	}

	for _, test := range tests {
		result := truncateText(test.input, test.maxLen)
		if result != test.expected {
			t.Errorf("truncateText(%q, %d) = %q, expected %q", test.input, test.maxLen, result, test.expected)
		}
	}
}
