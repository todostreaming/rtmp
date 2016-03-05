package chunk

import (
	"io"
	"sync"

	"github.com/WatchBeam/rtmp/spec"
)

const (
	// DefaultReadSize is the RTMP-defined default for chunk size, in bytes.
	DefaultReadSize int = 128
)

type Reader struct {
	// src is the io.Reader that the multiplexed chunks are read from.
	src io.Reader

	// bmu guards builders
	bmu sync.Mutex
	// builders maps the chunk stream ID to an associated builder. Once a
	// chunk has been fully read, this entry is removed.
	builders map[uint32]*Builder

	// rmu guards readSize
	rmu sync.Mutex
	// readSize refers to the maximum amount of bytes that can be read at
	// once for one chunk.
	readSize int

	// chunks is the internal, non-buffered channel of chunks.
	chunks chan *Chunk
	// errs is used to keep track of errors that occur during the decoding
	// process.
	errs chan error
	// closer is a non-buffered channel used to pass closing signal around.
	closer chan struct{}
}

func NewReader(src io.Reader, readSize int) *Reader {
	return &Reader{
		src:      src,
		readSize: readSize,
		builders: make(map[uint32]*Builder),
		chunks:   make(chan *Chunk),
		errs:     make(chan error),
		closer:   make(chan struct{}),
	}
}

// Chunks provides a read-only channel used to consume complete, parsed, RTMP
// chunks with. Chunks present in this channel are fully parsed. This channel is
// not buffered.
func (s *Stream) Chunks() <-chan *Chunk { return s.chunks }

// Errs provides a read-only channel of errors that occurred during the parsing
// process.
func (s *Stream) Errs() <-chan error { return s.errs }

// Close causes the Recv goroutine to return.
func (s *Stream) Close() { s.closer <- struct{}{} }

// ReadSize synchronously returns the maximum number of bytes that can be read
// at once from the `src` stream.
func (s *Stream) ReadSize() int {
	s.rmu.Lock()
	defer s.rmu.Unlock()

	return s.readSize
}

// SetReadSize synchronously sets the maximum number of bytes that can be read
// from the connection at once.
func (s *Stream) SetReadSize(size int) {
	s.rmu.Lock()
	defer s.rmu.Unlock()

	s.readSize = size
}

// Recv processes chunks from the input stream. It works by looping
// continuously, reading one partial (or complete) header at a time. It fetches
// the associated *chunk.Builder (or creates one) and appends as many bytes as
// it can before either (1) running out, or (2), encountering another chunk.
//
// If a chunk has been completely read, it is built and pushed over the channel.
//
// Recv runs within its own goroutine.
func (s *Stream) Recv() {
	for {
		select {
		case <-s.closer:
			break
		default:
			header := new(Header)
			if err := header.Read(s.src); err != nil {
				s.errs <- err
				continue
			}

			builder := s.builder(header)
			n := spec.Min(builder.BytesLeft(), s.ReadSize())

			if err := builder.Read(s.src, n); err != nil {
				s.errs <- err
				continue
			}

			if builder.BytesLeft() == 0 {
				chunk := builder.Build()

				if !s.updateChunkSize(chunk) {
					s.chunks <- chunk
				}

				s.removeBuilder(Header.BasicHeader.StreamId)
			}
		}
	}
}
