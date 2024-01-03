package chunk_test

import (
	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/mock"
)

type MockStream struct {
	mock.Mock
}

func (s *MockStream) In() <-chan *chunk.Chunk {
	args := s.Called()

	return args.Get(0).(chan *chunk.Chunk)
}
