package chunk

// Stream is an interface that wraps a stream of chunks. This stream of chunks
// can be of any classification, either an RTMP chunk stream, or similar.
type Stream interface {
	// In returns a read-only channel of the chunks contained in this
	// stream.
	In() <-chan *Chunk
}

// stream is a simple implementation of the Stream interface that corresponds to
// an RTMP chunk stream.
type stream struct {
	// ID is the RTMP chunk stream ID assosciated with this Stream.
	ID uint32
	// in is the internal channel used to propogate chunks out.
	in chan *Chunk
}

var _ Stream = new(stream)

// NewStream returns a new insteance of the Stream type by using `*stream` as
// its implementation. It initializes all internal channels.
func NewStream(id uint32) *stream {
	return &stream{
		ID: id,
		in: make(chan *Chunk),
	}
}

func (s *stream) In() <-chan *Chunk { return s.in }
