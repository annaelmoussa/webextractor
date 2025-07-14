package htmlparser

// NodeType représente le type d'un nœud dans l'arbre HTML.
type NodeType int

const (
	TextNode NodeType = iota
	DocumentNode
	ElementNode
	CommentNode
)

// Attribute représente un attribut HTML.
type Attribute struct {
	Key string // On peut ajouter une clé à un attribut exemple (class, id, etc.)
	Val string // On peut ajouter une valeur à un attribut exemple (red, blue, etc.)
}

// Node représente un nœud dans l'arbre HTML.
type Node struct {
	Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node // On peut ajouter des nœuds parents, enfants, etc.

	Type NodeType // On peut ajouter un type à un nœud exemple (TextNode, ElementNode, CommentNode, etc.)
	Data string // On peut ajouter du texte à un nœud exemple (Hello World)
	Attr []Attribute // On peut ajouter des attributs à un nœud exemple (class, id, etc.)
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