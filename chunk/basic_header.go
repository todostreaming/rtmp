package chunk

import (
	"io"

	"github.com/WatchBeam/rtmp/spec"
)

type BasicHeader struct {
	FormatId byte
	StreamId uint32
}

func (h *BasicHeader) Read(r io.Reader) error {
	b, err := spec.ReadByte(r)
	if err != nil {
		return err
	}

	h.FormatId = b >> 6

	if b&63 == 63 {
		tail, err := spec.ReadBytes(r, int(2))
		if err != nil {
			return err
		}

		h.StreamId = spec.Uint32(tail) + 64
	} else if b&63 == 0 {
		tail, err := spec.ReadBytes(r, int(1))
		if err != nil {
			return err
		}

		h.StreamId = spec.Uint32(tail) + 64
	} else {
		h.StreamId = spec.Uint32([]byte{b & 63})
	}

	return nil
}
