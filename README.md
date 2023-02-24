# goldmark-anchor

[![Go Reference](https://pkg.go.dev/badge/go.abhg.dev/goldmark/anchor.svg)](https://pkg.go.dev/go.abhg.dev/goldmark/anchor)
[![CI](https://github.com/abhinav/goldmark-anchor/actions/workflows/ci.yml/badge.svg)](https://github.com/abhinav/goldmark-anchor/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/abhinav/goldmark-anchor/branch/main/graph/badge.svg?token=NnC4qdjuDW)](https://codecov.io/gh/abhinav/goldmark-anchor)

goldmark-anchor is an extension for the [goldmark] Markdown parser
that adds support for anchors next to all headers.

  [goldmark]: http://github.com/yuin/goldmark

**Demo**:
A web-based demonstration of the extension is available at
<https://abhinav.github.io/goldmark-anchor/demo/>.

## Installation

```bash
go get go.abhg.dev/goldmark/anchor@latest
```

## Usage

To use goldmark-anchor, import the `anchor` package.

```go
import "go.abhg.dev/goldmark/anchor"
```

Then include the `anchor.Extender` in the list of extensions
when constructing your [`goldmark.Markdown`].

  [`goldmark.Markdown`]: https://pkg.go.dev/github.com/yuin/goldmark#Markdown

```go
goldmark.New(
  goldmark.WithParserOptions(
    parser.WithAutoHeadingID(), // read note
  ),
  goldmark.WithExtensions(
    // ...
    &anchor.Extender{},
  ),
).Convert(src, out)
```

This will begin adding '¶' anchors next to all headers in your Markdown files.

> **NOTE**:
> The example above adds the `parser.WithAutoHeadingID` option.
> Without this, or a custom implementation of `parser.IDs`,
> Goldmark will not add `id` attributes to headers.
> If a header does not have an `id`,
> then goldmark-anchor will not generate an anchor for it.

### Changing anchor text

Change the anchor text by setting the `Texter` field
of the `Extender` to an `anchor.Text` value.

```go
&anchor.Extender{
  Texter: anchor.Text("#"),
}
```

#### Dynamic anchor text

You can dynamically calculate anchor text
by supplying a custom implementation of `Texter`
to the `Extender`.

For example, the following `Texter` repeats '#' matching the header level,
providing anchors similar to Markdown `#`-style headers.

```go
type customTexter struct{}

func (*customTexter) AnchorText(h *anchor.HeaderInfo) []byte {
  return bytes.Repeat([]byte{'#'}, h.Level)
}

// Elsewhere:

&anchor.Extender{
  Texter: &customTexter{},
}
```

### Skipping headers

To skip headers, supply a custom `Texter` that returns an empty output
for the `AnchorText` method.

The following `Texter` will not render anchors for level 1 headers.

```go
type customTexter struct{}

func (*customTexter) AnchorText(h *anchor.HeaderInfo) []byte {
  if h.Level == 1 {
    return nil
  }
  return []byte("#")
}
```

### Changing anchor attributes

Change the anchor attributes by setting the `Attributer` field
of the `Extender`.

```go
&anchor.Extender{
  Attributer: anchor.Attributes{
    "class": "permalink",
  },
}
```

By default, goldmark-anchor will add `class="anchor"` to all anchors.
Set `Attributer` to `anchor.Attributes{}` to disable this.

```go
&anchor.Extender{
  Attributer: anchor.Attributes{},
}
```

### Changing anchor positioning

Anchors can appear either at the start of the header before the header text,
or at the end after the header text.

```html
<!-- Before -->
<h1><a href="#">#</a> Foo</h1>

<!-- After -->
<h1>Foo <a href="#">¶</a></h1>
```

You can choose the placement by setting the `Position` field of `Extender`.

```go
&anchor.Extender{
  Position: anchor.Before, // or anchor.After
}
```

By default, goldmark-anchor will place anchors after the header text.

## FAQ

### Why are no anchors being generated?

By default, Goldmark does not generate IDs for headings.
Since goldmark-anchor generates anchors only for headings with IDs,
this can result in no anchors being generated

You must enable heading ID generation using one of the following methods:

- set the [`parser.WithAutoHeadingID`] option
- supply your own [`parser.IDs`] implementation

Alternatively, if your document specifies heading attributes manually,
enable the [`parser.WithHeadingAttribute`] option and manually specify
heading IDs next to each heading.

  [`parser.WithAutoHeadingID`]: https://pkg.go.dev/github.com/yuin/goldmark/parser#WithAutoHeadingID
  [`parser.IDs`]: https://pkg.go.dev/github.com/yuin/goldmark/parser#IDs
  [`parser.WithHeadingAttribute`]: https://pkg.go.dev/github.com/yuin/goldmark/parser#WithHeadingAttribute

## License

This software is made available under the MIT license.
