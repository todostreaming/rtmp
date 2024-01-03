package conn

import (
	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/mock"
)

type MockChunker struct {
	mock.Mock
}

func (c *MockChunker) Chunk(m Marshallable) (*chunk.Chunk, error) {
	args := c.Called(m)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*chunk.Chunk), args.Error(1)
}
