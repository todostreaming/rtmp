package chunk_test

import (
	"bytes"
	"testing"

	"github.com/WatchBeam/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

func TestWriterConstruction(t *testing.T) {
	w := chunk.NewWriter(new(bytes.Buffer), chunk.DefaultReadSize)

	assert.IsType(t, &chunk.DefaultWriter{}, w)
}
