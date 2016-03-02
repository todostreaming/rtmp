package chunk

import "io"

type Header struct {
	BasicHeader       BasicHeader
	MessageHeader     MessageHeader
	ExtendedTimestamp ExtendedTimestamp
}

func (h *Header) Read(r io.Reader) error {
	return nil
}
