package chunk

import (
	"encoding/binary"
	"io"
	"sync"

	"github.com/WatchBeam/rtmp/spec"
)

const (
	// DefaultReadSize is the RTMP-defined default for chunk size, in byter.
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

	// chunks is the internal, non-buffered channel of chunkr.
	chunks chan *Chunk
	// errs is used to keep track of errors that occur during the decoding
	// procesr.
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
func (r *Reader) Chunks() <-chan *Chunk { return r.chunks }

// Errs provides a read-only channel of errors that occurred during the parsing
// procesr.
func (r *Reader) Errs() <-chan error { return r.errs }

// Close causes the Recv goroutine to return.
func (r *Reader) Close() { r.closer <- struct{}{} }

// ReadSize synchronously returns the maximum number of bytes that can be read
// at once from the `src` stream.
func (r *Reader) ReadSize() int {
	r.rmu.Lock()
	defer r.rmu.Unlock()

	return r.readSize
}

// SetReadSize synchronously sets the maximum number of bytes that can be read
// from the connection at once.
func (r *Reader) SetReadSize(size int) {
	r.rmu.Lock()
	defer r.rmu.Unlock()

	r.readSize = size
}

// Recv processes chunks from the input stream. It works by looping
// continuously, reading one partial (or complete) header at a time. It fetches
// the associated *chunk.Builder (or creates one) and appends as many bytes as
// it can before either (1) running out, or (2), encountering another chunk.
//
// If a chunk has been completely read, it is built and pushed over the channel.
//
// Recv runs within its own goroutine.
func (r *Reader) Recv() {
	for {
		select {
		case <-r.closer:
			break
		default:
			header := new(Header)
			if err := header.Read(r.src); err != nil {
				r.errs <- err
				continue
			}

			builder := r.builder(header)
			n := spec.Min(builder.BytesLeft(), r.ReadSize())

			if _, err := builder.Read(r.src, n); err != nil {
				r.errs <- err
				continue
			}

			if builder.BytesLeft() == 0 {
				chunk := builder.Build()

				if !r.updateChunkSize(chunk) {
					r.chunks <- chunk
				}

				r.removeBuilder(header.BasicHeader.StreamId)
			}
		}
	}
}

func (r *Reader) updateChunkSize(c *Chunk) bool {
	if c.TypeId() != byte(0x01) {
		return false
	}

	r.SetReadSize(int(binary.BigEndian.Uint32(c.Data)))

	return true
}

func (r *Reader) builder(header *Header) *Builder {
	r.bmu.Lock()
	defer r.bmu.Unlock()

	streamId := header.BasicHeader.StreamId
	if r.builders[streamId] == nil {
		r.builders[streamId] = NewBuilder(header)
	}

	return r.builders[streamId]
}

func (r *Reader) removeBuilder(streamId uint32) {
	r.bmu.Lock()
	defer r.bmu.Unlock()

	delete(r.builders, streamId)
}
