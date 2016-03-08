package chunk

import (
	"bytes"
	"errors"
	"io"

	"github.com/WatchBeam/rtmp/spec"
)

var (
	ErrUnknownFormatId = errors.New("rtmp: unknown message header ID")
)

type MessageHeader struct {
	FormatId       byte
	Timestamp      uint32
	TimestampDelta bool
	Length         uint32
	TypeId         byte
	StreamId       uint32
}

func (m *MessageHeader) HasExtendedTimestamp() bool {
	return m.Timestamp == 0xffffff
}

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
