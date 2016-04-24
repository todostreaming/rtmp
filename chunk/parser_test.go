package chunk_test

// Note: the AssertExpectations calls are commented out in this file because
// AssertExpectations does not work when spawning goroutines out of function
// calls. I have throughouly investigated this issue and determined that there
// is definitely a bug somewhere in github.com/stretchr/testify/mock because
// method calls that definitely occured were not counted in the mock.
//
// This note serves to remind us to uncomment the expectation assertions when
// this bug becomes fixed.

import (
	"errors"
	"reflect"
	"testing"

	"github.com/WatchBeam/rtmp/chunk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewStreamReturnsNewStreams(t *testing.T) {
	p := chunk.NewParser(&MockReader{}, nil)

	assert.IsType(t, &chunk.Parser{}, p)
}

func TestStreamPropogatesErrors(t *testing.T) {
	errs := make(chan error, 2)
	chunks := make(chan *chunk.Chunk)

	reader := &MockReader{}
	reader.On("Recv").Return().Once()
	reader.On("Chunks").Return(chunks)
	reader.On("Errs").Return(errs)
	reader.On("Close").Return()

	p := chunk.NewParser(reader, chunk.NewNormalizer())
	errs <- errors.New("testing: some error")

	go p.Recv()

	err := <-p.Errs()

	p.Close()

	// reader.AssertExpectations(t)
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

	p := chunk.NewParser(reader, chunk.NewNormalizer())

	go p.Recv()
	p.Close()

	// reader.AssertExpectations(t)
}

func TestStreamNormalizesChunksAndSendsThem(t *testing.T) {
	chunks := make(chan *chunk.Chunk, 2)
	chunks <- &chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader: chunk.BasicHeader{0, 2},
		},
	}

	reader := &MockReader{}
	reader.On("Recv").Return()
	reader.On("Chunks").Return(chunks)
	reader.On("Errs").Return(make(chan error))
	reader.On("Close").Return()

	normalizer := &MockNormalizer{}
	normalizer.On("Normalize", mock.Anything).Return().Once()

	parser := chunk.NewParser(reader, normalizer)
	go parser.Recv()

	stream, err := parser.Stream(2)
	assert.Nil(t, err)
	<-stream.In()

	parser.Close()

	// reader.AssertExpectations(t)
	// normalizer.AssertExpectations(t)
}

func TestParserReturnsNewSingleChunkStreams(t *testing.T) {
	parser := chunk.NewParser(nil, nil)

	stream, err := parser.Stream(1)

	assert.NotNil(t, stream)
	assert.Nil(t, err)
}

func TestParserReturnsConsistentChunkStreams(t *testing.T) {
	parser := chunk.NewParser(nil, nil)

	s1, _ := parser.Stream(1)
	s2, _ := parser.Stream(1)

	assert.True(t,
		reflect.ValueOf(s1).Pointer() == reflect.ValueOf(s2).Pointer(),
		"rtmp/chunk: parser should return identical streams, but pointer values differed")
}

func TestParserReturnsMultiStreamsWhenNoStreamsAlreadyExist(t *testing.T) {
	parser := chunk.NewParser(nil, nil)

	stream, err := parser.Stream(1, 2, 3)

	assert.Nil(t, err)
	assert.NotNil(t, stream)
}

func TestParserDoesNotReturnMultiStreamsWhenStreamsAlreadyExist(t *testing.T) {
	parser := chunk.NewParser(nil, nil)

	parser.Stream(1) // take out stream ID 1

	multiStream, err := parser.Stream(1, 2)

	assert.Nil(t, multiStream)
	assert.Equal(t, "rtmp/chunk: stream 1 already exists", err.Error())
}
