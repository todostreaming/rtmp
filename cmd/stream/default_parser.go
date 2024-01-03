package stream

import (
	"fmt"
	"io"

	"github.com/todostreaming/amf0/encoding"
)

var (
	// DefaultParser is the default, singleton instance of the Parser type.
	// It uses the SimpleParser type for its implementation, and is capable
	// of understanding all commands that are able to be sent over the
	// NetStream connection.
	//
	// For a complete list of commands that are supported, see the list
	// below.
	DefaultParser Parser = NewParser(map[string]CommandFactory{
		"play":         func() Command { return new(CommandPlay) },
		"play2":        func() Command { return new(CommandPlay2) },
		"deleteStream": func() Command { return new(CommandDeleteStream) },
		"receiveAudio": func() Command { return new(CommandReceiveAudio) },
		"receiveVideo": func() Command { return new(CommandReceiveVideo) },
		"publish":      func() Command { return new(CommandPublish) },
		"seek":         func() Command { return new(CommandSeek) },
		"pause":        func() Command { return new(CommandPause) },
	})
)

// CommandFactory is a factory type capabale of producing new instances of
// command types. By contract, the CommandFactory type should be pseudo-pure
// function within its parameterized type. In other words, a CommandFactory that
// returns a CommandPlay instance should always return new, unique instances of
// the CommandPlay type.
type CommandFactory func() Command

// SimpleParser is a simple implementation of the Parser type. Internally, it
// uses a map of strings to CommandFactory types in order to determine the type
// of, and then parse into a new instance of a Command from a stream of bytes.
type SimpleParser struct {
	// typs is the internal table in which the assosciation between strings
	// and CommandFactories is stored.
	typs map[string]CommandFactory
}

var _ Parser = new(SimpleParser)

// NewParser returns a new instance of the Parser type by using the
// *SimpleParser as its implementation.
//
// The returned parser is initialized with the given map[string]CommandFactory.
func NewParser(typs map[string]CommandFactory) *SimpleParser {
	return &SimpleParser{
		typs: typs,
	}
}

// Parse implements the Parse function in `type Parser interface`. It determines
// first the CommandHeader assosciated with the io.Reader, then creates a new
// instance of the corresponding command type and then parses into it.
//
// If an error is encountered in parsing, or if no matching command can be
// found, then an error will be returned.
func (p *SimpleParser) Parse(r io.Reader) (Command, error) {
	meta := new(CommandHeader)
	if err := encoding.Unmarshal(r, meta); err != nil {
		return nil, err
	}

	factory, ok := p.typs[meta.Name]
	if !ok {
		return nil, fmt.Errorf(
			"cmd/stream: unknown NetStream command %s", meta.Name)
	}

	cmd := factory()
	if err := encoding.Unmarshal(r, cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}
