package control_test

import (
	"bytes"
	"testing"

	"github.com/todostreaming/rtmp/control"
	"github.com/stretchr/testify/assert"
)

func TestEventReadsFromBuffer(t *testing.T) {
	e := new(control.Event)

	err := e.Read(bytes.NewReader([]byte{
		0x00, 0x03, 0x01, 0x02, 0x03,
	}))

	assert.Nil(t, err)
	assert.Equal(t, control.SetBufferLength, e.Type)
	assert.Equal(t, []byte{1, 2, 3}, e.Body)
}

func TestEventWritesToBuffer(t *testing.T) {
	buf := new(bytes.Buffer)
	e := &control.Event{
		Type: control.SetBufferLength,
		Body: []byte{4, 5, 6},
	}

	err := e.Write(buf)

	assert.Nil(t, err)
	assert.Equal(t, []byte{0x00, 0x03}, buf.Bytes()[:2])
	assert.Equal(t, []byte{0x04, 0x05, 0x06}, buf.Bytes()[2:])
}
