package chunk

import (
	"bytes"
	"errors"
	"io"

	"github.com/todostreaming/rtmp/spec"
)

var (
	// ErrUnknownFormatId is returned in a situation where the Format ID
	// belonging to an attempted read of a MessageHeader is not included
	// within the RTMP specification.
	ErrUnknownFormatId = errors.New("rtmp: unknown message header ID")
)

// MessageHeader represents the MessageHeader component of a chunk Header, as
// defined in the RTMP specification.
type MessageHeader struct {
	// FormatId is the FormatId of the Message header (should it be
	// read from or written to a chunk stream). It must be set from the
	// BasicHeader's FormatId elsewhere.
	FormatId byte
	// Timestamp is the 4-byte (perhaps partial) encoding of the timestamp
	// (or delta, if TimestampDelta field is true).
	Timestamp uint32
	// TimestampDelta is a bool that tells callers whether or not the
	// Timestamp field is a absolute, or relative timestamp.
	TimestampDelta bool
	// Length is the length in bytes of the associated chunk's payload.
	Length uint32
	// TypeId is the RTMP typeId as defined by the RTMP specification.
	TypeId byte
	// StreamId is the ID of the chunk stream that this header's chunk
	// belongs to.
	StreamId uint32
}

// HasExtendedTimestamp determines whether or not an ExtendedTimestamp header is
// necessary to encode the full timestamp.
func (m *MessageHeader) HasExtendedTimestamp() bool {
	return m.Timestamp == 0xffffff
}

// Read reads a type 0, 1, 2, or 3-format MessageHeader from the given
// io.Reader, using a process as defined in the RTMP specification.
func (m *MessageHeader) Read(r io.Reader) error {
	switch m.FormatId {
	case 0:
		buf, err := spec.ReadBytes(r, 11)
		if err != nil {
			return err
		}

		m.Timestamp = spec.Uint32(buf[:3])
		m.TimestampDelta = false
		m.Length = spec.Uint32(buf[3:6])
		m.TypeId = buf[6]
		m.StreamId = spec.LittleEndianUint32(buf[7:])
	case 1:
		buf, err := spec.ReadBytes(r, 7)
		if err != nil {
			return err
		}

		m.Timestamp = spec.Uint32(buf[:3])
		m.TimestampDelta = true
		m.Length = spec.Uint32(buf[3:6])
		m.TypeId = buf[6]
	case 2:
		buf, err := spec.ReadBytes(r, 3)
		if err != nil {
			return err
		}

		m.Timestamp = spec.Uint32(buf)
		m.TimestampDelta = true
	case 3:
		return nil
	default:
		return ErrUnknownFormatId
	}

	return nil
}

// Write encodes and writes the data held by this MessageHeader to the given
// io.Writer, returning any error encountered during the write as it occurs.
func (m *MessageHeader) Write(w io.Writer) error {
	buf := new(bytes.Buffer)

	switch m.FormatId {
	case 0:
		spec.PutUint24(m.Timestamp, buf)
		spec.PutUint24(m.Length, buf)
		spec.PutUint8(m.TypeId, buf)
		spec.LittleEndianPutUint32(m.StreamId, buf)
	case 1:
		spec.PutUint24(m.Timestamp, buf)
		spec.PutUint24(m.Length, buf)
		spec.PutUint8(m.TypeId, buf)
	case 2:
		spec.PutUint24(m.Timestamp, buf)
	}

	if _, err := w.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}
