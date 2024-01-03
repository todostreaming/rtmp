package control

import (
	"io"

	"github.com/todostreaming/rtmp/spec"
)

type AbortMessage struct {
	ChunkStreamId uint32
}

var _ Control = new(AbortMessage)

func (c *AbortMessage) TypeId() byte { return 0x2 }

func (c *AbortMessage) Read(r io.Reader) error {
	buf, err := spec.ReadBytes(r, 4)
	if err != nil {
		return err
	}

	c.ChunkStreamId = spec.Uint32(buf)

	return nil
}

func (c *AbortMessage) Write(w io.Writer) error {
	if _, err := spec.PutUint32(c.ChunkStreamId, w); err != nil {
		return err
	}

	return nil
}
