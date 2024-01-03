package chunk

import (
	"bytes"
	"io"
	"sync"

	"github.com/todostreaming/rtmp/spec"
)

// DefaultWriter provides a default implementation of chunk.Writer interface.
type DefaultWriter struct {
	// dest is the io.Writer where chunks are written to.
	dest io.Writer

	// wmu guards writeSize.
	wmu sync.Mutex
	// writeSize is the maximum payload length of a single chunk that can
	// be written without haveing to write multiple chunks.
	writeSize int
}

var _ Writer = new(DefaultWriter)

// WriteSize implements the WriteSize function defined in the Writer interface.
func (w *DefaultWriter) WriteSize() int {
	w.wmu.Lock()
	defer w.wmu.Unlock()

	return w.writeSize
}

// SetWriteSize implements the SetWriteSize function defined in the Writer interface.
func (w *DefaultWriter) SetWriteSize(writeSize int) {
	w.wmu.Lock()
	defer w.wmu.Unlock()

	w.writeSize = writeSize
}

// Write implements the Write function defined in the Writer interface.
func (w *DefaultWriter) Write(c *Chunk) error {
	payload := bytes.NewBuffer(c.Data)
	out := new(bytes.Buffer)

	c.Header.Write(out)
	for payload.Len() > 0 {
		io.CopyN(out, payload, int64(spec.Min(payload.Len(),
			w.WriteSize())))

		if payload.Len() > 0 {
			out.Write([]byte{byte(
				(3 << 6) | (c.Header.BasicHeader.StreamId & 63)),
			})
		}
	}

	if _, err := io.Copy(w.dest, out); err != nil {
		return err
	}

	return nil
}
