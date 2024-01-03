package chunk_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

func TestWriterConstruction(t *testing.T) {
	w := chunk.NewWriter(new(bytes.Buffer), chunk.DefaultReadSize)

	assert.IsType(t, &chunk.DefaultWriter{}, w)
}

func TestWritingOnce(t *testing.T) {
	buf := new(bytes.Buffer)
	c := &chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader:       chunk.BasicHeader{0, 18},
			MessageHeader:     chunk.MessageHeader{0, 1234, false, 8, 2, 3},
			ExtendedTimestamp: chunk.ExtendedTimestamp{},
		},
		Data: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
	}

	err := chunk.NewWriter(buf, 8).Write(c)

	assert.Nil(t, err)

	expected := new(bytes.Buffer)
	(&chunk.BasicHeader{0, 18}).Write(expected)
	(&chunk.MessageHeader{0, 1234, false, 8, 2, 3}).Write(expected)
	expected.Write([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07})

	assert.True(t, bytes.Equal(expected.Bytes(), buf.Bytes()),
		fmt.Sprintf("test: slice should be equal (%v, %v)", expected.Bytes(),
			buf.Bytes()))
}

func TestMultipleWrites(t *testing.T) {
	buf := new(bytes.Buffer)
	c := &chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader:       chunk.BasicHeader{0, 18},
			MessageHeader:     chunk.MessageHeader{0, 1234, false, 8, 2, 3},
			ExtendedTimestamp: chunk.ExtendedTimestamp{},
		},
		Data: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
	}

	err := chunk.NewWriter(buf, 4).Write(c)

	assert.Nil(t, err)

	expected := new(bytes.Buffer)
	(&chunk.BasicHeader{0, 18}).Write(expected)
	(&chunk.MessageHeader{0, 1234, false, 8, 2, 3}).Write(expected)
	expected.Write([]byte{0x00, 0x01, 0x02, 0x03})
	expected.Write([]byte{byte((3 << 6) | 18&63)})
	expected.Write([]byte{0x04, 0x05, 0x06, 0x07})

	assert.True(t, bytes.Equal(expected.Bytes(), buf.Bytes()),
		fmt.Sprintf("test: slice should be equal (%v, %v)", expected.Bytes(),
			buf.Bytes()))
}
