package handshake_test

import (
	"bytes"
	"testing"

	"github.com/todostreaming/rtmp/handshake"
	"github.com/stretchr/testify/assert"
)

func TestItConstructsVersionSequences(t *testing.T) {
	v := handshake.NewVersionSequence()

	assert.IsType(t, new(handshake.VerisonSequence), v)
}

func TestItReadsSupportedVersionNumbers(t *testing.T) {
	v := handshake.NewVersionSequence()
	err := v.Read(bytes.NewBuffer([]byte{0x3}))

	assert.Nil(t, err)
}

func TestItRejectsUnsupportedVersionNumbers(t *testing.T) {
	v := handshake.NewVersionSequence()
	err := v.Read(bytes.NewBuffer([]byte{0x4}))

	assert.Equal(t, "rtmp/handshake: unsupported version 4", err.Error())
}

func TestItWritesVersion(t *testing.T) {
	buf := new(bytes.Buffer)
	v := handshake.NewVersionSequence()

	assert.Nil(t, v.WriteTo(buf))
	assert.Equal(t, []byte{0x3}, buf.Bytes()[:1])
}
