package control_test

import (
	"testing"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/todostreaming/rtmp/control"
	"github.com/stretchr/testify/assert"
)

func TestChunkerConstruction(t *testing.T) {
	c := control.NewChunker()

	assert.IsType(t, &control.DefaultChunker{}, c)
}

func TestChunkingProducesCorrectChunks(t *testing.T) {
	c := control.NewChunker()
	ctrl := &control.Acknowledgement{5}

	out, err := c.Chunk(ctrl)

	assert.Nil(t, err)
	assert.IsType(t, &chunk.Chunk{}, out)

	assert.EqualValues(t, 0, out.Header.BasicHeader.FormatId)
	assert.EqualValues(t, control.ControlChunkStreamId,
		out.Header.BasicHeader.StreamId)

	assert.EqualValues(t, 0, out.Header.MessageHeader.FormatId)
	assert.EqualValues(t, 4, out.Header.MessageHeader.Length)
	assert.EqualValues(t, ctrl.TypeId(), out.Header.MessageHeader.TypeId)
	assert.EqualValues(t, control.ControlMessageStreamId,
		out.Header.MessageHeader.StreamId)

	assert.Equal(t, []byte{
		0, 0, 0, 5,
	}, out.Data)
}
