package chunk

type Stream struct {
	ID uint32
	in chan *Chunk
}

func NewStream(id uint32) *Stream {
	return &Stream{
		ID: id,
		in: make(chan *Chunk),
	}
}

func (s *Stream) In() <-chan *Chunk { return s.in }
