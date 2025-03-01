package anchor

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Extender adds support for anchors to a Goldmark Markdown parser.
//
// Use it by installing it into the [goldmark.Markdown] object upon creation.
// For example:
//
//	goldmark.New(
//		// ...
//		goldmark.WithExtensions(
//			// ...
//			&anchor.Extender{},
//		),
//	)
type Extender struct {
	// Texter determines the anchor text.
	//
	// Defaults to '¶' if unspecified.
	Texter Texter

	// Position specifies where the anchor will be placed in a header.
	//
	// Defaults to After.
	Position Position

	// Attributer determines the attributes
	// that will be associated with the anchor link.
	//
	// Defaults to adding a 'class="anchor"' attribute.
	Attributer Attributer

	// Unsafe specifies whether the Texter values will be escaped or not.
	// Setting this to true can lead to HTML injection if you don't handle
	// Texter values with care.
	//
	// Defaults to false.
	Unsafe bool
}

var _ goldmark.Extender = (*Extender)(nil)

// Extend extends the provided Goldmark Markdown.
func (e *Extender) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&Transformer{
				Texter:     e.Texter,
				Position:   e.Position,
				Attributer: e.Attributer,
			}, 100),
		),
	)
	md.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&Renderer{
				Position: e.Position,
				Unsafe:   e.Unsafe,
			}, 100),
		),
	)
}
