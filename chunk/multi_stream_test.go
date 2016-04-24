package chunk_test

import (
	"testing"

	"github.com/WatchBeam/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

func TestMultiStreamAppendsMultipleChunkStreams(t *testing.T) {
	c1 := make(chan *chunk.Chunk)
	c2 := make(chan *chunk.Chunk)

	s1 := new(MockStream)
	s1.On("In").Return(c1).Once()

	s2 := new(MockStream)
	s2.On("In").Return(c2).Once()

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
