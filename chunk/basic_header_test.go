package chunk_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

type BasicHeaderTestCase struct {
	IdealBuffer []byte
	Header      *chunk.BasicHeader
}

func (b *BasicHeaderTestCase) Assert(t *testing.T) {
	b.assertRead(t)
	b.assertWrite(t)
}

func (b *BasicHeaderTestCase) assertRead(t *testing.T) {
	h := new(chunk.BasicHeader)
	h.Read(bytes.NewBuffer(b.IdealBuffer))

	assert.Equal(t, b.Header.FormatId, h.FormatId,
		"chunk: basic header format ID should be equal, wasn't")
	assert.Equal(t, b.Header.StreamId, h.StreamId,
		"chunk: basic header stream ID should be equal, wasn't")
}

func (b *BasicHeaderTestCase) assertWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	err := b.Header.Write(buf)

	assert.Nil(t, err, "chunk: writing should return no error, did")
	assert.True(t, bytes.Equal(b.IdealBuffer, buf.Bytes()), fmt.Sprintf(
		"chunk: basic header output should be equal, wasn't (%v,%v)",
		b.IdealBuffer, buf.Bytes()))
}

func TestBasicHeaderReadWrite(t *testing.T) {
	for _, c := range []BasicHeaderTestCase{
		{
			IdealBuffer: []byte{
				(0x02 << 6) | 23,
			},
			Header: &chunk.BasicHeader{
				FormatId: 2,
				StreamId: 23,
			},
		}, {
			IdealBuffer: []byte{
				(0x02 << 6) | 0, 23,
			},
			Header: &chunk.BasicHeader{
				FormatId: 2,
				StreamId: 23 + 64,
			},
		}, {
			IdealBuffer: []byte{
				(0x02 << 6) | 63, byte(330 >> 8),
				byte(uint32(330>>0) & 0xff),
			},
			Header: &chunk.BasicHeader{
				FormatId: 2,
				StreamId: 330 + 64,
			},
		},
	} {
		c.Assert(t)
	}
}
