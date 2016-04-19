package stream

import "io"

// Parser is a functional interface responsible for turning commands sent over
// a stream of bytes into a Command, or an error.
type Parser interface {
	// Parse parses the next single command read off of the given io.Reader,
	// and returns the a new instance of the assosciated Command type,
	// initialized with all of the serialized data read from `r`.
	//
	// If the bytes read off of the io.Reader were not able to be parsed for
	// any reason, then an appropriate error is returned instead.
	Parse(r io.Reader) (Command, error)
}
