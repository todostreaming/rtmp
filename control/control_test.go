package control_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/todostreaming/rtmp/control"
	"github.com/stretchr/testify/assert"
)

func TestAbortControlReadWrite(t *testing.T) {
	n := rand.Uint32()
	tc := &ControlSequenceTestCase{
		Data: []byte{
			byte((n >> 24) & 0xff), byte((n >> 16) & 0xff),
			byte((n >> 8) & 0xff), byte((n >> 0) & 0xff),
		},
		Control: &control.AbortMessage{n},
	}

	tc.Assert(t)
}

func TestAcknowledgementControlReadWrite(t *testing.T) {
	n := rand.Uint32()
	tc := &ControlSequenceTestCase{
		Data: []byte{
			byte((n >> 24) & 0xff), byte((n >> 16) & 0xff),
			byte((n >> 8) & 0xff), byte((n >> 0) & 0xff),
		},
		Control: &control.Acknowledgement{n},
	}

	tc.Assert(t)
}

func TestSetChunkSizeInRangeReadWriteInRange(t *testing.T) {
	rand.Seed(time.Now().Unix())
	n := uint32(rand.Intn(0xffffff))

	tc := &ControlSequenceTestCase{
		Data: []byte{
			byte((n >> 24) & 0xff), byte((n >> 16) & 0xff),
			byte((n >> 8) & 0xff), byte((n >> 0) & 0xff),
		},
		Control: control.NewSetChunkSize(n),
	}

	tc.Assert(t)
}

func TestSetPeerBandwidth(t *testing.T) {
	n := rand.Uint32()
	l := control.LimitTypeSoft

	tc := &ControlSequenceTestCase{
		Data: []byte{
			byte((n >> 24) & 0xff), byte((n >> 16) & 0xff),
			byte((n >> 8) & 0xff), byte((n >> 0) & 0xff),
			byte(1),
		},
		Control: &control.SetPeerBandwidth{n, l},
	}

	tc.Assert(t)
}

func TestWindowAckSizeReadWrite(t *testing.T) {
	n := rand.Uint32()
	tc := &ControlSequenceTestCase{
		Data: []byte{
			byte((n >> 24) & 0xff), byte((n >> 16) & 0xff),
			byte((n >> 8) & 0xff), byte((n >> 0) & 0xff),
		},
		Control: &control.WindowAckSize{n},
	}

	tc.Assert(t)
}

type ControlSequenceTestCase struct {
	Data    []byte
	Control control.Control
}

func (c *ControlSequenceTestCase) Assert(t *testing.T) {
	c.assertRead(t)
}

func (c *ControlSequenceTestCase) assertRead(t *testing.T) {
	typ := reflect.TypeOf(c.Control).Elem()
	v := reflect.New(typ)

	ctrl := v.Interface().(control.Control)

	err := ctrl.Read(bytes.NewBuffer(c.Data))

	assert.Nil(t, err)
	assert.Equal(t, c.Control, ctrl)
}

func (c *ControlSequenceTestCase) assertWrite(t *testing.T) {
	buf := new(bytes.Buffer)

	err := c.Control.Write(buf)

	assert.Nil(t, err)
	assert.True(t, bytes.Equal(c.Data, buf.Bytes()), fmt.Sprintf(
		"control: payloads should be equal (%v, %v)", c.Data, buf.Bytes()))
}
