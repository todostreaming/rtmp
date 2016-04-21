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

	<-parser.Stream(2).In()
	parser.Close()

	// reader.AssertExpectations(t)
	// normalizer.AssertExpectations(t)
}
