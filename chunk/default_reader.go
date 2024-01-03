package chunk

import (
	"encoding/binary"
	"io"
	"sync"

	"github.com/todostreaming/rtmp/spec"
)

const (
	// DefaultReadSize is the RTMP-defined default for chunk size, in byter.
	DefaultReadSize int = 128
)

// DefaultReader provides an RTMP-compliant implementation to the Reader
// interface.
type DefaultReader struct {
	// src is the io.Reader that the multiplexed chunks are read from.
	src io.Reader

	// bmu guards builders
	bmu sync.Mutex
	// builders maps the chunk stream ID to an associated builder. Once a
	// chunk has been fully read, this entry is removed.
	builders map[uint32]*Builder

	// normalizer is the Normalizer used to normalize incoming headers.
	normalizer Normalizer

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

var _ Reader = new(DefaultReader)

// Chunks implements the `Chunks` func in the Reader interface.
func (r *DefaultReader) Chunks() <-chan *Chunk { return r.chunks }

// Errs implements the `Errs` func in the Reader interface.
func (r *DefaultReader) Errs() <-chan error { return r.errs }

// Close implements the `Close` func in the Reader interface.
func (r *DefaultReader) Close() { r.closer <- struct{}{} }

// ReadSize implements the `ReadSize` func in the Reader interface.
func (r *DefaultReader) ReadSize() int {
	r.rmu.Lock()
	defer r.rmu.Unlock()

	return r.readSize
}

// SetReadSize implements the `SetReadSize` func in the Reader interface.
func (r *DefaultReader) SetReadSize(size int) {
	r.rmu.Lock()
	defer r.rmu.Unlock()

	r.readSize = size
}

// Recv implements the `Recv` func in the Reader interface.
func (r *DefaultReader) Recv() {
	for {
		select {
		case <-r.closer:
			return
		default:
			header := new(Header)
			if err := header.Read(r.src); err != nil {
				r.errs <- err
				continue
			}
			header = r.normalizer.Normalize(header)

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

func (r *DefaultReader) updateChunkSize(c *Chunk) bool {
	if c.TypeId() != byte(0x01) {
		return false
	}

	r.SetReadSize(int(binary.BigEndian.Uint32(c.Data)))

	return true
}

func (r *DefaultReader) builder(header *Header) *Builder {
	r.bmu.Lock()
	defer r.bmu.Unlock()

	streamId := header.BasicHeader.StreamId
	if r.builders[streamId] == nil {
		r.builders[streamId] = NewBuilder(header)
	}

	return r.builders[streamId]
}

func (r *DefaultReader) removeBuilder(streamId uint32) {
	r.bmu.Lock()
	defer r.bmu.Unlock()

	delete(r.builders, streamId)
}
