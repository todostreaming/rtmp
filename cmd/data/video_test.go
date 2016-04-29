package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideoReturnsCorrectId(t *testing.T) {
	v := new(Video)

	assert.Equal(t, VideoTypeId, v.Id())
}

func TestvideoProducesCorrectCodecs(t *testing.T) {
	for _, c := range []struct {
		Control    byte
		VideoCodec VideoCodec
	}{
		{0x00, SorensenH263VideoCodec},
		{0x01, ScreenVideoVideoCodec},
		{0x02, On2VP6VideoCodec},
		{0x03, On2VP6AlphaVideoCodec},
		{0x04, ScreenVideo2VideoCodec},
		{0x05, H264VideoCodec},
	} {
		d := new(Video)
		d.data.data = []byte{c.Control}

		assert.Equal(t, c.VideoCodec, d.Codec())
	}
}

func TestVideoProducesCorrectTypes(t *testing.T) {
	for _, c := range []struct {
		Control   byte
		VideoType VideoType
	}{
		{0x0, KeyframeVideoType},
		{0x1, KeyframeVideoType},
		{0x2, KeyframeVideoType},
		{0x3, KeyframeVideoType},
		{0x4, KeyframeVideoType},
	} {
		d := new(Video)
		d.data.data = []byte{c.Control}

		assert.Equal(t, c.VideoType, d.Type())
	}
}
