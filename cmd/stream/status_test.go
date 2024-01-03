package stream_test

import (
	"testing"

	"github.com/todostreaming/amf0"
	"github.com/todostreaming/rtmp/chunk"
	"github.com/todostreaming/rtmp/cmd/stream"
	"github.com/stretchr/testify/assert"
)

var (
	ValidOnStatusHeader = []byte{
		0x02, 0x00, 0x08, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75,
		0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x05,
	}
)

func TestOnStatusCommandHeaderIsValid(t *testing.T) {
	assert.Equal(t, ValidOnStatusHeader, stream.OnStatusCommandHeader)
}

func TestNewStatusMakesNewStatuses(t *testing.T) {
	st := stream.NewStatus()

	assert.IsType(t, new(stream.Status), st)
}

func TestDataMarshalsTheStatusesData(t *testing.T) {
	st := stream.NewStatus()
	st.Arguments.Add("foo", amf0.NewString("bar"))

	data, err := st.Data()

	assert.Nil(t, err)
	assert.Equal(t, []byte{
		0x03, 0x00, 0x03, 0x66, 0x6f, 0x6f, 0x02, 0x00, 0x03, 0x62,
		0x61, 0x72, 0x00, 0x00, 0x09,
	}, data)
}

func TestAsChunkMarshalsTheStatusToChunks(t *testing.T) {
	st := stream.NewStatus()
	st.Arguments.Add("foo", amf0.NewString("bar"))

	c, err := st.AsChunk()

	expected := []byte{
		0x02, 0x00, 0x08, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75,
		0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x05, 0x03, 0x00, 0x03, 0x66, 0x6f, 0x6f, 0x02, 0x00, 0x03,
		0x62, 0x61, 0x72, 0x00, 0x00, 0x09,
	}

	assert.Nil(t, err)
	assert.Equal(t, chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader: chunk.BasicHeader{
				StreamId: 5,
			},
			MessageHeader: chunk.MessageHeader{
				Length:   uint32(len(expected)),
				TypeId:   0x14,
				StreamId: 1,
			},
		},
		Data: expected,
	}, *c)
}
