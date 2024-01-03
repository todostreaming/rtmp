package control

import (
	"io"

	"github.com/todostreaming/rtmp/spec"
)

type WindowAckSize struct {
	WindowAckSize uint32
}

var _ Control = new(WindowAckSize)

func (c *WindowAckSize) TypeId() byte { return 0x5 }

func (c *WindowAckSize) Read(r io.Reader) error {
	buf, err := spec.ReadBytes(r, 4)
	if err != nil {
		return err
	}

	c.WindowAckSize = spec.Uint32(buf)

	return nil
}

func (c *WindowAckSize) Write(w io.Writer) error {
	if _, err := spec.PutUint32(c.WindowAckSize, w); err != nil {
		return err
	}

	return nil
}
