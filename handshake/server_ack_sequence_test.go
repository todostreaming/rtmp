package handshake_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	"github.com/todostreaming/rtmp/handshake"
	"github.com/stretchr/testify/assert"
)

func TestServerAckSequenceReadsMatchingC2(t *testing.T) {
	s1 := new(handshake.AckPacket)
	rand.Read(s1.Payload[:])

	buf := bytes.NewBuffer([]byte{})
	buf.Write([]byte{0x0, 0x0, 0x0, 0x0}) // Time1
	buf.Write([]byte{0x0, 0x0, 0x0, 0x0}) // Padding
	buf.Write(s1.Payload[:])              // Matching payload

	s := handshake.NewServerAckSequence(s1)

	err := s.Read(buf)

	assert.Nil(t, err)
}

func TestServerAckSequenceReportsMismatchedChallenges(t *testing.T) {
	s1 := new(handshake.AckPacket)
	rand.Read(s1.Payload[:])

	var mismatched [handshake.PayloadLen]byte
	rand.Read(mismatched[:])

	buf := bytes.NewBuffer([]byte{})
	buf.Write([]byte{0x0, 0x0, 0x0, 0x0}) // Time1
	buf.Write([]byte{0x0, 0x0, 0x0, 0x0}) // Padding
	buf.Write(mismatched[:])              // Mismatched payload

	s := handshake.NewServerAckSequence(s1)

	err := s.Read(buf)

	assert.Equal(t, handshake.MismatchedChallengeErr, err)
}

func TestServerAckSequenceReportsErroredReads(t *testing.T) {
	s := handshake.NewServerAckSequence(nil)

	err := s.Read(bytes.NewBuffer([]byte{
	// Empty buffer ~> io.EOF
	}))

	assert.Equal(t, io.EOF, err)
}

func TestServerAckSequenceDoesNotWrite(t *testing.T) {
	s := handshake.NewServerAckSequence(nil)
	buf := bytes.NewBuffer([]byte{})

	err := s.WriteTo(buf)

	assert.Nil(t, err)
	assert.Empty(t, buf.Bytes())
}

func TestServerAckSequenceDoesNotHaveSubsequentSequences(t *testing.T) {
	s := handshake.NewServerAckSequence(nil)

	assert.Nil(t, s.Next())
}
