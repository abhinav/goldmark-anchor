package anchor

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var (
	_defaultTexter     = Text("¶")
	_defaultAttributer = Attributes{"class": "anchor"}
)

// HeaderInfo holds information about a header
// for which an anchor is being considered.
type HeaderInfo struct {
	// Level of the header.
	Level int

	// Identifier for the header on the page.
	// This will typically become part of the URL fragment.
	ID []byte
}

// Texter determines the anchor text.
//
// This is the clickable text displayed next to the header
// which tells readers that they can use it as an anchor to the header.
//
// By default, we will use the string '¶'.
type Texter interface {
	// AnchorText returns the anchor text
	// that should be used for the provided header info.
	//
	// If AnchorText returns an empty slice or nil,
	// an anchor will not be generated for this header.
	AnchorText(*HeaderInfo) []byte
}

// Text builds a Texter that uses a constant string
// as the anchor text.
//
// Pass this into [Extender] or [Transformer]
// to specify a custom anchor text.
//
//	anchor.Extender{
//		Texter: Text("#"),
//	}
func Text(s string) Texter {
	return textTexter(s)
}

type textTexter []byte

func (t textTexter) AnchorText(*HeaderInfo) []byte {
	return []byte(t)
}

// Position specifies where inside a heading we should place an anchor [Node].
type Position int

//go:generate stringer -type Position

const (
	// After places the anchor node after the heading text.
	//
	// This is the default.
	After Position = iota

	// Before places the anchor node before the heading text.
	Before
)

// Attributer determines attributes that will be attached to an anchor node.
//
// By default, we will add 'class="anchor"' to all nodes.
type Attributer interface {
	// AnchorAttributes returns the attributes
	// that should be attached to the anchor node
	// for the given header.
	//
	// If AnchorAttributes returns an empty map or nil,
	// no attributes will be added.
	AnchorAttributes(*HeaderInfo) map[string]string
}

// Attributes is an Attributer that uses a constant set of attributes
// for all anchor nodes.
//
// Pass this into [Extender] or [Transformer] to specify custom attributes.
//
//	anchor.Extender{
//		Attributer: Attributes{"class": "permalink"},
//	}
type Attributes map[string]string

var _ Attributer = Attributes{}

// AnchorAttributes reports the attributes associated with this object
// for all headers.
func (as Attributes) AnchorAttributes(*HeaderInfo) map[string]string {
	return as
}

// Transformer transforms a Goldmark Markdown AST,
// adding anchor [Node] objects for headers across the document.
type Transformer struct {
	// Texter determines the anchor text.
	//
	// Defaults to '¶' for all headers if unset.
	Texter Texter

	// Position specifies where the anchor will be placed in a header.
	//
	// Defaults to After.
	Position Position

	// Attributer determines the attributes
	// that will be associated with the anchor link.
	//
	// Defaults to adding a 'class="anchor"' attribute
	// for all headers if unset.
	Attributer Attributer
}

var _ parser.ASTTransformer = (*Transformer)(nil)

// Transform traverses and transforms the provided Markdown document.
//
// This method is typically called by Goldmark
// and should not need to be invoked directly.
func (t *Transformer) Transform(doc *ast.Document, reader text.Reader, pctx parser.Context) {
	tr := transform{
		Attributer: t.Attributer,
		Position:   t.Position,
		Texter:     t.Texter,
	}
	if tr.Attributer == nil {
		tr.Attributer = _defaultAttributer
	}
	if tr.Texter == nil {
		tr.Texter = _defaultTexter
	}

	ast.Walk(doc, tr.Visit)
}

// transform holds state for a single transformation traversal.
type transform struct {
	Texter     Texter
	Position   Position
	Attributer Attributer
}

func (t *transform) Visit(n ast.Node, enter bool) (ast.WalkStatus, error) {
	if !enter {
		return ast.WalkContinue, nil
	}
	h, ok := n.(*ast.Heading)
	if !ok {
		return ast.WalkContinue, nil
	}

	t.transform(h)
	return ast.WalkSkipChildren, nil
}

func (t *transform) transform(h *ast.Heading) {
	idattr, ok := h.AttributeString("id")
	if !ok {
		return
	}

	id, ok := idattr.([]byte)
	if !ok {
		return
	}

	info := HeaderInfo{
		Level: h.Level,
		ID:    id,
	}

	text := t.Texter.AnchorText(&info)
	if len(text) == 0 {
		return
	}

	n := &Node{
		ID:    id,
		Level: h.Level,
		Value: text,
	}

	for name, value := range t.Attributer.AnchorAttributes(&info) {
		n.SetAttributeString(name, []byte(value))
	}

	// If the header has no children yet, just append the anchor.
	if h.ChildCount() == 0 {
		h.AppendChild(h, n)
		return
	}

	if t.Position == Before {
		h.InsertBefore(h, h.FirstChild(), n)
	} else {
		h.InsertAfter(h, h.LastChild(), n)
	}
}
