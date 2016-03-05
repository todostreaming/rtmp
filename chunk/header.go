package chunk

import "io"

type Header struct {
	BasicHeader       BasicHeader
	MessageHeader     MessageHeader
	ExtendedTimestamp ExtendedTimestamp
}

func (h *Header) Read(r io.Reader) error {
	if err := h.BasicHeader.Read(r); err != nil {
		return err
	}

	h.MessageHeader.FormatId = h.BasicHeader.FormatId
	if err := h.MessageHeader.Read(r); err != nil {
		return err
	}

	if err := h.ExtendedTimestamp.Read(r); err != nil {
		return err
	}
}
