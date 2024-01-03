package chunk_test

import (
	"bytes"
	"testing"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

func TestHeaderReading(t *testing.T) {
	h := &chunk.Header{}

	timestamp := uint32(0xffffff)
	extended := uint32(1234)

	buf := bytes.NewBuffer([]byte{
		// Basic Header
		(0x02 << 6) | 63, byte(255 >> 8), 255,

		// Message header
		byte(timestamp >> 16), byte(timestamp >> 8), byte(timestamp),

		// Extended Timetsmp
		byte(extended >> 24), byte(extended >> 16), byte(extended >> 8),
		byte(extended),
	})

	err := h.Read(buf)

	assert.Nil(t, err)

	assert.Equal(t, byte(2), h.BasicHeader.FormatId)
	assert.Equal(t, uint32(255+64), h.BasicHeader.StreamId)
	assert.Equal(t, uint32(0xffffff), h.MessageHeader.Timestamp)
	assert.Equal(t, uint32(1234), h.ExtendedTimestamp.Delta)
}
