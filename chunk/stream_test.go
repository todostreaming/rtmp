package chunk_test

import (
	"bytes"
	"testing"

	"github.com/WatchBeam/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

func TestNewStreamReturnsNewStreams(t *testing.T) {
	s := chunk.NewStream(new(bytes.Buffer), stream.DefaultReadSize)

	assert.IsType(t, &chunk.Stream{}, s)
}
