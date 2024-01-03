package control

import "github.com/todostreaming/rtmp/chunk"

// Type Parser represents an interface capable of turning RTMP chunks into
// Control sequences.
type Parser interface {
	// Parse maps a RTMP chunk to a control sequence, and perhaps an error,
	// if one was encountered during this transaction.
	Parse(*chunk.Chunk) (Control, error)
}
