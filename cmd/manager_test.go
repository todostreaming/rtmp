package cmd

import (
	"reflect"
	"testing"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

type MockChunkStream struct {
	C chan *chunk.Chunk
}

var _ chunk.Stream = new(MockChunkStream)

func (m *MockChunkStream) In() <-chan *chunk.Chunk { return m.C }

func TestNewManagerMakesNewManagers(t *testing.T) {
	m := New(nil, nil)

	assert.IsType(t, new(Manager), m)
}

func TestManagerDispatchesToMatchingChannels(t *testing.T) {
	c := new(chunk.Chunk)

	c1 := make(chan *chunk.Chunk)
	c2 := make(chan *chunk.Chunk)

	cs := &MockChunkStream{make(chan *chunk.Chunk)}

	m := New(cs, nil)
	m.channels = map[Gate]chan<- *chunk.Chunk{
		new(TrueGate):  c2,
		new(FalseGate): c1,
	}

	go m.Dispatch(false)

	cs.C <- c

	assert.Equal(t, 0, len(c1))
	assert.Equal(t,
		reflect.ValueOf(c).Pointer(),
		reflect.ValueOf(<-c2).Pointer())
}
