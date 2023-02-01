package anchor

import (
	"strconv"

	"github.com/yuin/goldmark/ast"
)

// Kind is the NodeKind used by anchor nodes.
var Kind = ast.NewNodeKind("Anchor")

// Node is an anchor node in the Markdown AST.
type Node struct {
	ast.BaseInline

	// ID of the header this anchor is for.
	ID []byte

	// Level of the header that this anchor is for.
	Level int

	// Value is the text inside the anchor.
	// Typically this is a fixed string
	// like '¶' or '#'.
	Value []byte
}

// Kind reports that this is a Anchor node.
func (*Node) Kind() ast.NodeKind { return Kind }

// Dump dumps this node to stdout for debugging.
func (n *Node) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, map[string]string{
		"ID":    string(n.ID),
		"Value": string(n.Value),
		"Level": strconv.Itoa(n.Level),
	}, nil)
}
