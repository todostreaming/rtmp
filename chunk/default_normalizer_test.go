package chunk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNormalizer(t *testing.T) {
	n := NewNormalizer()

	assert.IsType(t, &DefaultNormalizer{}, n)
}

func TestNormalizingSkippedWhenNotNecessary(t *testing.T) {
	n := NewNormalizer()
	n.Normalize(&Header{
		MessageHeader: MessageHeader{
			StreamId: 3,
			Length:   8,
		},
	})

	h := &Header{
		BasicHeader: BasicHeader{FormatId: 0},
	}

	h = n.Normalize(h)

	assert.Equal(t, uint32(0), h.MessageHeader.StreamId)
	assert.Equal(t, uint32(0), h.MessageHeader.Length)
}

func TestNormalizingFillsInPartialHeadersWhenNecessary(t *testing.T) {
	n := NewNormalizer()
	n.Normalize(&Header{
		MessageHeader: MessageHeader{
			StreamId: 3,
			Length:   8,
		},
	})

	h := &Header{
		BasicHeader: BasicHeader{FormatId: 2},
		MessageHeader: MessageHeader{
			StreamId: 4,
		},
	}

	h = n.Normalize(h)

	assert.Equal(t, uint32(8), h.MessageHeader.Length)
}

func TestNormalizingFillsMissingHeaders(t *testing.T) {
	n := NewNormalizer()
	n.Normalize(&Header{
		BasicHeader: BasicHeader{FormatId: 0},
		MessageHeader: MessageHeader{
			Timestamp: 4,
			Length:    5,
		},
	})

	h := &Header{
		BasicHeader: BasicHeader{FormatId: 3},
	}

	h = n.Normalize(h)

	assert.Equal(t, uint32(4), h.MessageHeader.Timestamp)
	assert.Equal(t, uint32(5), h.MessageHeader.Length)
}
