package chunk

import (
	"io"

	"github.com/WatchBeam/rtmp/spec"
)

type ExtendedTimestamp struct {
	Delta uint32
}

func (t *ExtendedTimestamp) Read(r io.Reader) error {
	buf, err := spec.ReadBytes(r, 4)
	if err != nil {
		return err
	}

	t.Delta = spec.Uint32(buf)
	return nil
}
