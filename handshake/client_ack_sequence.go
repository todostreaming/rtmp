package handshake

import (
	"crypto/rand"
	"io"
)

// ClientAckSequence is a handshake sequence that verifies the challenge
// sequence sent by the client during a RTMP handshake. It is responsible for
// reading and responding to the C1 packet (with the C2 packet), and sending the
// server challenge in S1.
type ClientAckSequence struct {
	C1 *AckPacket
	S1 *AckPacket
}

var _ Sequence = new(ClientAckSequence)

// NewClientAckSequence initializes and returns a new *ClientAckSequence
// initialized with an empty C1 packet and a new S1 packet, initialized with the
// result of a "rand.Read" into its Payload header.
func NewClientAckSequence() *ClientAckSequence {
	c := &ClientAckSequence{
		C1: new(AckPacket),
		S1: new(AckPacket),
	}

	rand.Read(c.S1.Payload[:])

	return c
}

// Read implements the Sequence.Read function. It reads the C1 packet and
// returns any read error, if there was one. Otherwise, a value of "nil" is
// returned instead.
func (c *ClientAckSequence) Read(r io.Reader) error {
	if err := c.C1.Read(r); err != nil {
		return err
	}

	return nil
}

// WriteTo implements the Sequence.WriteTo function. It writes the S1 packet
// first (returning any errors if there is one), and then writes the S2 packet
// with the same data as was sent in the C1 packet (returning any error that was
// encountered).
//
// A successful call to Write constitutes a value of `nil` being returned.
func (c *ClientAckSequence) WriteTo(w io.Writer) error {
	if err := c.S1.Write(w); err != nil {
		return err
	}

	s2 := &AckPacket{
		Time1:   c.C1.Time1,
		Payload: c.C1.Payload,
	}

	if err := s2.Write(w); err != nil {
		return err
	}

	return nil
}

// Nex implements the Sequence.Next function.
func (c *ClientAckSequence) Next() Sequence {
	return NewServerAckSequence(c.S1)
}
