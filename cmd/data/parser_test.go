package data_test

import (
	"github.com/todostreaming/rtmp/chunk"
	"github.com/todostreaming/rtmp/cmd/data"
	"github.com/stretchr/testify/mock"
)

type MockParser struct {
	mock.Mock
}

var _ data.Parser = new(MockParser)

func (m *MockParser) Parse(c *chunk.Chunk) (data.Data, error) {
	args := m.Called(c)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(data.Data), args.Error(1)
}
