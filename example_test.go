package anchor_test

import (
	"log"
	"os"

	"github.com/yuin/goldmark/v2"
	"github.com/yuin/goldmark/v2/parser"
	"go.abhg.dev/goldmark/anchor"
)

func Example() {
	md := goldmark.New(
		// We need to enable automatic generation of heading IDs.
		// Otherwise, none of the headings will have IDs
		// which will leave goldmark-anchor
		// nothing to generate anchors for.
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			&anchor.Extender{},
		),
	)

	src := []byte("# Foo")
	if err := md.Convert(src, os.Stdout); err != nil {
		log.Fatal(err)
	}

	// Output:
	// <h1 id="foo">Foo <a class="anchor" href="#foo">¶</a></h1>
}
