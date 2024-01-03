package control_test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/todostreaming/rtmp/control"
	"github.com/stretchr/testify/assert"
)

func TestParserConstruction(t *testing.T) {
	p := control.NewParser()

	assert.IsType(t, &control.DefaultParser{}, p)
}

func TestTypeLookupForValidIds(t *testing.T) {
	p := control.NewParser()
	expected := reflect.TypeOf(control.SetChunkSize{})

	typ := p.TypeFor(1)

	assert.True(t, typ.AssignableTo(expected))
}

func TestTypeLookupForInvalidIds(t *testing.T) {
	p := control.NewParser()

	typ := p.TypeFor(30)

	assert.Nil(t, typ)
}

func TestParsingForInvalidChunks(t *testing.T) {
	p := control.NewParser()

	ctrl, err := p.Parse(&chunk.Chunk{
		Header: &chunk.Header{
			MessageHeader: chunk.MessageHeader{TypeId: 30},
		},
	})

	assert.Nil(t, ctrl)
	assert.Equal(t, "control: unknown control type (30)", err.Error())
}

func TestParsingForValidChunks(t *testing.T) {
	n := rand.Uint32()
	p := control.NewParser()

	ctrl, err := p.Parse(&chunk.Chunk{
		Header: &chunk.Header{
			MessageHeader: chunk.MessageHeader{TypeId: 3},
		},
		Data: []byte{
			byte((n >> 24) & 0xff), byte((n >> 16) & 0xff),
			byte((n >> 8) & 0xff), byte((n >> 0) & 0xff),
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, &control.Acknowledgement{n}, ctrl)
}
