package chunk

import (
	"errors"
	"io"
	"sync"
)

var (
	ErrTooManyBytes = errors.New("rtmp: read too many bytes into chunk")
)

// Builder represents a type capable of "stiching" together parts of a chunk
// received in a multiplexed chunk stream, to make a complete RTMP Chunk.
type Builder struct {
	// Header holds a pointer to the FIRST header that was received along
	// with the chunk that we are trying to build.
	//
	// As a note, if you want to ship off COMPLETE chunk headers (with the
	// missing data filled in), the chunk.Normalizer type may be of use to
	// you.
	Header *Header

	// pmu guards Builder.Payloads.
	pmu sync.Mutex
	// Payloads holds a slice of slices, each sub-slice being one part of
	// the payload received at a time. These slices are appended together to
	// form a complete chunk.
	Payloads [][]byte

	// lmu guards left
	lmu sync.Mutex
	// left returns the number of bytes left in the chunk.
	left int
}

// NewBuilder allocates and returns a pointer to a new instance of the Builder
// type. It is initialized with the given header as its chunk header, as well as
// the number of bytes left to be the `Length` field of the chunk's
// MessageHeader.
func NewBuilder(h *Header) *Builder {
	return &Builder{
		Header: h,
		left:   int(h.MessageHeader.Length),
	}
}

// Build builds a returns a Chunk formed by using the given header for the
// chunk's header, and the concatenation of all of the received payloads as the
// chunk body.
func (b *Builder) Build() *Chunk {
	b.pmu.Lock()
	defer b.pmu.Unlock()

	var payload []byte
	for _, partial := range b.Payloads {
		payload = append(payload, partial...)
	}

	return New(b.Header, payload)
}

// Read reads the given number of bytes from the specified io.Reader, and
// appends them as a slice. It returns the number of bytes read, and any error
// encountered (if applicable).
func (b *Builder) Read(r io.Reader, n int) (int, error) {
	buf := make([]byte, n)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return n, err
	}

	return b.Append(buf)
}

// Append appends the given slice to the payloads received thus far, subtracting
// the number of bytes read from the number of bytes left to be read. If too may
// bytes were appended (i.e., left < 0), then ErrTooManyBytes will be returned
// instead.
func (b *Builder) Append(slice []byte) (int, error) {
	b.pmu.lock()
	defer b.pmu.Unlock()

	b.Payloads = append(b.Payloads, slice)
	b.left -= len(slice)

	if b.left < 0 {
		return len(slice), ErrTooManyBytes
	}

	return len(slice), nil
}

// AddLeft adds the delta parameter to the amount of bytes left by using the
// guarding mutex.
func (b *Builder) AddLeft(delta int) {
	b.lmu.Lock()
	defer b.lmu.Unlock()

	b.left += delta
}

// BytesLeft returns the number of bytes stil to be read.
func (b *Builder) BytesLeft() int {
	b.lmu.Lock()
	defer b.lmu.Unlock()

	return b.left
}
