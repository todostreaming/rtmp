package stream

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type MockParser struct {
	mock.Mock
}

var _ Parser = new(MockParser)

func (p *MockParser) Parse(r io.Reader) (Command, error) {
	args := p.Called(r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(Command), args.Error(1)
}
