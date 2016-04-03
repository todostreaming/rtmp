package data

const (
	AudioTypeId byte = 0x08
)

const (
	UncompressedAudioCodec AudioCodec = iota
	ADPCMAudioCodec
	MP3AudioCodec
	UncompressedLittleEndianAudioCodec
	Nellymoser16AudioCodec
	Nellymoser8AudioCodec
	NellymoserAudioCodec
	G711AAudioCodec
	G711UAudioCodec
	HE_AACAudioCodec
	SPEEXAudioCodec
)

const (
	MonoAudioType AudioType = iota
	StereoAudioType
)

type (
	// AudioCodec represents a singleton definition of the Codec assosciated
	// with a specific frame of audio.
	AudioCodec byte

	// AudioType represents a singleton definition of the audio Type
	// assosicated with a specific frame of Audio.
	AudioType byte
)

// Audio implements the Data interface for a frame of Audio.
type Audio struct{ data }

var _ Data = new(Audio)

// Id implements the Data.Id function.
func (a *Audio) Id() byte { return AudioTypeId }

// Codec retrns the AudioCodec assosciated with this frame of audio.
func (a *Audio) Codec() AudioCodec { return AudioCodec((a.Control & 0xf0) >> 4) }

// Rate returns the rate of audio contained in this frame in units of kHz.
func (a *Audio) Rate() float32 {
	return float32(5.5) * float32(2^((a.Control&0x0c)>>2))
}

// Size returns the audio sizes in bits.
func (a *Audio) Size() int {
	return 8 * int(2^((a.Control&0x02)>>1))
}

// Type returns the AudioType assosciated with this frame of Audio.
func (a *Audio) Type() AudioType { return AudioType(a.Control & 0x01) }
