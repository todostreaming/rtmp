package chunk

import "io"

// Header represents an RTMP-compliant chunk header, consisting of a
// BasicHeader, a MessageHeader, and an ExtendedTimestamp.
type Header struct {
	// BasicHeader is the RTMP BasicHeader of the given Header.
	BasicHeader BasicHeader
	// MessageHeader is the RTMP MessageHeader of the given Header.
	MessageHeader MessageHeader
	// ExtendedTimestamp is the RTMP ExtendedTimestamp of the given Header.
	ExtendedTimestamp ExtendedTimestamp
}

// Read reads a Header (partial or complete) from the given io.Reader, by
// delegating into the Read methods of each component part of this header. If
// any error is encountered during the read, then execution will be halted, and
// the error will be returned immediately.
func (h *Header) Read(r io.Reader) error {
	if err := h.BasicHeader.Read(r); err != nil {
		return err
	}

	h.MessageHeader.FormatId = h.BasicHeader.FormatId
	if err := h.MessageHeader.Read(r); err != nil {
		return err
	}

	if h.MessageHeader.HasExtendedTimestamp() {
		if err := h.ExtendedTimestamp.Read(r); err != nil {
			return err
		}
	}

	return nil
}

// Write serializes and writes this Header to the given io.Writer by delegating
// into each of the component parts of the header. If any error is returned from
// the individual writes, then it will be returned immediately, and the Header
// CAN NOT be considered to be fully written.
func (h *Header) Write(w io.Writer) error {
	if err := h.BasicHeader.Write(w); err != nil {
		return err
	}

	if err := h.MessageHeader.Write(w); err != nil {
		return err
	}

	if h.MessageHeader.Timestamp == 0xffffff {
		if err := h.ExtendedTimestamp.Write(w); err != nil {
			return err
		}
	}

	return nil
}
