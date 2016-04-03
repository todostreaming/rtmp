package data_test

import (
	"testing"

	"github.com/WatchBeam/rtmp/cmd/data"
	"github.com/stretchr/testify/assert"
)

func TestAudioDataReturnsCorrectId(t *testing.T) {
	a := new(data.Audio)

	assert.Equal(t, data.AudioTypeId, a.Id())
}

func TestAudioDeterminesCorrectCodec(t *testing.T) {
	for _, c := range []struct {
		Control byte
		Codec   data.AudioCodec
	}{
		{0x00, data.UncompressedAudioCodec},
		{0x10, data.ADPCMAudioCodec},
		{0x20, data.MP3AudioCodec},
		{0x30, data.UncompressedLittleEndianAudioCodec},
		{0x40, data.Nellymoser16AudioCodec},
		{0x50, data.Nellymoser8AudioCodec},
		{0x60, data.NellymoserAudioCodec},
		{0x70, data.G711AAudioCodec},
		{0x80, data.G711UAudioCodec},
		{0x90, data.HE_AACAudioCodec},
		{0xa0, data.SPEEXAudioCodec},
	} {
		a := new(data.Audio)
		a.Control = c.Control

		assert.Equal(t, c.Codec, a.Codec())
	}
}

func TestAudioCanCalculateRate(t *testing.T) {
	a := new(data.Audio)
	a.Control = 0x0c

	assert.Equal(t, float32(5.5), a.Rate())
}

func TestAudioCanCalculateSize(t *testing.T) {
	a := new(data.Audio)
	a.Control = 0x02

	assert.Equal(t, 24, a.Size())
}

func TestAudioDeterminesCorrectType(t *testing.T) {
	for _, c := range []struct {
		Control byte
		Type    data.AudioType
	}{
		{0x00, data.MonoAudioType},
		{0x01, data.StereoAudioType},
	} {
		a := new(data.Audio)
		a.Control = c.Control

		assert.Equal(t, c.Type, a.Type())
	}
}
