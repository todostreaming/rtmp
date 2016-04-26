package chunk

import "github.com/stretchr/testify/mock"

type MockNormalizer struct {
	mock.Mock
}

var _ Normalizer = new(MockNormalizer)

func (n *MockNormalizer) Normalize(h *Header) *Header {
	args := n.Called(h)

	return args.Get(0).(*Header)
}
