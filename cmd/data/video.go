package data

const (
	VideoTypeId byte = 0x09
)

const (
	SorensenH263VideoCodec VideoCodec = iota
	ScreenVideoVideoCodec
	On2VP6VideoCodec
	On2VP6AlphaVideoCodec
	ScreenVideo2VideoCodec
	H264VideoCodec
)

const (
	KeyframeVideoType VideoType = iota
	InterframeVideoType
	DisposableInterframeVideoType
	GeneratedKeyFrameVideoType
	CommandFrameVideoType
)

type (
	// VideoType is a singleton representation of what type of video data is
	// encoded in this frame of Video.
	VideoType byte

	// VideoCodec is a singleton representation of which codec was used to
	// encoded the data contained in this frame of Video.
	VideoCodec byte
)

// Video is an implementation of the Data interface for video frames.
type Video struct{ data }

var _ Data = new(Video)

// Id implements Data.Id.
func (v *Video) Id() byte { return VideoTypeId }

// Codec returns the VideoCodec assosciated with this frame of Video.
func (v *Video) Codec() VideoCodec { return VideoCodec((v.Control() & 0x0f) >> 0) }

// Type returns the VideoType assosciated with this frame of Video.
func (v *Video) Type() VideoType { return VideoType((v.Control() & 0xf0) >> 4) }
