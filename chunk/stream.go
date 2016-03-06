package chunk

import "io"

type Stream struct {
	reader     *Reader
	normalizer Normalizer

	chunks chan *Chunk
	errs   chan error
	closer chan struct{}
}

func NewStream(src io.Reader) *Stream {
	return &Stream{
		reader:     NewReader(src, DefaultReadSize),
		normalizer: NewNormalizer(),
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
	defer s.reader.Close()

	for {
		select {
		case chunk := <-s.reader.Chunks():
			s.normalizer.Normalize(chunk)
			s.chunks <- chunk
		case err := <-s.reader.Errs():
			s.errs <- err
		case <-s.closer:
			break
		}
	}
}
