package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataReadsWithControlAndBody(t *testing.T) {
	d := new(data)

	err := d.Read([]byte{0x0, 0x1, 0x2, 0x3, 0x4})

	assert.Nil(t, err)
	assert.Equal(t, byte(0x0), d.Control)
	assert.Equal(t, []byte{0x1, 0x2, 0x3, 0x4}, d.Payload)
}

func TestDataReadsWithEmptyBody(t *testing.T) {
	d := new(data)

	err := d.Read([]byte{0x0})

	assert.Nil(t, err)
	assert.Equal(t, byte(0x0), d.Control)
	assert.Empty(t, d.Payload)
}

func TestDataDoesNotReadWhenMissingControl(t *testing.T) {
	d := new(data)
	err := d.Read([]byte{})

	assert.Equal(t, ErrControlMissing, err)
}
