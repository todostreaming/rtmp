package handshake

import "io"

// Handshaker preforms the RTMP handshake with the owned io.ReadWriter by
// cycling through a linked-list of sequences.
//
// Each sequence represents a set of operations to preform during the RTMP
// handshake. Additionally, each Sequence knows what sequence comes after it,
// and can instantiate any state on it when the next Sequence is asked for.
//
// A handshaker has an initial Sequence, which is, of course, the starting point
// for the complete Handshake operation.
type Handshaker struct {
	// rw is the io.ReadWriter that the Handshaker will handshake with.
	rw io.ReadWriter

	// initial represents the initial Sequence to begin the entire handshake
	// operation with.
	current Sequence
}

// Param wraps each argument passed to the constructor `func With`.
type Param struct {
	// Conn is the connection which the Handshaker will read and write to.
	// This parameter is required.
	Conn io.ReadWriter
	// Initial is the starting sequence. If not specified, RTMP will default
	// to the VersionSequence type, which is the initial sequence as
	// according to the RTMP specification.
	Initial Sequence
}

// With returns a new Handshaker initialized with the given Param.
func With(p *Param) *Handshaker {
	h := &Handshaker{
		rw: p.Conn,
	}

	if p.Initial != nil {
		h.current = p.Initial
	} else {
		h.current = NewVersionSequence()
	}

	return h
}

// Handshake preforms the handshake operation by cycling through each handshake
// sequence in the pseudo-linked-list of handshake sequences. Over each
// handshake sequence, it preforms the read side, and then the write side right
// afterwords. If either of these operations returns an error, then the entire
// handshake operation is halted, and the error is returned.
//
// If nil is returned, then the handshake completed succesfully without error,
// and the `current` handshake.Sequence will also be nil.
func (h *Handshaker) Handshake() error {
	for ; h.current != nil; h.current = h.current.Next() {
		if err := h.current.Read(h.rw); err != nil {
			return err
		}

		if err := h.current.Write(h.rw); err != nil {
			return err
		}
	}

	return nil
}
