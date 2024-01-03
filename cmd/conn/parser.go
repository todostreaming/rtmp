package conn

import (
	"io"

	"github.com/todostreaming/amf0"
)

// Parser is a functional interface responsible for parsing a body of data out
// of an io.Reader into a Receivable identified by the given `name`.
type Parser interface {
	// Parse parses a Receivable identified by `name` and with data living
	// on the io.Reader `r` into a Receivable, or an error, if the data was
	// incorrect, or not parse-able.
	Parse(name *amf0.String, r io.Reader) (Receivable, error)
}
