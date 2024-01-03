package data_test

import (
	"testing"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/todostreaming/rtmp/cmd/data"
	"github.com/stretchr/testify/assert"
)

func TestNewParserConstructsNewParsers(t *testing.T) {
	p := data.NewParser()

	assert.IsType(t, new(data.SimpleParser), p)
}

func TestParseParsesChunks(t *testing.T) {
	d, err := data.DefaultParser.Parse(&chunk.Chunk{
		Header: &chunk.Header{
			MessageHeader: chunk.MessageHeader{
				TypeId: 0x08,
			},
		},
		Data: []byte{0x0, 0x1},
	})

	assert.Nil(t, err)
	assert.Equal(t, byte(0x08), d.Id())
}

func TestParseErrsChunksWithMismatchedTypeIds(t *testing.T) {
	d, err := data.DefaultParser.Parse(&chunk.Chunk{
		Header: &chunk.Header{
			MessageHeader: chunk.MessageHeader{
				TypeId: 0x01,
			},
		},
		Data: []byte{0x0, 0x1},
	})

	assert.Equal(t, "rtmp/data: unknown type ID 1", err.Error())
	assert.Nil(t, d)
}

func TestParsePropogatesChunkReadingErrors(t *testing.T) {
	d, err := data.DefaultParser.Parse(&chunk.Chunk{
		Header: &chunk.Header{
			MessageHeader: chunk.MessageHeader{
				TypeId: 0x08,
			},
		},
		Data: []byte{},
	})

	assert.Nil(t, d)
	assert.Equal(t, data.ErrControlMissing, err)
}

func TestParseConstructsNewInstancesOfDatas(t *testing.T) {
	p := data.NewParser(
		func() data.Data { return &data.Audio{} },
	)

	d := p.New(0x08)

	assert.IsType(t, new(data.Audio), d)
}

func TestParseReturnsNilWhenNoMatchingDataExists(t *testing.T) {
	p := data.NewParser()

	d := p.New(0x08)

	assert.Nil(t, d)
}
