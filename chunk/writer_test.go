package chunk_test

import (
	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/mock"
)

type MockWriter struct {
	mock.Mock
}

var _ chunk.Writer = new(MockWriter)

func (w *MockWriter) Write(c *chunk.Chunk) error {
	args := w.Called(c)
	return args.Error(0)
}

func (w *MockWriter) WriteSize() int {
	args := w.Called()
	return args.Int(0)
}

func (w *MockWriter) SetWriteSize(size int) {
	w.Called(size)
}
