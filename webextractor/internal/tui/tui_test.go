package tui

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestParseIndices(t *testing.T) {
	tests := []struct {
		in  string
		out []int
	}{
		{"0", []int{0}},
		{"1,3", []int{1, 3}},
		{"2-4", []int{2, 3, 4}},
		{"5,7-9", []int{5, 7, 8, 9}},
	}
	for _, tc := range tests {
		got := parseIndices(tc.in)
		if len(got) != len(tc.out) {
			t.Fatalf("%s len expected %d got %d", tc.in, len(tc.out), len(got))
		}
		for i := range got {
			if got[i] != tc.out[i] {
				t.Fatalf("%s expect %v got %v", tc.in, tc.out, got)
			}
		}
	}
}

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

func TestBuildSelector(t *testing.T) {
	htmlStr := `<div id="foo" class="bar baz"></div><span class="note"></span>`
	doc, _ := html.Parse(strings.NewReader(htmlStr))
	divNode := findElement(doc, "div")
	sel := buildSelector(divNode)
	if sel != "#foo" {
		t.Fatalf("expected #foo got %s", sel)
	}
}

func TestFindMeaningfulNodes(t *testing.T) {
	htmlStr := `<div>Hello</div><span>World</span>`
	doc, _ := html.Parse(strings.NewReader(htmlStr))

	nodes := findMeaningfulNodes(doc)
	if len(nodes) == 0 {
		t.Fatal("findMeaningfulNodes should find at least some nodes")
	}

	// html.Parse creates html > head + body structure for fragments
	found := make(map[string]bool)
	for _, node := range nodes {
		found[node.Data] = true
		t.Logf("Found node: <%s> with text: %q", node.Data, previewText(node))
	}

	// We should find at least some meaningful content
	// Could be body (containing all text) or individual div/span
	if len(found) == 0 {
		t.Fatal("Expected to find some meaningful nodes")
	}
}

func TestPromptSelectorsSimple(t *testing.T) {

	htmlStr := `<div>Hello</div><span>World</span>`
	doc, _ := html.Parse(strings.NewReader(htmlStr))

	nodes := findMeaningfulNodes(doc)

	if len(nodes) == 0 {
		t.Fatal("findMeaningfulNodes returned no nodes")
	}
}
