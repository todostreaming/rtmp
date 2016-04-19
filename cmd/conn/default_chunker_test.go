package conn_test

import (
	"testing"

	"github.com/WatchBeam/amf0"
	"github.com/WatchBeam/amf0/encoding"
	"github.com/WatchBeam/rtmp/cmd/conn"
	"github.com/stretchr/testify/assert"
)

func TestNewChunkerReturnsNewChunkers(t *testing.T) {
	c := conn.NewChunker(0)

	assert.IsType(t, &conn.DefaultChunker{}, c)
}

func TestChunkersMarshalChunks(t *testing.T) {
	crsp := &conn.ConnectResponse{
		TransactionId: 12,
		Properties:    amf0.Object{amf0.NewPaired()},
		Information:   amf0.Object{amf0.NewPaired()},
	}

	marshalled, _ := encoding.Marshal(crsp)

	c, err := conn.NewChunker(13).Chunk(crsp)

	assert.Nil(t, err)
	assert.Equal(t, conn.ChunkStreamId, c.Header.BasicHeader.StreamId)
	assert.Equal(t, byte(0x14), c.Header.MessageHeader.TypeId)
	assert.Equal(t, uint32(len(marshalled)), c.Header.MessageHeader.Length)
	assert.Equal(t, marshalled, c.Data)
}
