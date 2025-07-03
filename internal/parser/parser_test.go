package parser

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

const sampleHTML = `<html><body><div id="main" class="content highlight"><p>First</p><p class="note">Second</p></div><span class="note">Third</span></body></html>`

func getDoc() *html.Node {
	doc, _ := html.Parse(strings.NewReader(sampleHTML))
	return doc
}

func TestCompileAndFindAll(t *testing.T) {
	doc := getDoc()

	tests := []struct {
		sel   string
		count int
	}{
		{"div", 1},
		{"p", 2},
		{".note", 2},
		{"#main", 1},
		{"span", 1},
	}
	for _, tc := range tests {
		nodes := FindAll(doc, tc.sel)
		if len(nodes) != tc.count {
			t.Fatalf("selector %s expected %d got %d", tc.sel, tc.count, len(nodes))
		}
	}
}

func TestTextContent(t *testing.T) {
	doc := getDoc()
	nodes := FindAll(doc, "p")
	if len(nodes) != 2 {
		t.Fatalf("expected 2 p nodes")
	}
	txt := TextContent(nodes[0])
	if txt != "First" {
		t.Fatalf("got %s", txt)
	}
}

func TestFindLinks(t *testing.T) {
	htmlWithLinks := `
	<html>
		<body>
			<a href="https://example.com">Example Link</a>
			<a href="/internal-link">Internal Link</a>
			<a href="mailto:test@example.com">Email Link</a>
			<a>Link without href</a>
			<div>
				<a href="#section">Section Link</a>
			</div>
		</body>
	</html>
	`

	doc, err := html.Parse(strings.NewReader(htmlWithLinks))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	links := FindLinks(doc)

	expectedCount := 4
	if len(links) != expectedCount {
		t.Fatalf("Expected %d links, got %d", expectedCount, len(links))
	}

	expectedLinks := map[string]string{
		"https://example.com":     "Example Link",
		"/internal-link":          "Internal Link",
		"mailto:test@example.com": "Email Link",
		"#section":                "Section Link",
	}

	found := make(map[string]string)
	for _, link := range links {
		found[link.Href] = link.Text
	}

	for expectedHref, expectedText := range expectedLinks {
		if text, exists := found[expectedHref]; !exists {
			t.Errorf("Expected link %s not found", expectedHref)
		} else if text != expectedText {
			t.Errorf("Expected text '%s' for link %s, got '%s'", expectedText, expectedHref, text)
		}
	}
}
