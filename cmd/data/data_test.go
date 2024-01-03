package data

import (
	"testing"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

func TestDataReadsChunk(t *testing.T) {
	c := &chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader: chunk.BasicHeader{0, 1},
		},
		Data: []byte{0x0, 0x1, 0x2, 0x3, 0x4},
	}
	d := new(data)

	err := d.Read(c)

	assert.Nil(t, err)
	assert.Equal(t, c.Header, d.header)
	assert.Equal(t, byte(0x0), d.Control())
	assert.Equal(t, []byte{0x1, 0x2, 0x3, 0x4}, d.Payload())
}

func TestDataReadsWithEmptyBody(t *testing.T) {
	d := new(data)

	err := d.Read(&chunk.Chunk{
		Data: []byte{0x0},
	})

	assert.Nil(t, err)
	assert.Equal(t, byte(0x0), d.Control())
	assert.Empty(t, d.Payload())
}

func TestDataDoesNotReadWhenMissingControl(t *testing.T) {
	d := new(data)
	err := d.Read(&chunk.Chunk{
		Data: []byte{},
	})

	assert.Equal(t, ErrControlMissing, err)
}
