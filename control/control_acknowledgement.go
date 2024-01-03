package control

import (
	"io"

	"github.com/todostreaming/rtmp/spec"
)

type Acknowledgement struct {
	SequenceNumber uint32
}

var _ Control = new(Acknowledgement)

func (c *Acknowledgement) TypeId() byte { return 0x3 }

func (c *Acknowledgement) Read(r io.Reader) error {
	buf, err := spec.ReadBytes(r, 4)
	if err != nil {
		return err
	}

	c.SequenceNumber = spec.Uint32(buf)

	return nil
}

func (c *Acknowledgement) Write(w io.Writer) error {
	if _, err := spec.PutUint32(c.SequenceNumber, w); err != nil {
		return err
	}

	return nil
}
