package chunk

// Stream represents a stream of all chunks over a given Reader and Writer pair.
// It works by reading and writing chunks, normalizing them as they come in.
// Every chunk that is received from the Stream can be assumed to have complete
// and total data owned by its message header, as data from the previous
// received header (overall, and of the given chunk stream) will be used to
// normalize it (see type Normalizer).
type Stream struct {
	// reader is the Reader that chunks are read from.
	reader Reader
	// writer is the Writer that chunks are writte to.
	writer Writer
	// normalizer is the Normalizer type that produces complete chunks.
	normalizer Normalizer

	// in holds a channel of chunks to be read from.
	in chan *Chunk
	// out holds a channel of chunks to be written to.
	out chan *Chunk
	// errs holds a channel of all errors encountered during the read/write
	// process.
	errs chan error
	// closer holds a channel that closes the Stream when anything is
	// written to it.
	closer chan struct{}
}

// NewStream allocates and returns a pointer to a new instance of the Stream
// type initialized with the given Reader, Writer and Normalizer. All internal
// channels are initialized and opened during this time.
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

// In returns a channel of `*Chunk`s that is written to when a complete chunk
// has been recieved from the member reader. Any chunks that are read from this
// channel can be assumed to have a complete serialization of their
// MessageHeader, using the rules set forth in the RTMP specification.
func (s *Stream) In() <-chan *Chunk { return s.in }

// Out returns a write-only channel of `*Chunk`s that are written to the
// underlying Writer when written to.
func (s *Stream) Out() chan<- *Chunk { return s.out }

// Errs returns a channel of errors that is written to if an error is
// encountered in reading, writing, or normalizing.
func (s *Stream) Errs() <-chan error { return s.errs }

// Close stops the reading and writing processes immediately after the next item
// in either the In() or Out() channel has been processed. It also closes down
// any/all channels owned by this Stream.
func (s *Stream) Close() { s.closer <- struct{}{} }

// Recv processes the input, output, closer, channels, and the normalizing of
// the input.
//
// Recv runs within its own goroutine.
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
