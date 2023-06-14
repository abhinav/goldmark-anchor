package anchor

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func TestTransform(t *testing.T) {
	t.Parallel()

	const defaultValue = "¶"

	type anchor struct {
		ID       string
		Level    int
		Value    string
		Position Position
	}

	tests := []struct {
		desc string
		give []string
		want []*anchor

		pos  Position
		text Texter
	}{
		{
			desc: "simple",
			give: []string{
				"# Foo bar",
				"",
				"## Bar baz qux",
				"",
				"#### Quux",
			},
			want: []*anchor{
				{
					ID:       "foo-bar",
					Level:    1,
					Value:    defaultValue,
					Position: After,
				},
				{
					ID:       "bar-baz-qux",
					Level:    2,
					Value:    defaultValue,
					Position: After,
				},
				{
					ID:       "quux",
					Level:    4,
					Value:    defaultValue,
					Position: After,
				},
			},
		},
		{
			desc: "custom position and text",
			text: texterFunc(func(i *HeaderInfo) string {
				return strings.Repeat("#", i.Level)
			}),
			pos: Before,
			give: []string{
				"# Foo",
				"",
				"## Bar",
			},
			want: []*anchor{
				{
					ID:       "foo",
					Level:    1,
					Value:    "#",
					Position: Before,
				},
				{
					ID:       "bar",
					Level:    2,
					Value:    "##",
					Position: Before,
				},
			},
		},
		{
			desc: "skip empty id",
			text: texterFunc(func(i *HeaderInfo) string {
				if string(i.ID) == "skip-me" {
					return ""
				}
				return "#"
			}),
			give: []string{
				"# Foo",
				"",
				"## Skip me",
				"",
				"### Bar",
			},
			want: []*anchor{
				{
					ID:       "foo",
					Level:    1,
					Value:    "#",
					Position: After,
				},
				nil,
				{
					ID:       "bar",
					Level:    3,
					Value:    "#",
					Position: After,
				},
			},
		},
		{
			desc: "no title yet",
			give: []string{
				"#",
				"",
			},
			want: []*anchor{
				{
					ID:       "heading",
					Level:    1,
					Value:    defaultValue,
					Position: After,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			p := goldmark.New().Parser()
			p.AddOptions(
				parser.WithAutoHeadingID(),
				parser.WithASTTransformers(
					util.Prioritized(&Transformer{
						Position: tt.pos,
						Texter:   tt.text,
					}, 100),
				),
			)

			src := []byte(strings.Join(tt.give, "\n") + "\n")
			got := p.Parse(text.NewReader(src))

			var gotAnchors []*anchor
			err := ast.Walk(got, func(n ast.Node, enter bool) (ast.WalkStatus, error) {
				if !enter {
					return ast.WalkContinue, nil
				}

				h, ok := n.(*ast.Heading)
				if !ok {
					return ast.WalkContinue, nil
				}

				an, gotPos := findAnchor(h)
				if an == nil {
					gotAnchors = append(gotAnchors, nil)
					return ast.WalkSkipChildren, nil
				}

				gotAnchors = append(gotAnchors, &anchor{
					ID:       string(an.ID),
					Level:    an.Level,
					Value:    string(an.Value),
					Position: gotPos,
				})

				return ast.WalkSkipChildren, nil
			})
			require.NoError(t, err)

			assert.Equal(t, tt.want, gotAnchors)
		})
	}
}

func TestTransform_noHeadingIDs(t *testing.T) {
	t.Parallel()

	p := goldmark.New().Parser()
	p.AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&Transformer{}, 100),
		),
	)

	src := []byte("# Foo\n\n# Bar\n\n# Baz\n")
	got := p.Parse(text.NewReader(src))

	err := ast.Walk(got, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		an, ok := n.(*Node)
		assert.False(t, ok, "unexpected Node: %#v", an)
		return ast.WalkContinue, nil
	})
	require.NoError(t, err)
}

func TestTransform_badIDAttribute(t *testing.T) {
	t.Parallel()

	var h ast.Heading
	h.SetAttributeString("id", []string{"not", "a", "[]byte"})

	tr := transform{
		Texter:     _defaultTexter,
		Attributer: _defaultAttributer,
	}
	assert.NotPanics(t, func() {
		tr.transform(&h)
	})
}

func findAnchor(h *ast.Heading) (an *Node, pos Position) {
	if an, ok := h.LastChild().(*Node); ok {
		return an, After
	}
	if an, ok := h.FirstChild().(*Node); ok {
		return an, Before
	}
	return nil, 0
}

type texterFunc func(*HeaderInfo) string

func (f texterFunc) AnchorText(i *HeaderInfo) []byte {
	return []byte(f(i))
}
