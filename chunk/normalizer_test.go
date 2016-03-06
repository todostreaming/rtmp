package chunk_test

import (
	"github.com/WatchBeam/rtmp/chunk"
	"github.com/stretchr/testify/mock"
)

type MockNormalizer struct {
	mock.Mock
}

var _ chunk.Normalizer = new(MockNormalizer)

func (n *MockNormalizer) Normalize(chunk *chunk.Chunk) {
	n.Called(chunk)
}

func (n *MockNormalizer) Last() *chunk.Chunk {
	args := n.Called()
	return args.Get(0).(*chunk.Chunk)
}

func (n *MockNormalizer) SetLast(chunk *chunk.Chunk) {
	n.Called(chunk)
}

func (n *MockNormalizer) Header(streamId uint32) *chunk.Header {
	args := n.Called(streamId)
	return args.Get(0).(*chunk.Header)
}

func (n *MockNormalizer) StoreHeader(header *chunk.Header) {
	n.Called(header)
}
