package chunk

type Stream struct {
	reader     Reader
	writer     Writer
	normalizer Normalizer

	in     chan *Chunk
	out    chan *Chunk
	errs   chan error
	closer chan struct{}
}

func NewStream(reader Reader, writer Writer, normalizer Normalizer) *Stream {
	return &Stream{
		reader:     reader,
		writer:     writer,
		normalizer: normalizer,
		in:         make(chan *Chunk),
		out:        make(chan *Chunk),
		errs:       make(chan error),
		closer:     make(chan struct{}),
	}
}

func (s *Stream) In() <-chan *Chunk  { return s.in }
func (s *Stream) Out() chan<- *Chunk { return s.out }
func (s *Stream) Errs() <-chan error { return s.errs }
func (s *Stream) Close()             { s.closer <- struct{}{} }

func (s *Stream) Recv() {
	go s.reader.Recv()

	for {
		select {
		case in := <-s.reader.Chunks():
			s.normalizer.Normalize(in)
			s.in <- in
		case out := <-s.out:
			if err := s.writer.Write(out); err != nil {
				s.errs <- err
			}
		case err := <-s.reader.Errs():
			s.errs <- err
		case <-s.closer:
			s.reader.Close()

			close(s.in)
			close(s.out)
			close(s.errs)
			close(s.closer)

			return
		}
	}

}
