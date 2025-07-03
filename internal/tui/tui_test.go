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

	if info.URL != "https://test.com" {
		t.Errorf("Expected URL 'https://test.com', got '%s'", info.URL)
	}

	if info.Title != "Test Page" {
		t.Errorf("Expected title 'Test Page', got '%s'", info.Title)
	}

	if len(info.H1) != 1 || info.H1[0] != "Main Title" {
		t.Errorf("Expected H1 ['Main Title'], got %v", info.H1)
	}

	if len(info.H2) != 1 || info.H2[0] != "Subtitle" {
		t.Errorf("Expected H2 ['Subtitle'], got %v", info.H2)
	}

	if len(info.Paragraphs) != 2 {
		t.Errorf("Expected 2 paragraphs, got %d", len(info.Paragraphs))
	}

	if len(info.Links) != 2 {
		t.Errorf("Expected 2 links, got %d", len(info.Links))
	}

	if len(info.Images) != 1 {
		t.Errorf("Expected 1 image, got %d", len(info.Images))
	}
	if info.Images[0].Alt != "Test Image" {
		t.Errorf("Expected image alt 'Test Image', got '%s'", info.Images[0].Alt)
	}

	if len(info.Lists) != 1 {
		t.Errorf("Expected 1 list, got %d", len(info.Lists))
	}
}

func TestPrintFunctions(t *testing.T) {
	elements := []SelectableElement{
		{Index: 0, Type: "title", Content: "Test Title"},
		{Index: 1, Type: "h1", Content: "Test H1"},
		{Index: 2, Type: "p", Content: "Test paragraph"},
	}

	selected := make([]bool, len(elements))
	selected[0] = true

	pageInfo := PageInfo{
		URL:   "https://test.com",
		Title: "Test Page",
	}

	state := SelectionState{
		Elements: elements,
		Selected: selected,
		PageInfo: pageInfo,
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printSelectableElements panicked: %v", r)
		}
	}()

	printSelectableElements(state)
	printSelectionStatus(state)
	printSelectionMenu()
	printHelp()
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

	expectedCount := 9
	if len(elements) != expectedCount {
		t.Errorf("Expected %d elements, got %d", expectedCount, len(elements))
	}

	if elements[0].Type != "title" || elements[0].Data.(string) != "Test Title" {
		t.Errorf("Expected first element to be title 'Test Title', got type=%s data=%v", elements[0].Type, elements[0].Data)
	}

	for i, elem := range elements {
		if elem.Index != i {
			t.Errorf("Expected element %d to have index %d, got %d", i, i, elem.Index)
		}
	}

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
		{"0", 5, []int{0}, false},
		{"3", 5, []int{3}, false},
		{"4", 5, []int{4}, false},
		{"0,2,4", 5, []int{0, 2, 4}, false},
		{"1,3", 5, []int{1, 3}, false},
		{"0-2", 5, []int{0, 1, 2}, false},
		{"1-3", 5, []int{1, 2, 3}, false},
		{"0,2-4", 5, []int{0, 2, 3, 4}, false},
		{"0,3,1-2", 5, []int{0, 3, 1, 2}, false},

		{"", 5, nil, true},
		{"5", 5, nil, true},
		{"-1", 5, nil, true},
		{"0-5", 5, nil, true},
		{"3-1", 5, nil, true},
		{"a", 5, nil, true},
		{"0,a", 5, nil, true},
		{"0-", 5, nil, true},
		{"0-2-4", 5, nil, true},
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

		if len(result) != len(test.expected) {
			t.Errorf("For input '%s': expected %d indices, got %d", test.input, len(test.expected), len(result))
			continue
		}

		resultSet := make(map[int]bool)
		for _, idx := range result {
			resultSet[idx] = true
		}

		expectedSet := make(map[int]bool)
		for _, idx := range test.expected {
			expectedSet[idx] = true
		}

		for expectedIdx := range expectedSet {
			if !resultSet[expectedIdx] {
				t.Errorf("For input '%s': expected index %d not found in result %v", test.input, expectedIdx, result)
			}
		}

		for resultIdx := range resultSet {
			if !expectedSet[resultIdx] {
				t.Errorf("For input '%s': unexpected index %d found in result %v", test.input, resultIdx, result)
			}
		}
	}
}

func TestHandleIndexSelection(t *testing.T) {
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

	err := handleIndexSelection("0,2", &state)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !state.Selected[0] {
		t.Error("Expected index 0 to be selected")
	}
	if !state.Selected[2] {
		t.Error("Expected index 2 to be selected")
	}
	if state.Selected[1] {
		t.Error("Expected index 1 to NOT be selected")
	}

	err = handleIndexSelection("0", &state)
	if err != nil {
		t.Fatalf("Unexpected error for already selected index: %v", err)
	}

	err = handleIndexSelection("99", &state)
	if err == nil {
		t.Error("Expected error for out of bounds index")
	}

	err = handleIndexSelection("invalid", &state)
	if err == nil {
		t.Error("Expected error for invalid format")
	}
}

func TestHandleFinishWithSelections(t *testing.T) {
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

	state.Selected[0] = true
	state.Selected[1] = true
	state.Selected[3] = true

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

	if title, ok := result.SelectedData["title"].(string); !ok || title != "Test Title" {
		t.Errorf("Expected title 'Test Title' in SelectedData, got %v", result.SelectedData["title"])
	}

	if h1List, ok := result.SelectedData["h1"].([]string); !ok || len(h1List) != 1 || h1List[0] != "H1 Title" {
		t.Errorf("Expected H1 ['H1 Title'] in SelectedData, got %v", result.SelectedData["h1"])
	}

	if paragraphs, ok := result.SelectedData["paragraphs"].([]string); !ok || len(paragraphs) != 1 || paragraphs[0] != "Paragraph 1" {
		t.Errorf("Expected paragraphs ['Paragraph 1'] in SelectedData, got %v", result.SelectedData["paragraphs"])
	}

	if _, exists := result.SelectedData["h2"]; exists {
		t.Error("Expected H2 to NOT be in SelectedData")
	}

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

	result, err := handleLinkNavigation("L0", links)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.NextURL != "https://example.com" {
		t.Errorf("Expected NextURL 'https://example.com', got '%s'", result.NextURL)
	}

	_, err = handleLinkNavigation("L", links)
	if err == nil {
		t.Error("Expected error for invalid format")
	}

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
		{"abcd", 3, "..."},
		{"a very long text", 10, "a very ..."},
	}

	for _, test := range tests {
		result := truncateText(test.input, test.maxLen)
		if result != test.expected {
			t.Errorf("truncateText(%q, %d) = %q, expected %q", test.input, test.maxLen, result, test.expected)
		}
	}
}
