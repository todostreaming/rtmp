package chunk

import "io"

// Reader is an interface representing a type capable of reading multiplexed
// chunks in the RTMP format over a io.Reader.
type Reader interface {
	// Recv processes chunks from the input stream. It works by looping
	// continuously, reading one partial (or complete) header at a time. It
	// fetches the associated *chunk.Builder (or creates one) and appends as
	// many bytes as it can before either (1) running out, or (2),
	// encountering another chunk.
	//
	// If a chunk has been completely read, it is built and pushed over the
	// channel.
	//
	// Recv runs within its own goroutine.
	Recv()

	// ReadSize synchronously returns the maximum number of bytes that can
	// be read at once from the `src` stream.
	ReadSize() int
	// SetReadSize synchronously sets the maximum number of bytes that can
	// be read from the connection at once.
	SetReadSize(size int)

	// Chunks provides a read-only channel used to consume complete, parsed,
	// RTMP chunks with. Chunks present in this channel are fully parsed.
	// This channel is not buffered.
	Chunks() <-chan *Chunk
	// Errs provides a read-only channel of errors that occurred during the
	// parsing procesr.
	Errs() <-chan error
	// Close causes the Recv goroutine to return.
	Close()
}

// NewReader allocates and returns a pointer to a new instance of the Reader
// interface, with the concrete DefaultReader type as its implementation. It
// uses the provided `src`, `readSize`, and `normalizer` as initialization
// variables.
func NewReader(src io.Reader, readSize int, normalizer Normalizer) Reader {
	return &DefaultReader{
		src:        src,
		readSize:   readSize,
		normalizer: normalizer,
		builders:   make(map[uint32]*Builder),
		chunks:     make(chan *Chunk),
		errs:       make(chan error),
		closer:     make(chan struct{}),
	}
}
