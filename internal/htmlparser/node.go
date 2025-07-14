package htmlparser

// NodeType représente le type d'un nœud dans l'arbre HTML.
type NodeType int

const (
	ErrorNode NodeType = iota
	TextNode
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
)

// Attribute représente un attribut HTML.
type Attribute struct {
	Key string
	Val string
}

// Node représente un nœud dans l'arbre HTML.
type Node struct {
	Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node

	Type NodeType
	Data string
	Attr []Attribute
}

// AppendChild ajoute un nœud enfant à la fin des enfants du nœud donné.
func (n *Node) AppendChild(child *Node) {
	if child.Parent != nil {
		child.Parent.RemoveChild(child)
	}
	child.Parent = n
	if n.LastChild != nil {
		n.LastChild.NextSibling = child
		child.PrevSibling = n.LastChild
	} else {
		n.FirstChild = child
	}
	n.LastChild = child
}

// RemoveChild supprime un nœud enfant du nœud donné.
func (n *Node) RemoveChild(child *Node) {
	if child.Parent != n {
		return
	}
	if n.FirstChild == child {
		n.FirstChild = child.NextSibling
	}
	if child.NextSibling != nil {
		child.NextSibling.PrevSibling = child.PrevSibling
	}
	if n.LastChild == child {
		n.LastChild = child.PrevSibling
	}
	if child.PrevSibling != nil {
		child.PrevSibling.NextSibling = child.NextSibling
	}
	child.Parent = nil
	child.PrevSibling = nil
	child.NextSibling = nil
}