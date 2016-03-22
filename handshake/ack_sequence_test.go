package handshake_test

import (
	"bytes"
	"testing"

	"github.com/WatchBeam/rtmp/handshake"
	"github.com/stretchr/testify/assert"
)

func TestAckSequenceConstruction(t *testing.T) {
	var payload [1528]byte
	a := handshake.NewAckSequence(payload)

	assert.IsType(t, new(handshake.AckSequence), a)
}

func TestAckSequenceReadsMatchingPayloads(t *testing.T) {
	c1 := &handshake.AckPacket{0, 0, RandomPayload()}
	c2 := &handshake.AckPacket{0, 0, RandomPayload()}

	buf := new(bytes.Buffer)
	c1.Write(buf) // Write C1
	c2.Write(buf) // Write C2

	a := handshake.NewAckSequence(c2.Payload) // Matching C2 Payload
	err := a.Read(buf)

	assert.Nil(t, err)
}

func TestAckSequencesErrsMismatchedPayloads(t *testing.T) {
	c1 := handshake.NewAckSequence(RandomPayload())
	c2 := handshake.NewAckSequence(RandomPayload())

	buf := new(bytes.Buffer)
	c1.Write(buf) // Write C1
	c2.Write(buf) // Write C2

	a := handshake.NewAckSequence(RandomPayload()) // Mismatched C2 Payload
	err := a.Read(buf)

	assert.Equal(t, handshake.ErrMismatchedPayload, err)
}

func TestAckSequenceWritesC1Payload(t *testing.T) {
	buf := new(bytes.Buffer)
	a := handshake.NewAckSequence(RandomPayload())
	a.C1.Payload = RandomPayload()

	a.Write(buf)

	assert.Equal(t, a.C1.Payload[:], buf.Bytes()[8:])
}

func TestAckSequenceHasNoNext(t *testing.T) {
	var payload [1528]byte
	a := handshake.NewAckSequence(payload)

	next := a.Next()

	assert.Nil(t, next)
}
