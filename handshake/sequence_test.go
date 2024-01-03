package handshake_test

import (
	"io"

	"github.com/todostreaming/rtmp/handshake"
	"github.com/stretchr/testify/mock"
)

type MockSequence struct {
	mock.Mock
}

var _ handshake.Sequence = new(MockSequence)

func (s *MockSequence) Read(r io.Reader) error {
	return s.Called(r).Error(0)
}

func (s *MockSequence) WriteTo(w io.Writer) error {
	return s.Called(w).Error(0)
}

func (s *MockSequence) Next() handshake.Sequence {
	next := s.Called().Get(0)
	if next == nil {
		return nil
	}

	return next.(handshake.Sequence)
}
