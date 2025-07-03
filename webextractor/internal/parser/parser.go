package parser

import (
	"strings"

	"golang.org/x/net/html"
)

// MatchFunc decides if a given node matches our selector.
type MatchFunc func(n *html.Node) bool

// Compile builds a MatchFunc from a simple selector string.
// Supported syntax:
//
//	tag          — matches element tag name
//	.class       — matches element with given class (space separated list)
//	#id          — matches element with given id.
func Compile(selector string) MatchFunc {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return func(_ *html.Node) bool { return false }
	}
	switch selector[0] {
	case '.':
		class := selector[1:]
		return func(n *html.Node) bool {
			if n.Type != html.ElementNode {
				return false
			}
			for _, a := range n.Attr {
				if a.Key == "class" {
					classes := strings.Fields(a.Val)
					for _, c := range classes {
						if c == class {
							return true
						}
					}
					return false
				}
			}
			return false
		}
	case '#':
		id := selector[1:]
		return func(n *html.Node) bool {
			if n.Type != html.ElementNode {
				return false
			}
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == id {
					return true
				}
			}
			return false
		}
	default:
		tag := strings.ToLower(selector)
		return func(n *html.Node) bool {
			return n.Type == html.ElementNode && n.Data == tag
		}
	}
}

// FindAll traverses the DOM tree depth-first and returns nodes that match the selector.
func FindAll(root *html.Node, selector string) []*html.Node {
	matcher := Compile(selector)
	var out []*html.Node
	var rec func(*html.Node)
	rec = func(n *html.Node) {
		if matcher(n) {
			out = append(out, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			rec(c)
		}
	}
	rec(root)
	return out
}

// TextContent returns the concatenation of all text descendants of n.
func TextContent(n *html.Node) string {
	var b strings.Builder
	var rec func(*html.Node)
	rec = func(nd *html.Node) {
		if nd.Type == html.TextNode {
			b.WriteString(strings.TrimSpace(nd.Data))
			b.WriteString(" ")
		}
		for c := nd.FirstChild; c != nil; c = c.NextSibling {
			rec(c)
		}
	}
	rec(n)
	return strings.TrimSpace(b.String())
}

// Link represents a hyperlink with its text and URL.
type Link struct {
	Href string
	Text string
}

// FindLinks traverses the HTML tree and extracts all hyperlinks.
func FindLinks(n *html.Node) []Link {
	var links []Link
	if n.Type == html.ElementNode && n.Data == "a" {
		var href string
		for _, a := range n.Attr {
			if a.Key == "href" {
				href = a.Val
				break
			}
		}
		if href != "" {
			links = append(links, Link{
				Href: href,
				Text: strings.TrimSpace(TextContent(n)),
			})
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, FindLinks(c)...)
	}
	return links
}
