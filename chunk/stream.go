package chunk

type Stream struct {
	reader     Reader
	normalizer Normalizer

	chunks chan *Chunk
	errs   chan error
	closer chan struct{}
}

func NewStream(reader Reader, normalizer Normalizer) *Stream {
	return &Stream{
		reader:     reader,
		normalizer: normalizer,
		chunks:     make(chan *Chunk),
		errs:       make(chan error),
		closer:     make(chan struct{}),
	}
}

func (s *Stream) Chunks() <-chan *Chunk { return s.chunks }
func (s *Stream) Errs() <-chan error    { return s.errs }
func (s *Stream) Close()                { s.closer <- struct{}{} }

func (s *Stream) Recv() {
	go s.reader.Recv()

	for {
		select {
		case chunk := <-s.reader.Chunks():
			s.normalizer.Normalize(chunk)
			s.chunks <- chunk
		case err := <-s.reader.Errs():
			s.errs <- err
		case <-s.closer:
			s.reader.Close()

			close(s.chunks)
			close(s.errs)
			close(s.closer)

			return
		}
	}

}
