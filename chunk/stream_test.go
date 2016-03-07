package chunk_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/WatchBeam/rtmp/chunk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewStreamReturnsNewStreams(t *testing.T) {
	s := chunk.NewStream(&MockReader{}, &MockWriter{}, nil)

	assert.IsType(t, &chunk.Stream{}, s)
}

func TestStreamPropogatesErrors(t *testing.T) {
	errs := make(chan error, 2)
	chunks := make(chan *chunk.Chunk)

	reader := &MockReader{}
	reader.On("Recv").Return().Once()
	reader.On("Chunks").Return(chunks)
	reader.On("Errs").Return(errs)
	reader.On("Close").Return()

	s := chunk.NewStream(reader, &MockWriter{}, chunk.NewNormalizer())
	errs <- errors.New("testing: some error")

	go s.Recv()
	defer s.Close()

	err := <-s.Errs()

	assert.Equal(t, "testing: some error", err.Error())
}

func TestStreamClosesAfterSignalSent(t *testing.T) {
	errs := make(chan error, 2)
	chunks := make(chan *chunk.Chunk)

	reader := &MockReader{}
	reader.On("Recv").Return().Once()
	reader.On("Chunks").Return(chunks)
	reader.On("Errs").Return(errs)
	reader.On("Close").Return().Once()

	s := chunk.NewStream(reader, &MockWriter{}, chunk.NewNormalizer())
	go s.Recv()

	s.Close()
	<-time.After(1 * time.Millisecond) // HACK: wait for all the things

	reader.AssertExpectations(t)
}

func TestStreamNormalizesChunksAndSendsThem(t *testing.T) {
	chunks := make(chan *chunk.Chunk, 2)
	chunks <- new(chunk.Chunk)

	reader := &MockReader{}
	reader.On("Recv").Return()
	reader.On("Chunks").Return(chunks)
	reader.On("Errs").Return(make(chan error))
	reader.On("Close").Return()

	normalizer := &MockNormalizer{}
	normalizer.On("Normalize", mock.Anything).Return().Once()

	stream := chunk.NewStream(reader, &MockWriter{}, normalizer)
	go stream.Recv()

	<-stream.In()
	stream.Close()
	<-time.After(1 * time.Millisecond) // HACK: wait for all the things

	reader.AssertExpectations(t)
	normalizer.AssertExpectations(t)
}

func TestStreamWritesChunksToWriter(t *testing.T) {
	writer := &MockWriter{}
	writer.On("Write", mock.Anything).Return(nil)

	stream := chunk.NewStream(
		chunk.NewReader(new(bytes.Buffer), chunk.DefaultReadSize),
		writer,
		&MockNormalizer{},
	)
	go stream.Recv()

	stream.Out() <- &chunk.Chunk{}

	writer.AssertExpectations(t)
}

func TestStreamPropogatesWriteErrors(t *testing.T) {
	writer := &MockWriter{}
	writer.On("Write", mock.Anything).Return(errors.New("test: error"))

	stream := chunk.NewStream(
		chunk.NewReader(new(bytes.Buffer), chunk.DefaultReadSize),
		writer,
		&MockNormalizer{},
	)
	go stream.Recv()

	stream.Out() <- &chunk.Chunk{}
	err := <-stream.Errs()

	assert.Equal(t, "test: error", err.Error())
	writer.AssertExpectations(t)
}
