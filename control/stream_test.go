package control_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/WatchBeam/rtmp/chunk"
	"github.com/WatchBeam/rtmp/control"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	TestChunk = &chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader:       chunk.BasicHeader{0, 18},
			MessageHeader:     chunk.MessageHeader{0, 1234, false, 8, 2, 3},
			ExtendedTimestamp: chunk.ExtendedTimestamp{},
		},
		Data: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
	}
)

func newStreamWithChunk(c ...*chunk.Chunk) *chunk.Stream {
	buf := new(bytes.Buffer)

	w := chunk.NewWriter(buf, chunk.DefaultReadSize)
	for _, cc := range c {
		w.Write(cc)
	}

	chunks := chunk.NewStream(
		chunk.NewReader(buf, chunk.DefaultReadSize), nil,
		chunk.NewNormalizer())
	go chunks.Recv()

	return chunks
}

func TestStreamConstruction(t *testing.T) {
	s := control.NewStream(nil, nil, nil)

	assert.IsType(t, &control.Stream{}, s)
}

func TestSuccesfulChunkParsingPushesAControl(t *testing.T) {
	parser := &MockParser{}
	parser.On("Parse", mock.Anything).Return(&control.Acknowledgement{}, nil)

	stream := control.NewStream(newStreamWithChunk(TestChunk), parser, nil)
	go stream.Recv()

	ctrl := <-stream.In()

	assert.Equal(t, &control.Acknowledgement{}, ctrl)
	parser.AssertExpectations(t)
}

func TestFailedChunkParsingPushesAnError(t *testing.T) {
	parser := &MockParser{}
	parser.On("Parse", mock.Anything).
		Return(new(control.Acknowledgement), errors.New("test"))

	stream := control.NewStream(newStreamWithChunk(TestChunk), parser, nil)
	go stream.Recv()

	err := <-stream.Errs()

	assert.Equal(t, "test", err.Error())
	parser.AssertExpectations(t)
}

func TestWritingAControlChunksIt(t *testing.T) {
	ctrl := new(control.Acknowledgement)

	chunker := &MockChunker{}
	chunker.On("Chunk", ctrl).Return(new(chunk.Chunk), nil)

	stream := control.NewStream(
		newStreamWithChunk(), nil, chunker,
	)
	go stream.Recv()

	stream.Out() <- ctrl

	chunker.AssertExpectations(t)
}

func TestWritingAControlErrorsWhenErrored(t *testing.T) {
	ctrl := new(control.Acknowledgement)

	chunker := &MockChunker{}
	chunker.On("Chunk", ctrl).Return(new(chunk.Chunk), errors.New("test"))

	stream := control.NewStream(
		newStreamWithChunk(), nil, chunker,
	)
	go stream.Recv()

	stream.Out() <- ctrl

	assert.Equal(t, "test", (<-stream.Errs()).Error())
	chunker.AssertExpectations(t)
}
