package parser

import (
	"strings"

	"webextractor/internal/htmlparser"
)

// MatchFunc décide si un nœud donné correspond à notre sélecteur.
type MatchFunc func(n *htmlparser.Node) bool

// Compile construit une MatchFunc depuis une chaîne de sélecteur simple.
func Compile(selector string) MatchFunc {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return func(_ *htmlparser.Node) bool { return false }
	}
	switch selector[0] {
	case '.':
		class := selector[1:]
		return func(n *htmlparser.Node) bool {
			if n.Type != htmlparser.ElementNode {
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
		return func(n *htmlparser.Node) bool {
			if n.Type != htmlparser.ElementNode {
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
		return func(n *htmlparser.Node) bool {
			return n.Type == htmlparser.ElementNode && n.Data == tag
		}
	}
}

// FindAll parcourt l'arbre DOM en profondeur et retourne les nœuds qui correspondent au sélecteur.
func FindAll(root *htmlparser.Node, selector string) []*htmlparser.Node {
	matcher := Compile(selector)
	var out []*htmlparser.Node
	var rec func(*htmlparser.Node)
	rec = func(n *htmlparser.Node) {
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

// TextContent retourne la concaténation de tous les descendants texte de n.
func TextContent(n *htmlparser.Node) string {
	var b strings.Builder
	var rec func(*htmlparser.Node)
	rec = func(nd *htmlparser.Node) {
		if nd.Type == htmlparser.TextNode {
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

// Link représente un hyperlien avec son texte et son URL.
type Link struct {
	Href string
	Text string
}

// FindLinks parcourt l'arbre HTML et extrait tous les hyperliens.
func FindLinks(n *htmlparser.Node) []Link {
	var links []Link
	if n.Type == htmlparser.ElementNode && n.Data == "a" {
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
