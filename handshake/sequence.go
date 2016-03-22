package handshake

import "io"

// Sequence represents a single step in the entire RTMP handshake operation.
// Each sequence is assumed to have two components, the read side, and the write
// side. A sequence is also assumed to knows about which sequence comes after
// it, thus making the Sequence struct a sort of "pseudo-linked-list".
type Sequence interface {
	// Read preforms the read-operation associated with this Sequence.
	// Since this is an "incoming" sort of operation, it makes sense that
	// any persistent state on a implementing Sequence type would be
	// initialized here.
	//
	// If an error is encountered during this operation, then it should be
	// returned here, as well.
	Read(r io.Reader) error

	// Write preforms the write-operation assosciatd with this Sequence. It
	// writes any data to the connecting client (represented by the `w
	// io.Writer` parameter) necessary to complete this current handshake
	// Sequence.
	//
	// If an error was encountered during this operation, then it should be
	// returned here, as well.
	Write(w io.Writer) error

	// Next returns the Sequence that should follow this one. If no sequence
	// follows this one, then a value of nil should be returned, instead.
	Next() Sequence
}
