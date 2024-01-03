package handshake_test

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"testing"

	"github.com/todostreaming/rtmp/handshake"
	"github.com/stretchr/testify/assert"
)

func RandomPayload() [handshake.PayloadLen]byte {
	var payload [handshake.PayloadLen]byte

	rand.Read(payload[:])

	return payload
}

func TestItReadsValuesFromReaders(t *testing.T) {
	a := new(handshake.AckPacket)
	buf := new(bytes.Buffer)

	payload := RandomPayload()

	binary.Write(buf, binary.BigEndian, uint32(11))
	binary.Write(buf, binary.BigEndian, uint32(22))
	buf.Write(payload[:])

	err := a.Read(buf)

	assert.Nil(t, err)
	assert.Equal(t, uint32(11), a.Time1)
	assert.Equal(t, uint32(22), a.Time2)

	assert.Equal(t, payload, a.Payload)
}

func TestItWritesValuesToWriters(t *testing.T) {
	buf := new(bytes.Buffer)
	a := &handshake.AckPacket{
		Time1:   33,
		Time2:   44,
		Payload: RandomPayload(),
	}

	assert.Nil(t, a.Write(buf))

	assert.Equal(t, []byte{0, 0, 0, 33}, buf.Bytes()[:4])
	assert.Equal(t, []byte{0, 0, 0, 44}, buf.Bytes()[4:8])
	assert.Equal(t, a.Payload[:], buf.Bytes()[8:])
}
