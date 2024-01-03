package control

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/todostreaming/rtmp/chunk"
)

// UnknownControlType is an Error representing a scenario where an unknown
// Control Type ID was read from an io.Reader.
type UnknownControlType byte

var _ error = new(UnknownControlType)

// Error implements the `func Error` in the `type error interface`.
func (e UnknownControlType) Error() string {
	return fmt.Sprintf("control: unknown control type (%v)", byte(e))
}

// DefaultParser provides a default implementation of the Parser type.
type DefaultParser struct {
	// controls maps control sequence IDs to their respective reflect.Type
	controls map[byte]reflect.Type
}

var _ Parser = new(DefaultParser)

// NewParser returns a new instance of the Parser type (using the DefaultParser
// implementation) initialized with the Controls variable.
func NewParser() *DefaultParser {
	p := &DefaultParser{
		controls: make(map[byte]reflect.Type),
	}

	for _, c := range Controls {
		p.controls[c.TypeId()] = reflect.TypeOf(c).Elem()
	}

	return p
}

// Parse implements the Parse function as defined in the Parser interface.
func (p *DefaultParser) Parse(chunk *chunk.Chunk) (Control, error) {
	id := chunk.Header.MessageHeader.TypeId

	t := p.TypeFor(id)
	if t == nil {
		return nil, UnknownControlType(id)
	}

	c := reflect.New(t).Interface().(Control)
	if err := c.Read(bytes.NewBuffer(chunk.Data)); err != nil {
		return nil, err
	}

	return c, nil
}

// TypeFor returns the de-referenced reflect.Type assosicated with a given
// Control Sequence ID. If no matching type is found, nil is returned instead.
func (p *DefaultParser) TypeFor(id byte) reflect.Type {
	return p.controls[id]
}
