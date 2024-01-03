package stream

import (
	"bytes"
	"errors"
	"testing"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewNetStreamConstructsNetStreams(t *testing.T) {
	s := New(make(chan *chunk.Chunk), chunk.NoopWriter)

	assert.IsType(t, new(NetStream), s)
}

func TestNetStreamParsesChunksSuccessfully(t *testing.T) {
	parser := &MockParser{}
	parser.On("Parse", mock.Anything).
		Return(new(CommandPlay), nil).Once()

	chunks := make(chan *chunk.Chunk)
	s := New(chunks, chunk.NoopWriter)
	s.parser = parser

	go s.Listen()
	chunks <- new(chunk.Chunk)

	cmd := <-s.In()

	parser.AssertExpectations(t)
	assert.Equal(t, new(CommandPlay), cmd)
}

func TestNetStreamPropogatesChunkParsingErrors(t *testing.T) {
	parser := &MockParser{}
	parser.On("Parse", mock.Anything).
		Return(nil, errors.New("foo")).Once()

	chunks := make(chan *chunk.Chunk)
	s := New(chunks, chunk.NoopWriter)
	s.parser = parser

	go s.Listen()
	chunks <- new(chunk.Chunk)

	parser.AssertExpectations(t)
	assert.Equal(t, "foo", (<-s.Errs()).Error())
}

func TestStreamSendsOnStatusUpdates(t *testing.T) {
	buf := new(bytes.Buffer)
	writer := chunk.NewWriter(buf, chunk.DefaultReadSize)

	chunks := make(chan *chunk.Chunk)
	s := New(chunks, writer)

	go s.Listen()

	s.Status() <- NewStatus()

	assert.Equal(t, 0, len(s.Errs()))
	assert.NotEmpty(t, buf.Bytes())
}
