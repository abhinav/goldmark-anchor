// demo implements a WASM module that can be used to format markdown
// with the goldmark-anchor extension.
package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/anchor"
)

func main() {
	js.Global().Set("formatMarkdown", js.FuncOf(formatMarkdown))
	select {}
}

func formatMarkdown(_ js.Value, args []js.Value) any {
	input := args[0].Get("markdown").String()
	var pos anchor.Position
	switch s := args[0].Get("position").String(); s {
	case "before":
		pos = anchor.Before
	case "after":
		pos = anchor.After
	default:
		return fmt.Sprintf("invalid position: %q", s)
	}

	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			&anchor.Extender{
				Position: pos,
			},
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		return err.Error()
	}
	return buf.String()
}
