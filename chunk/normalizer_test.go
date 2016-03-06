package chunk_test

import (
	"testing"

	"github.com/WatchBeam/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

func TestNewNormalizer(t *testing.T) {
	n := chunk.NewNormalizer()

	assert.IsType(t, &chunk.Normalizer{}, n)
}

func TestInitialChunkIsNil(t *testing.T) {
	n := chunk.NewNormalizer()
	c := n.Last()

	assert.Nil(t, c)
}

func TestSettingLastChunk(t *testing.T) {
	n := chunk.NewNormalizer()
	chunk := &chunk.Chunk{}

	n.SetLast(chunk)

	assert.Equal(t, chunk, n.Last())
}

func TestStoringHeaders(t *testing.T) {
	n := chunk.NewNormalizer()
	h := &chunk.Header{
		BasicHeader: chunk.BasicHeader{
			StreamId: 0x01,
		},
	}

	n.StoreHeader(h)

	assert.Equal(t, h, n.Header(0x01))
}

func TestNormalizingStoresChunks(t *testing.T) {
	n := chunk.NewNormalizer()
	n.Normalize(&chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader: chunk.BasicHeader{
				StreamId: 12,
			},
			MessageHeader: chunk.MessageHeader{
				StreamId: 3,
				Length:   8,
			},
		},
	})

	last := n.Last()
	header := n.Header(12)

	assert.Equal(t, uint32(3), last.Header.MessageHeader.StreamId)
	assert.Equal(t, uint32(8), last.Header.MessageHeader.Length)
	assert.Equal(t, uint32(3), header.MessageHeader.StreamId)
	assert.Equal(t, uint32(8), header.MessageHeader.Length)
}

func TestNormalizingSkippedWhenNotNecessary(t *testing.T) {
	n := chunk.NewNormalizer()
	n.SetLast(&chunk.Chunk{
		Header: &chunk.Header{
			MessageHeader: chunk.MessageHeader{
				StreamId: 3,
				Length:   8,
			},
		},
	})

	chunk := &chunk.Chunk{
		Header: &chunk.Header{
			MessageHeader: chunk.MessageHeader{
				TypeId: 0,
			},
		},
	}

	n.Normalize(chunk)

	assert.Equal(t, uint32(0), chunk.Header.MessageHeader.StreamId)
	assert.Equal(t, uint32(0), chunk.Header.MessageHeader.Length)
}

func TestNormalizingFillsInPartialHeadersWhenNecessary(t *testing.T) {
	n := chunk.NewNormalizer()
	n.SetLast(&chunk.Chunk{
		Header: &chunk.Header{
			MessageHeader: chunk.MessageHeader{
				StreamId: 3,
				Length:   8,
			},
		},
	})

	chunk := &chunk.Chunk{
		Header: &chunk.Header{
			MessageHeader: chunk.MessageHeader{
				TypeId: 2,
			},
		},
	}
	n.Normalize(chunk)

	assert.Equal(t, uint32(3), chunk.Header.MessageHeader.StreamId)
	assert.Equal(t, uint32(8), chunk.Header.MessageHeader.Length)
}

func TestNormalizingFillsMissingHeaders(t *testing.T) {
	n := chunk.NewNormalizer()
	n.StoreHeader(&chunk.Header{
		MessageHeader: chunk.MessageHeader{
			TypeId:    3,
			Timestamp: 4,
			Length:    5,
		},
	})

	chunk := &chunk.Chunk{
		Header: &chunk.Header{
			MessageHeader: chunk.MessageHeader{
				TypeId: 3,
			},
		},
	}
	n.Normalize(chunk)

	assert.Equal(t, byte(3), chunk.Header.MessageHeader.TypeId)
	assert.Equal(t, uint32(4), chunk.Header.MessageHeader.Timestamp)
	assert.Equal(t, uint32(5), chunk.Header.MessageHeader.Length)
}
