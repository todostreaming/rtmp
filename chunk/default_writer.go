package chunk

import (
	"bytes"
	"io"
	"sync"

	"github.com/WatchBeam/rtmp/spec"
)

type DefaultWriter struct {
	dest io.Writer

	wmu       sync.Mutex
	writeSize int
}

var _ Writer = new(DefaultWriter)

func (w *DefaultWriter) WriteSize() int {
	w.wmu.Lock()
	defer w.wmu.Unlock()

	return w.writeSize
}

func (w *DefaultWriter) SetWriteSize(writeSize int) {
	w.wmu.Lock()
	defer w.wmu.Unlock()

	w.writeSize = writeSize
}

func (w *DefaultWriter) Write(c *Chunk) error {
	if err := c.Header.Write(w.dest); err != nil {
		return err
	}

	payload := bytes.NewBuffer(c.Data)
	for payload.Len() > 0 {
		len := payload.Len()
		n := spec.Min(len, w.writeSize)
		if _, err := io.CopyN(w.dest, payload, int64(n)); err != nil {
			return err
		}

		// HACK(taylor): move this up to the chunk level
		if payload.Len() > 0 {
			partialHeader := []byte{byte(
				(3 << 6) | (c.Header.BasicHeader.StreamId & 63)),
			}

			if _, err := w.dest.Write(partialHeader); err != nil {
				return err
			}
		}
	}

	return nil
}
