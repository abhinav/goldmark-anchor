package anchor

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNode_Kind(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Kind, new(Node).Kind())
}

func TestNode_Dump(t *testing.T) {
	getStdout := hijackStdout(t)
	(&Node{
		ID:    []byte("foo-bar"),
		Level: 1,
		Value: []byte("#"),
	}).Dump(nil, 0)
	got := getStdout()

	assert.Contains(t, got, "Anchor {\n")
	assert.Contains(t, got, "    Value: #\n")
	assert.Contains(t, got, "    Level: 1\n")
	assert.Contains(t, got, "    ID: foo-bar\n")
	assert.Contains(t, got, "}\n")
}

func hijackStdout(t testing.TB) func() string {
	stdout := os.Stdout
	t.Cleanup(func() {
		os.Stdout = stdout
	})

	r, w, err := os.Pipe()
	done := make(chan struct{})
	var buff bytes.Buffer
	go func() {
		defer close(done)
		_, err := io.Copy(&buff, r)
		assert.NoError(t, err)
	}()

	require.NoError(t, err)
	os.Stdout = w
	return func() string {
		os.Stdout = stdout
		assert.NoError(t, w.Close())
		<-done
		return buff.String()
	}
}
