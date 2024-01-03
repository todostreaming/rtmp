package conn

import (
	"fmt"
	"io"

	"github.com/todostreaming/amf0"
	"github.com/todostreaming/amf0/encoding"
)

var (
	// DefaultParser is the primary singleton instance of the Parser type.
	// It comes preloaded with all RTMP-related Receivable types, which
	// currently include the connect, createStream, releaseStream, and
	// FCPublish packets.
	//
	// It is recommended that this be used as the primary parcer in any type
	// that requires it.
	DefaultParser Parser = NewParser(map[string]ReceviableFactory{
		"connect": func() Receivable {
			return &ConnectCommand{
				Metadata: amf0.NewObject(),
			}
		},

		"createStream": func() Receivable {
			return &CreateStreamCommand{
				Metadata: amf0.NewObject(),
			}
		},

		"releaseStream": func() Receivable { return new(ReleaseCommand) },

		"FCPublish": func() Receivable { return new(FCPublishCommand) },

		"getStreamLength": func() Receivable {
			return new(GetStreamLength)
		},
	})
)

// ReceviableFactory is a convenience type that returns a new instance of a
// particular Receivable type. As a rule, implementing functions must be
// "pure-ish" meaning they return a consistent sub-type of Receviable.
//
// For examples, see the code used to instantiate the DefaultParser instance
// above.
type ReceviableFactory func() Receivable

// SimpleParser is a default implementation of the Parser type.
type SimpleParser struct {
	// typs maps the type "name" of a Receivable to its corresponding
	// ReceviableFactory so that it can be looked up later and instantiated
	// quickly.
	typs map[string]ReceviableFactory
}

var _ Parser = new(SimpleParser)

// NewParser returns an new instance of the Parser type, using SimpleParser as
// its implementation.
//
// Parsers created using this method are initialized with the given `typs` map.
func NewParser(typs map[string]ReceviableFactory) *SimpleParser {
	return &SimpleParser{
		typs: typs,
	}
}

// Parse implements the Parser.Parse method. It parses a Receivable type out of
// the given AMF identifier (parsed from the io.Reader `r` as a source), or an
// error in the following cases:
//
//   1) no corresponding command could be found
//   2) an error occured during unmarshalling (see todostreaming/rtmp)
//
// Otherwise the Receivable type is returned succesfully, and no error is
// returned.
func (p *SimpleParser) Parse(name *amf0.String, r io.Reader) (Receivable, error) {
	str := string(*name)

	factory := p.typs[str]
	if factory == nil {
		return nil, fmt.Errorf(
			"rtmp/cmd/conn: unknown command name: %v", str)
	}

	v := factory()
	if err := encoding.Unmarshal(r, v); err != nil {
		return nil, err
	}

	return v, nil
}
