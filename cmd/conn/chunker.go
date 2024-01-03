package conn

import "github.com/todostreaming/rtmp/chunk"

// Chunker is a functional interface designed to turn ConnSendable types into
// RTMP `*chunk.Chunk`s.
type Chunker interface {
	// Chunk marshals a Marshallable into the *chunk.Chunk type, using the
	// protocol defined in the RTMP specification.
	Chunk(m Marshallable) (*chunk.Chunk, error)
}
