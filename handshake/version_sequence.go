package handshake

import (
	"crypto/rand"
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
	// S1 is the S1 packet that will be sent. It is generated in the
	// NewVersionSequence() pseudo-constructor.
	S1 *AckPacket
	// Supported is the supported version byte that this server can handle.
	Supported byte
}

var _ Sequence = new(VerisonSequence)

// NewVersionSequence instantiates and returns a pointer to a new instance of
// the VerisonSequence type. A call to this method also initializes a random
// payload into the S1 packet by reading from `rand.Reader`.
func NewVersionSequence() *VerisonSequence {
	v := &VerisonSequence{
		S1:        new(AckPacket),
		Supported: SupportedRTMPVersion,
	}

	rand.Read(v.S1.Payload[:])

	return v
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

// Write writes out the supported version as well as the S1 payload, returning
// any errors it encountered, if there were any.
func (v *VerisonSequence) Write(w io.Writer) error {
	if _, err := w.Write([]byte{v.Supported}); err != nil {
		return err
	}

	if err := v.S1.Write(w); err != nil {
		return err
	}

	return nil
}

// Next returns the AckSequence, which is the next logical sequence according to
// the RTMP specification. It initializes next with the same random challenge
// payload that was sent in S1, for comparison to C2.
func (v *VerisonSequence) Next() Sequence {
	return NewAckSequence(v.S1.Payload)
}
