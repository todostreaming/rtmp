package chunk_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

func NewReader(r io.Reader) chunk.Reader {
	if r == nil {
		r = new(bytes.Buffer)
	}

	return chunk.NewReader(
		r, chunk.DefaultReadSize, chunk.NoopNormalizer,
	)
}

func TestReaderConstruction(t *testing.T) {
	r := chunk.NewReader(
		new(bytes.Buffer),
		chunk.DefaultReadSize,
		chunk.NoopNormalizer,
	)

	assert.IsType(t, &chunk.DefaultReader{}, r)
}

func TestGetReadSize(t *testing.T) {
	r := NewReader(nil)

	size := r.ReadSize()

	assert.Equal(t, chunk.DefaultReadSize, size)
}

func TestSetReadSize(t *testing.T) {
	r := NewReader(nil)

	r.SetReadSize(6)

	assert.Equal(t, 6, r.ReadSize())
}

func TestReadSingleChunkSinglePass(t *testing.T) {
	b := new(bytes.Buffer)
	c := &chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader:       chunk.BasicHeader{0, 18},
			MessageHeader:     chunk.MessageHeader{0, 1234, false, 8, 2, 3},
			ExtendedTimestamp: chunk.ExtendedTimestamp{},
		},
		Data: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
	}

	chunk.NewWriter(b, chunk.DefaultReadSize).Write(c)

	r := NewReader(b)
	go r.Recv()

	read := <-r.Chunks()

	assert.Equal(t, 0, len(r.Errs()))
	assert.Equal(t, c, read)
}

func TestReadSingleChunkMultiPass(t *testing.T) {
	b := new(bytes.Buffer)
	c := &chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader:       chunk.BasicHeader{0, 18},
			MessageHeader:     chunk.MessageHeader{0, 1234, false, 8, 2, 3},
			ExtendedTimestamp: chunk.ExtendedTimestamp{},
		},
		Data: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
	}

	chunk.NewWriter(b, 4).Write(c)

	r := chunk.NewReader(b, 4, chunk.NoopNormalizer)
	go r.Recv()

	read := <-r.Chunks()

	assert.Equal(t, 0, len(r.Errs()))
	assert.Equal(t, c, read)
}

func TestReadMultiChunkSinglePass(t *testing.T) {
	b := new(bytes.Buffer)
	c1 := &chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader:       chunk.BasicHeader{0, 18},
			MessageHeader:     chunk.MessageHeader{0, 1234, false, 8, 2, 3},
			ExtendedTimestamp: chunk.ExtendedTimestamp{},
		},
		Data: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
	}
	c2 := &chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader:       chunk.BasicHeader{0, 19},
			MessageHeader:     chunk.MessageHeader{0, 1234, false, 8, 2, 3},
			ExtendedTimestamp: chunk.ExtendedTimestamp{},
		},
		Data: []byte{0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
	}

	chunk.NewWriter(b, chunk.DefaultReadSize).Write(c1)
	chunk.NewWriter(b, chunk.DefaultReadSize).Write(c2)

	r := NewReader(b)
	go r.Recv()

	r1 := <-r.Chunks()
	r2 := <-r.Chunks()

	assert.Equal(t, 0, len(r.Errs()))
	assert.Equal(t, c1, r1)
	assert.Equal(t, c2, r2)
}

func TestReadMultiChunkMultiPass(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.Write([]byte{
		18, 0, 4, 210, 0, 0, 8, 2, 3, 0, 0, 0, 0, 1, 2, 3, // Chunk: 1, Part: 1
		19, 0, 4, 210, 0, 0, 8, 2, 3, 0, 0, 0, 8, 9, 10, 11, // Chunk: 2, Part: 1
		byte((3 << 6) | 18&63), 4, 5, 6, 7, // Chunk: 1, Part: 2
		byte((3 << 6) | 19&63), 12, 13, 14, 15,
	})

	c1 := &chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader:       chunk.BasicHeader{0, 18},
			MessageHeader:     chunk.MessageHeader{0, 1234, false, 8, 2, 3},
			ExtendedTimestamp: chunk.ExtendedTimestamp{},
		},
		Data: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
	}
	c2 := &chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader:       chunk.BasicHeader{0, 19},
			MessageHeader:     chunk.MessageHeader{0, 1234, false, 8, 2, 3},
			ExtendedTimestamp: chunk.ExtendedTimestamp{},
		},
		Data: []byte{0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
	}

	r := chunk.NewReader(buf, 4, chunk.NoopNormalizer)
	go r.Recv()

	r1 := <-r.Chunks()
	r2 := <-r.Chunks()

	assert.Equal(t, 0, len(r.Errs()))
	assert.Equal(t, c1, r1)
	assert.Equal(t, c2, r2)
}
