package data

import "github.com/todostreaming/rtmp/chunk"

// Type Parser is an interface which contains a single method, Parse. It is
// responsible for turning a *chunk.Chunk into an appropriate Data type, or
// returning an error, otherwise.
type Parser interface {
	// Parse is the single method placed on Parser, and it is responsible
	// for parsing a *chunk.Chunk into a piece of Data, or an error.
	Parse(*chunk.Chunk) (Data, error)
}
