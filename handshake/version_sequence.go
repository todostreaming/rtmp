package handshake

import (
	"fmt"
	"io"
)

const (
	// SupportedRTMPVersion is the default RMTP version tag this is
	// supported out in the wild.
	SupportedRTMPVersion byte = 3
)

// VerisonSequence is the sequence responsible for exchanging and acknowledging
// the RTMP versions across the network.
type VerisonSequence struct {
	// Supported is the supported version byte that this server can handle.
	Supported byte
}

var _ Sequence = new(VerisonSequence)

// NewVersionSequence instantiates and returns a pointer to a new instance of
// the VerisonSequence type.
func NewVersionSequence() *VerisonSequence {
	return &VerisonSequence{
		Supported: SupportedRTMPVersion,
	}
}

// Read reads the version byte off of the io.Reader, returning an error if
// either there was an error reading, or a version mismatch. Otherwise, a value
// of nil was retuend.
func (v *VerisonSequence) Read(r io.Reader) error {
	var b [1]byte
	if _, err := r.Read(b[:]); err != nil {
		return err
	}

	if b[0] != v.Supported {
		return fmt.Errorf(
			"rtmp/handshake: unsupported version %v", b[0])
	}

	return nil
}

// Write writes out the version that this RTMP server supports. If any write
// error was encountered, it will be returned immediately. Otherwise, a value of
// nil is returned.
func (v *VerisonSequence) Write(w io.Writer) error {
	if _, err := w.Write([]byte{v.Supported}); err != nil {
		return err
	}

	return nil
}

// Next returns the ClientAckSequence, which is the next step in the RTMP
// handshake, according to the specification.
func (v *VerisonSequence) Next() Sequence {
	return NewClientAckSequence()
}
