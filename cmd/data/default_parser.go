package data

import (
	"fmt"

	"github.com/todostreaming/amf0"
	"github.com/todostreaming/rtmp/chunk"
)

var (
	// DefaultParser is a singleton instance of the Parser type (using the
	// SimpleParser type as implementation) that contains references to both
	// Data implementations: Audio and Video.
	DefaultParser = NewParser(
		func() Data { return &Audio{} },
		func() Data { return &Video{} },
		func() Data { return &DataFrame{Arguments: amf0.NewArray()} },
	)
)

// DataFactory is a factory type that produces new instances of a given Data
// type.
type DataFactory func() Data

// SimpleParser provides a default implementation of the Parser eype.
type SimpleParser struct {
	typs map[byte]DataFactory
}

// NewParser creates and returns an instance of the *SimpleParser type. It is
// initialized with the given Data implementations.
func NewParser(factories ...DataFactory) *SimpleParser {
	p := &SimpleParser{
		typs: make(map[byte]DataFactory),
	}

	for _, f := range factories {
		p.typs[f().Id()] = f
	}

	return p
}

var _ Parser = new(SimpleParser)

// Parse implements the Parser.Parser function.
func (p *SimpleParser) Parse(c *chunk.Chunk) (Data, error) {
	d := p.New(c.Header.MessageHeader.TypeId)
	if d == nil {
		return nil, fmt.Errorf("rtmp/data: unknown type ID %v",
			c.Header.MessageHeader.TypeId)
	}

	if err := d.Read(c); err != nil {
		return nil, err
	}

	return d, nil
}

// New returns a new instance of the Data implementation keyed by the given ID.
// If no matching Data implementation was found, then a value of a nil is
// returned instead.
func (p *SimpleParser) New(id byte) Data {
	f := p.typs[id]
	if f == nil {
		return nil
	}

	return f()
}
