package handshake_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/todostreaming/rtmp/handshake"
	"github.com/stretchr/testify/assert"
)

func TestItConstructsNewHandshakers(t *testing.T) {
	h := handshake.With(&handshake.Param{
		Conn: new(bytes.Buffer),
	})

	assert.IsType(t, new(handshake.Handshaker), h)
}

func TestItCallsHandshakeReadWrite(t *testing.T) {
	conn := new(bytes.Buffer)

	initial := new(MockSequence)
	initial.On("Read", conn).Return(nil).Once()
	initial.On("WriteTo", conn).Return(nil).Once()
	initial.On("Next").Return(nil).Once()

	h := handshake.With(&handshake.Param{
		Conn:    conn,
		Initial: initial,
	})

	assert.Nil(t, h.Handshake())
	initial.AssertExpectations(t)
}

func TestItAbortsOnReadErrors(t *testing.T) {
	conn := new(bytes.Buffer)

	initial := new(MockSequence)
	initial.On("Read", conn).Return(errors.New("foo")).Once()

	h := handshake.With(&handshake.Param{
		Conn:    conn,
		Initial: initial,
	})

	assert.Equal(t, "foo", h.Handshake().Error())
	initial.AssertExpectations(t)
}

func TestItAbortsOnWriteErrors(t *testing.T) {
	conn := new(bytes.Buffer)

	initial := new(MockSequence)
	initial.On("Read", conn).Return(nil).Once()
	initial.On("WriteTo", conn).Return(errors.New("foo")).Once()

	h := handshake.With(&handshake.Param{
		Conn:    conn,
		Initial: initial,
	})

	assert.Equal(t, "foo", h.Handshake().Error())
	initial.AssertExpectations(t)
}

func TestItCyclesToNextSequence(t *testing.T) {
	conn := new(bytes.Buffer)

	next := new(MockSequence)
	next.On("Read", conn).Return(nil).Once()
	next.On("WriteTo", conn).Return(nil).Once()
	next.On("Next").Return(nil).Once()

	initial := new(MockSequence)
	initial.On("Read", conn).Return(nil).Once()
	initial.On("WriteTo", conn).Return(nil).Once()
	initial.On("Next").Return(next).Once()

	h := handshake.With(&handshake.Param{
		Conn:    conn,
		Initial: initial,
	})

	assert.Nil(t, h.Handshake())

	initial.AssertExpectations(t)
	next.AssertExpectations(t)
}
