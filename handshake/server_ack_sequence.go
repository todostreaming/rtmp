package handshake

import (
	"bytes"
	"errors"
	"io"
)

var (
	// MismatchedChallengeErr is an error which is returned in a situtation
	// when the challenge sequence read from the client in C2 differs from
	// the challenge sequence in S1.
	MismatchedChallengeErr = errors.New("rtmp/handshake: mismatched challenge")
)

// ServerAckSequence is a type implementing the handshake.Sequence interface and
// is repsonseible for reading and verifying the C2 packet written by the client
// against the S1 packet from the server.
type ServerAckSequence struct {
	// S1 is the packet which C2 should acknowledge.
	S1 *AckPacket
}

var _ Sequence = new(ServerAckSequence)

// NewServerAckSequence returns a new *ServerAckSequence initialized with the
// given S1 packet.
func NewServerAckSequence(S1 *AckPacket) *ServerAckSequence {
	return &ServerAckSequence{S1}
}

// Read implements the Handshake.Read method by reading the C2 packet and
// comparing it to the stored S1 packet. If a read error occured while reading
// C2, then it will be returned. If the payloads were not equal, then
// MismatchedChallengeErr will be returned. Otherwise, in the successful case,
// a value of nil will be returned.
func (s *ServerAckSequence) Read(r io.Reader) error {
	c2 := new(AckPacket)
	if err := c2.Read(r); err != nil {
		return err
	}

	if !bytes.Equal(s.S1.Payload[:], c2.Payload[:]) {
		return MismatchedChallengeErr
	}

	return nil
}

// WriteTo implements the Sequence.WriteTo function. Since there is nothing to
// write, a value of nil is always returned here.
func (s *ServerAckSequence) WriteTo(w io.Writer) error { return nil }

// Next implements the Sequence.Next function. Since there is no next function,
// this function always returns nil.
func (s *ServerAckSequence) Next() Sequence { return nil }
