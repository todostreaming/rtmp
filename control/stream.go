package control

import "github.com/todostreaming/rtmp/chunk"

// Stream represents an RTMP-compliant bi-directional transfer of RTMP control
// sequences. It parses control sequences out of a chunk.Stream, and writes them
// back when they are sent into the stream.
type Stream struct {
	chunks chunk.Stream
	writer chunk.Writer

	in     chan Control
	out    chan Control
	errs   chan error
	closer chan struct{}

	parser  Parser
	chunker Chunker
}

// NewStream returns a new instance of the Stream type initialized with the
// given chunk stream, parser, and chunker.
func NewStream(chunks chunk.Stream, writer chunk.Writer,
	parser Parser, chunker Chunker) *Stream {

	return &Stream{
		chunks: chunks,
		writer: writer,

		in:     make(chan Control),
		out:    make(chan Control),
		errs:   make(chan error),
		closer: make(chan struct{}),

		parser:  parser,
		chunker: chunker,
	}
}

// In is written to when incoming control sequences are read off of the stream
func (s *Stream) In() <-chan Control { return s.in }

// Out is written to by callers when they want to write a control sequence to
// the stream
func (s *Stream) Out() chan<- Control { return s.out }

// Errs is written to when an error is encountered from the chunk stream, or an
// error is encountered in chunking or parsing.
func (s *Stream) Errs() <-chan error { return s.errs }

// Close stops the Recv goroutine.
func (s *Stream) Close() { s.closer <- struct{}{} }

// Recv processes input from all channels, as well as the incoming and outgoing
// chunk streams.
//
// Recv runs within its own goroutine.
func (s *Stream) Recv() {
	defer func() {
		close(s.in)
		close(s.out)
		close(s.errs)
		close(s.closer)
	}()

	for {
		select {
		case <-s.closer:
			return
		case c := <-s.chunks.In():
			control, err := s.parser.Parse(c)
			if err != nil {
				s.errs <- err
				continue
			}

			s.in <- control
		case control := <-s.out:
			chunk, err := s.chunker.Chunk(control)
			if err != nil {
				s.errs <- err
				continue
			}

			if err = s.writer.Write(chunk); err != nil {
				s.errs <- err
				continue
			}
		}
	}
}
