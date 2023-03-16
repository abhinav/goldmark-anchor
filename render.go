package anchor

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// Renderer renders anchor [Node]s.
type Renderer struct {
	// Position specifies where in the header text
	// the anchor is being added.
	Position Position
}

var _ renderer.NodeRenderer = (*Renderer)(nil)

// RegisterFuncs registers functions against the provided goldmark Registerer.
func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(Kind, r.RenderNode)
}

// RenderNode renders an anchor node.
// Goldmark will invoke this method when it encounters a Node.
func (r *Renderer) RenderNode(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// If position is Before, we need to add the anchor when entering;
	// otherwise when exiting.
	if (r.Position == Before) != entering {
		return ast.WalkContinue, nil
	}

	n := node.(*Node)
	if len(n.ID) == 0 {
		return ast.WalkContinue, nil
	}

	// Add leading/trailing ' ' depending on position.
	if r.Position == Before {
		defer w.WriteByte(' ')
	} else {
		w.WriteByte(' ')
	}

	w.WriteString("<a")
	html.RenderAttributes(w, node, nil)
	w.WriteString(` href="#`)
	w.Write(util.EscapeHTML(n.ID))
	w.WriteString(`">`)
	w.Write(util.EscapeHTML(n.Value))
	w.WriteString("</a>")

	return ast.WalkContinue, nil
}
