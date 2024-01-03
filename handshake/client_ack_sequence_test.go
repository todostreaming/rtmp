package handshake_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	"github.com/todostreaming/rtmp/handshake"
	"github.com/stretchr/testify/assert"
)

func TestItReadsC1Packet(t *testing.T) {
	payload := payload()

	buf := bytes.NewBuffer([]byte{})
	buf.Write([]byte{0x0, 0x0, 0x0, 0x1}) // Time 1
	buf.Write([]byte{0x0, 0x0, 0x0, 0x0}) // Padding
	buf.Write(payload[:])

	c := handshake.NewClientAckSequence()
	err := c.Read(buf)

	assert.Nil(t, err)
	assert.Equal(t, &handshake.AckPacket{
		Time1:   1,
		Time2:   0,
		Payload: payload,
	}, c.C1)
}

func TestItErorrsOnBadC1Packets(t *testing.T) {
	buf := bytes.NewBuffer([]byte{
	// Empty C1 ~> io.EOF
	})

	c := handshake.NewClientAckSequence()
	err := c.Read(buf)

	assert.Equal(t, io.EOF, err)
}

func TestItWritesS1AndMatchingS2(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})

	c := handshake.NewClientAckSequence()
	c.C1 = &handshake.AckPacket{
		Payload: payload(),
	}

	start := 4 + 4 + 1528 + 4 + 4
	end := start + 1528

	err := c.WriteTo(buf)

	assert.Nil(t, err)
	assert.Len(t, buf.Bytes(), 2*(4+4+1528))
	assert.Equal(t, c.C1.Payload[:], buf.Bytes()[start:end])
}

func payload() [handshake.PayloadLen]byte {
	var b [handshake.PayloadLen]byte

	rand.Read(b[:])

	return b
}
