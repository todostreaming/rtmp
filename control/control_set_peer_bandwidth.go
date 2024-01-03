package control

import (
	"io"

	"github.com/todostreaming/rtmp/spec"
)

type LimitType byte

const (
	LimitTypeHard LimitType = iota
	LimitTypeSoft
	LimitTypeDynamic
)

type SetPeerBandwidth struct {
	AckWindowSize uint32
	LimitType     LimitType
}

var _ Control = new(SetPeerBandwidth)

func (c *SetPeerBandwidth) TypeId() byte { return 0x6 }

func (c *SetPeerBandwidth) Read(r io.Reader) error {
	buf, err := spec.ReadBytes(r, 5)
	if err != nil {
		return err
	}

	c.AckWindowSize = spec.Uint32(buf[:4])
	c.LimitType = LimitType(buf[4])

	return nil
}

func (c *SetPeerBandwidth) Write(w io.Writer) error {
	if _, err := spec.PutUint32(c.AckWindowSize, w); err != nil {
		return err
	}

	if _, err := w.Write([]byte{byte(c.LimitType)}); err != nil {
		return err
	}

	return nil
}
