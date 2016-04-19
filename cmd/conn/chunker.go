package conn

import "github.com/WatchBeam/rtmp/chunk"

// Chunker is a functional interface designed to turn ConnSendable types into
// RTMP `*chunk.Chunk`s.
type Chunker interface {
	// Chunk marshals a ConnSendable into the *chunk.Chunk type, using the
	// protocol defined in the RTMP specification.
	Chunk(s Sendable) (*chunk.Chunk, error)
}
