package stream_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/todostreaming/amf0/encoding"
	"github.com/todostreaming/rtmp/cmd/stream"
	"github.com/stretchr/testify/assert"
)

var (
	ValidCommandHeader = []byte{
		0x02, 0x00, 0x08, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75,
		0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x05,
	}
)

func TestCommandHeaderRead(t *testing.T) {
	var header stream.CommandHeader
	buf := bytes.NewReader(ValidCommandHeader)

	err := encoding.Unmarshal(buf, &header)

	assert.Nil(t, err)
	assert.Equal(t, stream.CommandHeader{
		Name:          "onStatus",
		TransactionId: 0,
		Arguments:     nil,
	}, header)
}

func TestInvalidCommandHeaderRead(t *testing.T) {
	var header stream.CommandHeader
	buf := bytes.NewReader([]byte{
	// Invalid payload, empty
	})

	err := encoding.Unmarshal(buf, &header)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, stream.CommandHeader{}, header)
}
