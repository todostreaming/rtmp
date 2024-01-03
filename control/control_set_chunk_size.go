package control

import (
	"io"

	"github.com/todostreaming/rtmp/spec"
)

type SetChunkSize struct {
	chunkSize uint32
}

var _ Control = new(SetChunkSize)

func NewSetChunkSize(size uint32) *SetChunkSize {
	return &SetChunkSize{
		chunkSize: size,
	}
}

func (c *SetChunkSize) TypeId() byte { return 0x1 }

func (c *SetChunkSize) Read(r io.Reader) error {
	buf, err := spec.ReadBytes(r, 4)
	if err != nil {
		return err
	}

	c.chunkSize = spec.Uint32(buf)

	return nil
}

func (c *SetChunkSize) Write(w io.Writer) error {
	if _, err := spec.PutUint32(c.ChunkSize(), w); err != nil {
		return err
	}

	return nil
}

func (c *SetChunkSize) ChunkSize() uint32 {
	if c.chunkSize > 0xffffff {
		return 0xffffff
	}

	return c.chunkSize
}
