package data_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/todostreaming/rtmp/cmd/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewStreamConstructsNewStreams(t *testing.T) {
	s := data.NewStream(make(chan *chunk.Chunk), chunk.NoopWriter)

	assert.IsType(t, new(data.Stream), s)
}

func TestRecvPushesDataWhenSuccessful(t *testing.T) {
	s := data.NewStream(make(chan *chunk.Chunk), chunk.NoopWriter)

	parser := &MockParser{}
	parser.On("Parse", mock.Anything).Return(new(data.Audio), nil).Once()
	s.SetParser(parser)

	go s.Recv()
	s.Chunks() <- new(chunk.Chunk)

	assert.Equal(t, new(data.Audio), <-s.Out())
	parser.AssertExpectations(t)
}

func TestRecvEmitsAnErrorWhenNotSuccessful(t *testing.T) {
	s := data.NewStream(make(chan *chunk.Chunk), chunk.NoopWriter)

	parser := &MockParser{}
	parser.On("Parse", mock.Anything).Return(nil, errors.New("foo")).Once()
	s.SetParser(parser)

	go s.Recv()
	s.Chunks() <- new(chunk.Chunk)

	assert.Equal(t, "foo", (<-s.Errs()).Error())
}

func TestRecvWritesToChunkWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	writer := chunk.NewWriter(buf, 4096)

	ch := &chunk.Chunk{
		Header: new(chunk.Header),
		Data:   []byte{0x0, 0x1, 0x2, 0x3},
	}

	d := new(MockData)
	d.On("Marshal").Return(ch, nil).Once()

	s := data.NewStream(make(chan *chunk.Chunk), writer)
	go s.Recv()

	s.In() <- d

	assert.Equal(t, []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01, 0x02, 0x03,
	}, buf.Bytes())
}

type MockData struct {
	mock.Mock
}

var _ data.Data = new(MockData)

func (d *MockData) Id() byte {
	return d.Called().Get(0).(byte)
}

func (d *MockData) Read(c *chunk.Chunk) error {
	return d.Called(c).Error(0)
}

func (d *MockData) Marshal() (*chunk.Chunk, error) {
	args := d.Called()

	return args.Get(0).(*chunk.Chunk), args.Error(1)
}
