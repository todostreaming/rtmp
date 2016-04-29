package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAudioDataReturnsCorrectId(t *testing.T) {
	a := new(Audio)

	assert.Equal(t, AudioTypeId, a.Id())
}

func TestAudioDeterminesCorrectCodec(t *testing.T) {
	for _, c := range []struct {
		Control byte
		Codec   AudioCodec
	}{
		{0x00, UncompressedAudioCodec},
		{0x10, ADPCMAudioCodec},
		{0x20, MP3AudioCodec},
		{0x30, UncompressedLittleEndianAudioCodec},
		{0x40, Nellymoser16AudioCodec},
		{0x50, Nellymoser8AudioCodec},
		{0x60, NellymoserAudioCodec},
		{0x70, G711AAudioCodec},
		{0x80, G711UAudioCodec},
		{0x90, HE_AACAudioCodec},
		{0xa0, SPEEXAudioCodec},
	} {
		a := new(Audio)
		a.data.data = []byte{c.Control}

		assert.Equal(t, c.Codec, a.Codec())
	}
}

func TestAudioCanCalculateRate(t *testing.T) {
	a := new(Audio)
	a.data.data = []byte{0x0c}

	assert.Equal(t, float32(5.5), a.Rate())
}

func TestAudioCanCalculateSize(t *testing.T) {
	a := new(Audio)
	a.data.data = []byte{0x02}

	assert.Equal(t, 24, a.Size())
}

func TestAudioDeterminesCorrectType(t *testing.T) {
	for _, c := range []struct {
		Control byte
		Type    AudioType
	}{
		{0x00, MonoAudioType},
		{0x01, StereoAudioType},
	} {
		a := new(Audio)
		a.data.data = []byte{c.Control}

		assert.Equal(t, c.Type, a.Type())
	}
}
