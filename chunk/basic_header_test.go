package chunk_test

import (
	"bytes"
	"testing"

	"github.com/WatchBeam/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

func TestBasicHeaderTypeOneRead(t *testing.T) {
	h := new(chunk.BasicHeader)
	h.Read(bytes.NewBuffer([]byte{
		(0x02 << 6) | 23,
	}))

	assert.Equal(t, byte(2), h.FormatId)
	assert.Equal(t, uint32(23), h.StreamId)
}

func TestBasicHeaderTypeTwoRead(t *testing.T) {
	h := new(chunk.BasicHeader)
	h.Read(bytes.NewBuffer([]byte{
		(0x02 << 6) | 0, 23,
	}))

	assert.Equal(t, byte(2), h.FormatId)
	assert.Equal(t, uint32(23+64), h.StreamId)
}

func TestBasicHeaderTypeThreeRead(t *testing.T) {
	h := new(chunk.BasicHeader)
	h.Read(bytes.NewBuffer([]byte{
		(0x02 << 6) | 63, byte(255 >> 8), 255,
	}))
	assert.Equal(t, byte(2), h.FormatId)
	assert.Equal(t, uint32(255+64), h.StreamId)
}
