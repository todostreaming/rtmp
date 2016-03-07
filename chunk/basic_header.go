package chunk

import (
	"bytes"
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

func (h *BasicHeader) Write(w io.Writer) error {
	buf := make([]byte, 1)
	buf[0] = b.FormatId << 6

	var csId byte

	switch {
	case h.StreamId < 64:
		buf[0] |= (byte(h.StreamId) & 63)
	case h.StreamId < 320:
		csId = 1
		buf = append(buf, byte(h.StreamId-64))
	default:
		csId = 2

		tmp := new(bytes.Buffer)
		if _, err := spec.PutUint16(uint16(h.StreamId-64), tmp); err !=
			nil {

			return err
		}

		buf = append(buf, tmp.Bytes()...)
	}

	buf[0] |= csId

	if _, err := w.Write(buf); err != nil {
		return err
	}

	return nil
}
