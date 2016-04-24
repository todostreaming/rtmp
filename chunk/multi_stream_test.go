package chunk_test

import (
	"testing"

	"github.com/WatchBeam/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

func TestMultiStreamAppendsMultipleChunkStreams(t *testing.T) {
	s1, c1 := newMockStream()
	s2, c2 := newMockStream()

	ms := chunk.NewMultiStream()
	ms.Append(s1, s2)

	o1 := &chunk.Chunk{Data: []byte{1}}
	o2 := &chunk.Chunk{Data: []byte{2}}

	go func() {
		c2 <- o2
		c1 <- o1
	}()

	assert.Equal(t, o2, <-ms.In())
	assert.Equal(t, o1, <-ms.In())
}

func newMockStream() (*MockStream, chan *chunk.Chunk) {
	s := new(MockStream)
	c := make(chan *chunk.Chunk)

	s.On("In").Return(c).Once()

	return s, c
}
