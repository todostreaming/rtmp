package conn

import (
	"io"

	"github.com/todostreaming/amf0"
	"github.com/stretchr/testify/mock"
)

type MockParser struct {
	mock.Mock
}

var _ Parser = new(MockParser)

func (p *MockParser) Parse(name *amf0.String, r io.Reader) (Receivable, error) {
	args := p.Called(name, r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(Receivable), args.Error(1)
}
