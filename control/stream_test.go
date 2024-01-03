package control_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/todostreaming/rtmp/control"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	TestChunk = &chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader:       chunk.BasicHeader{0, 2},
			MessageHeader:     chunk.MessageHeader{0, 1234, false, 8, 2, 3},
			ExtendedTimestamp: chunk.ExtendedTimestamp{},
		},
		Data: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
	}
)

func newStreamWithChunk(streamId uint32, chunks ...*chunk.Chunk) chunk.Stream {
	buf := new(bytes.Buffer)
	for _, c := range chunks {
		chunk.NewWriter(buf, chunk.DefaultReadSize).Write(c)
	}

	parser := chunk.NewParser(chunk.NewReader(
		buf, chunk.DefaultReadSize, chunk.NoopNormalizer,
	))
	go parser.Recv()

	st, _ := parser.Stream(streamId)
	return st
}

func TestStreamConstruction(t *testing.T) {
	s := control.NewStream(nil, nil, nil, nil)

	assert.IsType(t, &control.Stream{}, s)
}

func TestSuccesfulChunkParsingPushesAControl(t *testing.T) {
	parser := &MockParser{}
	parser.On("Parse", mock.Anything).Return(&control.Acknowledgement{}, nil)

	stream := control.NewStream(newStreamWithChunk(2, TestChunk), nil, parser, nil)
	go stream.Recv()

	ctrl := <-stream.In()

	assert.Equal(t, &control.Acknowledgement{}, ctrl)
	parser.AssertExpectations(t)
}

func TestFailedChunkParsingPushesAnError(t *testing.T) {
	parser := &MockParser{}
	parser.On("Parse", mock.Anything).
		Return(new(control.Acknowledgement), errors.New("test"))

	stream := control.NewStream(newStreamWithChunk(2, TestChunk), nil, parser, nil)
	go stream.Recv()

	err := <-stream.Errs()

	assert.Equal(t, "test", err.Error())
	parser.AssertExpectations(t)
}

func TestWritingAControlChunksIt(t *testing.T) {
	ctrl := new(control.Acknowledgement)

	chunker := &MockChunker{}
	chunker.On("Chunk", ctrl).Return(&chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader:   chunk.BasicHeader{0, 2},
			MessageHeader: chunk.MessageHeader{},
		},
	}, nil)

	stream := control.NewStream(
		newStreamWithChunk(2), chunk.NewWriter(ioutil.Discard,
			chunk.DefaultReadSize), nil, chunker,
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
		newStreamWithChunk(2), nil, nil, chunker,
	)
	go stream.Recv()

	stream.Out() <- ctrl

	assert.Equal(t, "test", (<-stream.Errs()).Error())
	chunker.AssertExpectations(t)
}
