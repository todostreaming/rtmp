package cmd

import (
	"testing"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

type (
	TrueGate  struct{}
	FalseGate struct{}
)

var (
	_ Gate = new(TrueGate)
	_ Gate = new(FalseGate)
)

func (t *TrueGate) Open(_ *chunk.Chunk) bool  { return true }
func (f *FalseGate) Open(_ *chunk.Chunk) bool { return false }

func TestStreamIdGateIsOpenForMatchingStreams(t *testing.T) {
	gate := &StreamIdGate{1}

	open := gate.Open(&chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader: chunk.BasicHeader{0, 1},
		},
	})

	assert.True(t, open)
}

func TestStreamIdGateIsClosedForMismatchedStreams(t *testing.T) {
	gate := &StreamIdGate{1}

	open := gate.Open(&chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader: chunk.BasicHeader{0, 2},
		},
	})

	assert.False(t, open)
}

func TestTypeIdGateIsOpenForMatchingTypes(t *testing.T) {
	gate := &TypeIdGate{1}

	open := gate.Open(&chunk.Chunk{
		Header: &chunk.Header{
			MessageHeader: chunk.MessageHeader{TypeId: 1},
		},
	})

	assert.True(t, open)
}

func TestTypeIdGateIsClosedForMismatchedTypes(t *testing.T) {
	gate := &TypeIdGate{1}

	open := gate.Open(&chunk.Chunk{
		Header: &chunk.Header{
			MessageHeader: chunk.MessageHeader{TypeId: 2},
		},
	})

	assert.False(t, open)
}

func TestMessageStreamGateIsOpenForMatchingStreams(t *testing.T) {
	gate := &MessageStreamGate{1}

	open := gate.Open(&chunk.Chunk{
		Header: &chunk.Header{
			MessageHeader: chunk.MessageHeader{StreamId: 1},
		},
	})

	assert.True(t, open)
}

func TestMessageStreamGateIsClosedForMismatchedStreams(t *testing.T) {
	gate := &MessageStreamGate{1}

	open := gate.Open(&chunk.Chunk{
		Header: &chunk.Header{
			MessageHeader: chunk.MessageHeader{StreamId: 2},
		},
	})

	assert.False(t, open)
}

func TestUnionGateIsOpenWhenAllChildrenAreOpen(t *testing.T) {
	gate := NewUnionGate(new(TrueGate), new(TrueGate))

	open := gate.Open(new(chunk.Chunk))

	assert.True(t, open)
}

func TestUnionGateIsClosedWhenAnyChildrenAreClosed(t *testing.T) {
	gate := NewUnionGate(new(TrueGate), new(FalseGate))

	open := gate.Open(new(chunk.Chunk))

	assert.False(t, open)
}

func TestUnionGateIsClosedWithNoChildren(t *testing.T) {
	gate := NewUnionGate()

	open := gate.Open(new(chunk.Chunk))

	assert.False(t, open)
}

func TestAnyGateIsOpenWhenAnyChildrenAreOpen(t *testing.T) {
	gate := NewAnyGate(new(TrueGate), new(FalseGate))

	open := gate.Open(new(chunk.Chunk))

	assert.True(t, open)
}

func TestAnyGateIsClosedWhenNoChildrenAreOpen(t *testing.T) {
	gate := NewAnyGate(new(FalseGate), new(FalseGate))

	open := gate.Open(new(chunk.Chunk))

	assert.False(t, open)
}
