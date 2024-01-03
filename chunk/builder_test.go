package chunk_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

func TestNewReturnsNewBuilder(t *testing.T) {
	h := &chunk.Header{}
	b := chunk.NewBuilder(h)

	assert.IsType(t, &chunk.Builder{}, b)
	assert.Equal(t, h, b.Header)
}

func TestReadAppendsBytes(t *testing.T) {
	slice := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	builder := chunk.NewBuilder(&chunk.Header{
		MessageHeader: chunk.MessageHeader{Length: 8},
	})

	n, err := builder.Read(bytes.NewBuffer(slice), 8)

	assert.Nil(t, err)
	assert.Equal(t, 8, n)
	assert.Equal(t, slice, builder.Payloads[0])
}

func TestReadDoesNotAppendFailedReads(t *testing.T) {
	slice := []byte{}
	builder := chunk.NewBuilder(&chunk.Header{
		MessageHeader: chunk.MessageHeader{Length: 8},
	})

	n, err := builder.Read(bytes.NewBuffer(slice), 8)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, n)
	assert.Empty(t, builder.Payloads)
}

func TestAppendAddsSingleSliceWithinBounds(t *testing.T) {
	b := chunk.NewBuilder(&chunk.Header{
		MessageHeader: chunk.MessageHeader{Length: 8},
	})
	slice := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}

	n, err := b.Append(slice)

	assert.Nil(t, err)
	assert.Equal(t, 8, n)
	assert.Equal(t, slice, b.Payloads[0])
}

func TestAppendAddsNSliceWithinBounds(t *testing.T) {
	b := chunk.NewBuilder(&chunk.Header{
		MessageHeader: chunk.MessageHeader{Length: 8},
	})

	for i, slice := range [][]byte{
		[]byte{0x00, 0x01, 0x02, 0x03},
		[]byte{0x04, 0x05, 0x06, 0x07},
	} {
		n, err := b.Append(slice)

		assert.Nil(t, err)
		assert.Equal(t, len(slice), n)
		assert.Equal(t, slice, b.Payloads[i])
	}

	assert.Equal(t, 2, len(b.Payloads))
}

func TestChunkBuildersBuildValidChunks(t *testing.T) {
	header := new(chunk.Header)
	builder := &chunk.Builder{
		Header: header,
		Payloads: [][]byte{
			[]byte{0x01}, []byte{0x02},
		},
	}

	c := builder.Build()

	assert.IsType(t, &chunk.Chunk{}, c)
	assert.Equal(t, header, c.Header)
	assert.Equal(t, []byte{0x01, 0x02}, c.Data)
}
