package htmlparser

import (
	"io"
	"strings"
)

// Parse analyse le HTML depuis un reader et retourne le nœud racine.
func Parse(r io.Reader) (*Node, error) {
	tokenizer := NewTokenizer(r)
	doc := &Node{Type: DocumentNode}
	stack := []*Node{doc}

	for {
		tokenType := tokenizer.Next()
		if tokenType == ErrorToken {
			break
		}

		// On récupère le token actuel
		token := tokenizer.Token()
		// On récupère le nœud courant
		current := stack[len(stack)-1]

		switch tokenType {
		case TextToken:
			if strings.TrimSpace(token.Data) != "" {
				textNode := &Node{
					Type: TextNode,
					Data: token.Data,
				}
				current.AppendChild(textNode)
			}

		case StartTagToken:
			element := &Node{
				Type: ElementNode,
				Data: token.Data,
				Attr: token.Attr,
			}
			current.AppendChild(element)
			// Si la balise n'est pas auto-fermante, on l'ajoute à la pile
			if !isSelfClosing(token.Data) {
				// On ajoute le nœud à la pile
				stack = append(stack, element)
			}

		case SelfClosingTagToken:
			element := &Node{
				Type: ElementNode,
				Data: token.Data,
				Attr: token.Attr,
			}
			current.AppendChild(element)

		case EndTagToken:
			if len(stack) > 1 {
				if stack[len(stack)-1].Data == token.Data {
					stack = stack[:len(stack)-1]
				} else {
					for i := len(stack) - 1; i >= 1; i-- {
						if stack[i].Data == token.Data {
							stack = stack[:i]
							break
						}
					}
				}
			}

		case CommentToken:
			commentNode := &Node{
				Type: CommentNode,
				Data: token.Data,
			}
			current.AppendChild(commentNode)
		}
	}

	return doc, nil
}

// isSelfClosing retourne true pour les balises HTML auto-fermantes.
func isSelfClosing(tag string) bool {
	// On définit les balises auto-fermantes
	selfClosing := map[string]bool{
		"area": true, "base": true, "br": true, "col": true,
		"embed": true, "hr": true, "img": true, "input": true,
		"link": true, "meta": true, "param": true, "source": true,
		"track": true, "wbr": true,
	}
	// On retourne true si la balise est auto-fermante
	return selfClosing[strings.ToLower(tag)]
}
