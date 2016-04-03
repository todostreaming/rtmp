package data_test

import (
	"errors"
	"testing"

	"github.com/WatchBeam/rtmp/chunk"
	"github.com/WatchBeam/rtmp/cmd/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewStreamConstructsNewStreams(t *testing.T) {
	s := data.NewStream()

	assert.IsType(t, new(data.Stream), s)
}

func TestRecvPushesDataWhenSuccessful(t *testing.T) {
	s := data.NewStream()

	parser := &MockParser{}
	parser.On("Parse", mock.Anything).Return(new(data.Audio), nil).Once()
	s.SetParser(parser)

	go s.Recv()
	s.Chunks() <- new(chunk.Chunk)

	assert.Equal(t, new(data.Audio), <-s.In())
	parser.AssertExpectations(t)
}

func TestRecvEmitsAnErrorWhenNotSuccessful(t *testing.T) {
	s := data.NewStream()

	parser := &MockParser{}
	parser.On("Parse", mock.Anything).Return(nil, errors.New("foo")).Once()
	s.SetParser(parser)

	go s.Recv()
	s.Chunks() <- new(chunk.Chunk)

	assert.Equal(t, "foo", (<-s.Errs()).Error())
}
