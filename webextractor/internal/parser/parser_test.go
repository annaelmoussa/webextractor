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
