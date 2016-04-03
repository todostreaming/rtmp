package data

import (
	"fmt"
	"reflect"

	"github.com/WatchBeam/rtmp/chunk"
)

var (
	// DefaultParser is a singleton instance of the Parser type (using the
	// SimpleParser type as implementation) that contains references to both
	// Data implementations: Audio and Video.
	DefaultParser = NewParser([]Data{
		new(Audio), new(Video),
	})
)

// SimpleParser provides a default implementation of the Parser eype.
type SimpleParser struct {
	typs map[byte]reflect.Type
}

// NewParser creates and returns an instance of the *SimpleParser type. It is
// initialized with the given Data implementations.
func NewParser(typs []Data) *SimpleParser {
	p := &SimpleParser{
		typs: make(map[byte]reflect.Type),
	}

	for _, t := range typs {
		p.typs[t.Id()] = reflect.TypeOf(t).Elem()
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

	if err := d.Read(c.Data); err != nil {
		return nil, err
	}

	return d, nil
}

// New returns a new instance of the Data implementation keyed by the given ID.
// If no matching Data implementation was found, then a value of a nil is
// returned instead.
func (p *SimpleParser) New(id byte) Data {
	typ := p.typs[id]
	if typ == nil {
		return nil
	}

	return reflect.New(typ).Interface().(Data)
}
