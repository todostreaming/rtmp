package chunk_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

type MessageHeaderTestCase struct {
	IdealBuffer []byte
	Header      *chunk.MessageHeader
}

func (m *MessageHeaderTestCase) Assert(t *testing.T) {
	m.assertRead(t)
	m.assertWrite(t)
}

func (m *MessageHeaderTestCase) assertRead(t *testing.T) {
	h := &chunk.MessageHeader{
		FormatId: m.Header.FormatId,
	}
	err := h.Read(bytes.NewBuffer(m.IdealBuffer))

	assert.Nil(t, err, "message header: read err should not exist")

	assert.Equal(t, m.Header.FormatId, h.FormatId)
	assert.Equal(t, m.Header.Timestamp, h.Timestamp)
	assert.Equal(t, m.Header.TimestampDelta, h.TimestampDelta)
	assert.Equal(t, m.Header.Length, h.Length)
	assert.Equal(t, m.Header.TypeId, h.TypeId)
	assert.Equal(t, m.Header.StreamId, h.StreamId)
}

func (m *MessageHeaderTestCase) assertWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	err := m.Header.Write(buf)

	assert.Nil(t, err)
	assert.True(t, bytes.Equal(m.IdealBuffer, buf.Bytes()), fmt.Sprintf(
		"message header: slice should be equal, wasn't (%v,%v)",
		m.IdealBuffer, buf.Bytes()))
}

func TestMessageHeaderKnowsWhenMoreTimestampIsNeeded(t *testing.T) {
	h := &chunk.MessageHeader{
		Timestamp: 0xffffff,
	}

	assert.True(t, h.HasExtendedTimestamp())
}

func TestMessageHeaderTypeZeroReadWrite(t *testing.T) {
	timestamp := uint32(123423)
	length := uint32(456152)
	typeId := byte(7)
	streamId := uint32(12352)

	buf := []byte{
		byte(timestamp >> 16), byte(timestamp >> 8), byte(timestamp),
		byte(length >> 16), byte(length >> 8), byte(length),
		byte(typeId),
		byte(streamId >> 0), byte(streamId >> 8), byte(streamId >> 16),
		byte(streamId >> 24),
	}

	c := &MessageHeaderTestCase{
		IdealBuffer: buf,
		Header: &chunk.MessageHeader{
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

func TestMessageHeaderTypeOneReadWrite(t *testing.T) {
	timestamp := uint32(123423)
	length := uint32(456152)
	typeId := byte(7)

	buf := []byte{
		byte(timestamp >> 16), byte(timestamp >> 8), byte(timestamp),
		byte(length >> 16), byte(length >> 8), byte(length),
		byte(typeId),
	}

	c := &MessageHeaderTestCase{
		IdealBuffer: buf,
		Header: &chunk.MessageHeader{
			FormatId:       1,
			Timestamp:      timestamp,
			TimestampDelta: true,
			Length:         length,
			TypeId:         typeId,
		},
	}

	c.Assert(t)
}

func TestMessageHeaderTypeTwoReadWrite(t *testing.T) {
	timestamp := uint32(123423)
	buf := []byte{
		byte(timestamp >> 16), byte(timestamp >> 8), byte(timestamp),
	}

	c := &MessageHeaderTestCase{
		IdealBuffer: buf,
		Header: &chunk.MessageHeader{
			FormatId:       2,
			Timestamp:      timestamp,
			TimestampDelta: true,
		},
	}

	c.Assert(t)
}
