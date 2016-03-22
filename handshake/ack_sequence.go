package handshake

import (
	"bytes"
	"errors"
	"io"
)

var (
	// ErrMismatchedPayload is returned when the payload read during
	// AckSequence.Read is not equal to the S1Payload that was expected.
	ErrMismatchedPayload error = errors.New(
		"rtmp/handshake: mismatched payload")
)

type AckSequence struct {
	// S1Payload is the payload sent during S1.
	S1Payload [1528]byte

	// C1 and C2 are read during the Read operation and represent the first
	// and second packets sent back from the client during the Ack sequence.
	C1, C2 *AckPacket
}

var _ Sequence = new(AckSequence)

// NewAckSequence instantiates and returns a pointer to a new instance of the
// AckSequence type. The returned AckSequence is initialized with the S1Payload,
// as well as zero-value instances for C1 and C2.
func NewAckSequence(s1Payload [1528]byte) *AckSequence {
	return &AckSequence{
		S1Payload: s1Payload,
		C1:        new(AckPacket),
		C2:        new(AckPacket),
	}
}

// Read	reads both C1 and C2, as well as preform the comparison of S1 and C2. If
// an error occurs during the reads, or a mismatched payload was encountered the
// appropriate error will be returned immediately, and the execution will halt.
func (a *AckSequence) Read(r io.Reader) error {
	if err := a.C1.Read(r); err != nil {
		return err
	}

	c2 := new(AckPacket)
	if err := c2.Read(r); err != nil {
		return err
	}

	if !bytes.Equal(a.S1Payload[:], c2.Payload[:]) {
		return ErrMismatchedPayload
	}

	return nil
}

// Write writes the matching S2 packet which contains the payload found in S1.
// If an error was encountered during the write operation on the S2 packet, then
// it will be returned immediately, otherwise, a value of nil will be returned
// instead.
func (a *AckSequence) Write(w io.Writer) error {
	// TODO(taylor): Time1, Time2 should be the time that C1 was received,
	// but most implementations don't seem to care about this.
	s2 := &AckPacket{
		Payload: a.C1.Payload,
	}

	if err := s2.Write(w); err != nil {
		return err
	}

	return nil
}

// Next returns the sequence to run after the AckSequence. Per the RTMP
// specification, no sequence is to run after the AckSequence, therefore a value
// of nil is returned.
func (a *AckSequence) Next() Sequence { return nil }
