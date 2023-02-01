package anchor_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/anchor"
	"gopkg.in/yaml.v3"
)

func TestIntegration(t *testing.T) {
	t.Parallel()

	testdata, err := os.ReadFile("testdata/tests.yaml")
	require.NoError(t, err)

	var tests []struct {
		Desc string `yaml:"desc"`
		Give string `yaml:"give"`
		Want string `yaml:"want"`

		Pos   string            `yaml:"pos"` // "before" or "after"
		Text  string            `yaml:"text"`
		Attrs map[string]string `yaml:"attrs"`
	}
	require.NoError(t, yaml.Unmarshal(testdata, &tests))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.Desc, func(t *testing.T) {
			t.Parallel()

			var ext anchor.Extender

			if len(tt.Text) > 0 {
				ext.Texter = anchor.Text(tt.Text)
			}

			switch strings.ToLower(tt.Pos) {
			case "":
				// No customization
			case "before":
				ext.Position = anchor.Before
			case "after":
				ext.Position = anchor.After
			default:
				t.Fatalf("unknown position %q", tt.Pos)
			}

			if len(tt.Attrs) > 0 {
				ext.Attributer = anchor.Attributes(tt.Attrs)
			}

			md := goldmark.New(
				goldmark.WithExtensions(&ext),
				goldmark.WithParserOptions(
					parser.WithAutoHeadingID(),
				),
			)

			var got bytes.Buffer
			require.NoError(t, md.Convert([]byte(tt.Give), &got))
			assert.Equal(t,
				strings.TrimSuffix(tt.Want, "\n"),
				strings.TrimSuffix(got.String(), "\n"),
			)
		})
	}
}
