package chunk_test

import (
	"testing"
	"time"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

const (
	CHAN_CLOSE_TIMEOUT = 100 * time.Millisecond
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

func TestAwaitCloseWaitsForAllChildrenToClose(t *testing.T) {
	done := make(chan bool)

	s1, c1 := newMockStream()
	s2, c2 := newMockStream()

	multi := chunk.NewMultiStream()
	multi.Append(s1, s2)

	go func() {
		multi.AwaitClose()
		done <- true
	}()

	close(c1)

	select {
	case <-time.After(CHAN_CLOSE_TIMEOUT):
	case <-done:
		t.Fatal("rtmp/chunk: MutliStream.AwaitClose should wait for " +
			"children to die")
	}

	close(c2)

	select {
	case <-time.After(CHAN_CLOSE_TIMEOUT):
		t.Fatal("rtmp/chunk: MutliStream.AwaitClose should have " +
			"closed after all children were closed")
	case <-done:
	}
}

func newMockStream() (*MockStream, chan *chunk.Chunk) {
	s := new(MockStream)
	c := make(chan *chunk.Chunk)

	s.On("In").Return(c).Once()

	return s, c
}
