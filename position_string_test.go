package anchor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosition_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc string
		give Position
		want string
	}{
		{desc: "before", give: Before, want: "Before"},
		{desc: "after", give: After, want: "After"},
		{desc: "unknown", give: 42, want: "Position(42)"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.give.String())
		})
	}
}
