package conn

import (
	"bytes"
	"errors"
	"testing"

	"github.com/todostreaming/amf0"
	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	CreateStreamPayload = []byte{
		0x02, 0x00, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53,
		0x74, 0x72, 0x65, 0x61, 0x6d, 0x00, 0x40, 0x10, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x05,
	}
)

func TestNewNetConnectionMakesNewNetConnections(t *testing.T) {
	cnx := NewNetConnection(nil, nil)

	assert.IsType(t, new(NetConn), cnx)
}

func TestNetConnectionReceivablesAreWrittenOut(t *testing.T) {
	chunks := make(chan *chunk.Chunk, 1)
	chunks <- &chunk.Chunk{
		Data: CreateStreamPayload,
	}

	nc := NewNetConnection(chunks, nil)
	go nc.Listen()

	assert.Equal(t, 0, len(nc.Errs()))
	assert.Equal(t, &CreateStreamCommand{
		TransactionId: 4,
		Metadata:      amf0.NewObject(),
	}, <-nc.In())
}

func TestMalformedNamesPropogateErrors(t *testing.T) {
	chunks := make(chan *chunk.Chunk, 1)
	chunks <- &chunk.Chunk{
		Data: []byte{
			0xff, // <invalid type ID>
		},
	}

	nc := NewNetConnection(chunks, nil)
	go nc.Listen()

	assert.NotNil(t, <-nc.Errs())
}

func TestWronglyTypedNamesPropogateErrors(t *testing.T) {
	chunks := make(chan *chunk.Chunk, 1)
	chunks <- &chunk.Chunk{
		Data: []byte{
			0x05, // <invalid type ID>
		},
	}

	nc := NewNetConnection(chunks, nil)
	go nc.Listen()

	assert.Equal(t,
		"rtmp/conn: wrong type for AMF header: *amf0.Null (expected amf0.String)",
		(<-nc.Errs()).Error())
}

func TestParserErrorsArePropogated(t *testing.T) {
	chunks := make(chan *chunk.Chunk, 1)
	chunks <- &chunk.Chunk{
		Data: CreateStreamPayload,
	}

	p := new(MockParser)
	p.On("Parse", mock.Anything, mock.Anything).Return(
		nil, errors.New("parse err")).Once()

	nc := NewNetConnection(chunks, nil)
	nc.parser = p
	go nc.Listen()

	assert.Equal(t, "parse err", (<-nc.Errs()).Error())

}

func TestSendablesAreWrittenToChunkStream(t *testing.T) {
	buf := new(bytes.Buffer)

	nc := NewNetConnection(
		make(chan *chunk.Chunk),
		chunk.NewWriter(buf, chunk.DefaultReadSize))
	go nc.Listen()

	nc.Out() <- &CreateStreamResponse{
		TransactionId: 1,
	}

	assert.Equal(t, 0, len(nc.Errs()))
	assert.Equal(t, []byte{
		0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1d, 0x14, 0x00, 0x00,
		0x00, 0x00, 0x02, 0x00, 0x07, 0x5f, 0x72, 0x65, 0x73, 0x75,
		0x6c, 0x74, 0x00, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,
	}, buf.Bytes())

}

func TestChunkingErrorsArePropogated(t *testing.T) {
	out := new(CreateStreamResponse)
	chunker := new(MockChunker)
	chunker.On("Chunk", out).Return(nil, errors.New("foo")).Once()

	nc := NewNetConnection(
		make(chan *chunk.Chunk),
		chunk.NewWriter(new(bytes.Buffer), chunk.DefaultReadSize))
	nc.chunker = chunker

	go nc.Listen()

	nc.Out() <- out
	assert.Equal(t, "foo", (<-nc.Errs()).Error())
}
