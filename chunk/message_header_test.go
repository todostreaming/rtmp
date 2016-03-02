package chunk_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/WatchBeam/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

func TestMessageHeaderKnowsWhenMoreTimestampIsNeeded(t *testing.T) {
	h := &chunk.MessageHeader{
		Timestamp: 0xffffff,
	}

	assert.True(t, h.HasExtendedTimestamp())
}

type MessageHeaderTestCase struct {
	Buffer  []byte
	Control chunk.MessageHeader
}

func (c *MessageHeaderTestCase) Assert(t *testing.T) {
	h := chunk.MessageHeader{
		FormatId: c.Control.FormatId,
	}
	err := h.Read(bytes.NewBuffer(c.Buffer))

	assert.Nil(t, err)
	assert.Equal(t, c.Control, h)
}

func TestMessageHeaderTypeZeroReads(t *testing.T) {
	timestamp := uint32(123423)
	length := uint32(456152)
	typeId := byte(7)
	streamId := rand.Uint32()

	buf := []byte{
		byte(timestamp >> 16), byte(timestamp >> 8), byte(timestamp),
		byte(length >> 16), byte(length >> 8), byte(length),
		byte(typeId),
		byte(streamId >> 24), byte(streamId >> 16), byte(streamId >> 8), byte(streamId),
	}

	c := &MessageHeaderTestCase{
		Buffer: buf,
		Control: chunk.MessageHeader{
			FormatId:       0,
			Timestamp:      timestamp,
			TimestampDelta: false,
			Length:         length,
			TypeId:         typeId,
			StreamId:       streamId,
		},
	}

	c.Assert(t)
}

func TestMessageHeaderTypeOneReads(t *testing.T) {
	timestamp := uint32(123423)
	length := uint32(456152)
	typeId := byte(7)

	buf := []byte{
		byte(timestamp >> 16), byte(timestamp >> 8), byte(timestamp),
		byte(length >> 16), byte(length >> 8), byte(length),
		byte(typeId),
	}

	c := &MessageHeaderTestCase{
		Buffer: buf,
		Control: chunk.MessageHeader{
			FormatId:       1,
			Timestamp:      timestamp,
			TimestampDelta: true,
			Length:         length,
			TypeId:         typeId,
		},
	}

	c.Assert(t)
}

func TestMessageHeaderTypeTwoReads(t *testing.T) {
	timestamp := uint32(123423)
	buf := []byte{
		byte(timestamp >> 16), byte(timestamp >> 8), byte(timestamp),
	}

	c := &MessageHeaderTestCase{
		Buffer: buf,
		Control: chunk.MessageHeader{
			FormatId:       2,
			Timestamp:      timestamp,
			TimestampDelta: true,
		},
	}

	c.Assert(t)
}
