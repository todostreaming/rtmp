package chunk_test

import (
	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/mock"
)

type MockReader struct {
	mock.Mock
}

var _ chunk.Reader = new(MockReader)

func (r *MockReader) Recv() {
	r.Called()
}
func (r *MockReader) ReadSize() int {
	args := r.Called()
	return args.Int(0)
}

func (r *MockReader) SetReadSize(size int) {
	r.Called(size)
}

func (r *MockReader) Chunks() <-chan *chunk.Chunk {
	args := r.Called()
	return args.Get(0).(chan *chunk.Chunk)
}

func (r *MockReader) Errs() <-chan error {
	args := r.Called()
	return args.Get(0).(chan error)
}

func (r *MockReader) Close() {
	r.Called()
}
