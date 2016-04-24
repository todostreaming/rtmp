package chunk_test

import (
	"reflect"
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

	o1 := new(chunk.Chunk)
	o2 := new(chunk.Chunk)

	go func() {
		c2 <- o2
		c1 <- o1
	}()

	assert.True(t, pointerEquality(o2, <-ms.In()), "did not expect chunk #2")
	assert.True(t, pointerEquality(o1, <-ms.In()), "did not expect chunk #1")
}

func pointerEquality(v1, v2 interface{}) bool {
	return reflect.ValueOf(v1).Pointer() == reflect.ValueOf(v2).Pointer()
}
