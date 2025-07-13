package htmlparser

// NodeType represents the type of a node in the HTML tree.
type NodeType int

const (
	ErrorNode NodeType = iota
	TextNode
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
)

// Attribute represents an HTML attribute.
type Attribute struct {
	Key string
	Val string
}

// Node represents a node in the HTML tree.
type Node struct {
	Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node

	Type NodeType
	Data string
	Attr []Attribute
}

// AppendChild adds a child node to the end of the given node's children.
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

// RemoveChild removes a child node from the given node.
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