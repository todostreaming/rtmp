package chunk

import (
	"errors"
	"io"
	"sync"
)

var (
	ErrTooManyBytes = errors.New("rtmp: read too many bytes into chunk")
)

type Builder struct {
	Header   *Header
	Payloads [][]byte

	lmu  sync.Mutex
	left int
}

func (b *Builder) AddLeft(left int) {
	b.lmu.Lock()
	defer b.lmu.Unlock()

	b.left += left
}

func (b *Builder) BytesLeft() int {
	b.lmu.Lock()
	defer b.lmu.Unlock()

	return b.left
}

func NewBuilder(h *Header) *Builder {
	return &Builder{
		Header: h,
		left:   int(h.MessageHeader.Length),
	}
}

func (b *Builder) Read(r io.Reader, n int) (int, error) {
	buf := make([]byte, n)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return n, err
	}

	return b.Append(buf)
}

func (b *Builder) Append(slice []byte) (int, error) {
	b.Payloads = append(b.Payloads, slice)
	b.left -= len(slice)

	if b.left < 0 {
		return len(slice), ErrTooManyBytes
	}

	return len(slice), nil
}

func (b *Builder) Build() *Chunk {
	var payload []byte
	for _, partial := range b.Payloads {
		payload = append(payload, partial...)
	}

	return New(b.Header, payload)
}
