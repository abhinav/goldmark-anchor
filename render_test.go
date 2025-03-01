package anchor

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func TestRenderer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc   string
		give   Node
		attrs  map[string]string
		pos    Position
		want   string
		unsafe bool
	}{
		{desc: "empty ID"},
		{
			desc: "after",
			give: Node{
				ID:    []byte("hello"),
				Value: []byte("#"),
			},
			want: ` <a href="#hello">#</a>`,
		},
		{
			desc: "before",
			pos:  Before,
			give: Node{
				ID:    []byte("hello"),
				Value: []byte("#"),
			},
			want: `<a href="#hello">#</a> `,
		},
		{
			desc: "attributes",
			give: Node{
				ID:    []byte("hello"),
				Value: []byte("#"),
			},
			attrs: map[string]string{"foo": "bar"},
			want:  ` <a foo="bar" href="#hello">#</a>`,
		},
		{
			desc: "attributes",
			give: Node{
				ID:    []byte("hello"),
				Value: []byte("<unsafe></unsafe>"),
			},
			attrs: map[string]string{"foo": "bar"},
			want:  ` <a foo="bar" href="#hello">&lt;unsafe&gt;&lt;/unsafe&gt;</a>`,
		},
		{
			desc: "attributes",
			give: Node{
				ID:    []byte("hello"),
				Value: []byte("<unsafe></unsafe>"),
			},
			attrs:  map[string]string{"foo": "bar"},
			want:   ` <a foo="bar" href="#hello"><unsafe></unsafe></a>`,
			unsafe: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			anchorR := Renderer{
				Position: tt.pos,
				Unsafe:   tt.unsafe,
			}
			r := renderer.NewRenderer(
				renderer.WithNodeRenderers(
					util.Prioritized(&anchorR, 100),
				),
			)

			node := tt.give
			for k, v := range tt.attrs {
				node.SetAttributeString(k, []byte(v))
			}

			var buff bytes.Buffer
			require.NoError(t,
				r.Render(&buff, nil /* src */, &node))
			assert.Equal(t, tt.want, buff.String())
		})
	}
}
