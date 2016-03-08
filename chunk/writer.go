package chunk

import "io"

// Writer represents an interface capable of chunking and writing RTMP Chunks
// according to a given (mutable) WriteSize.
//
// TODO(taylor): this interface is a little weird in that it doesn't own a
// io.Writer. Let's see what people think of this and then go from there.
type Writer interface {
	// Write writes a chunk to the owned io.Writer, returning any error that
	// it encountered along the way. Should the payload of the chunk exceed
	// the maximum write size (the value held in WriteSize(), then the chunk
	// should be split into multiple parts, seperated by a partial message
	// header.
	Write(*Chunk) error

	// WriteSize returns the maximum length of a chunk's payload before that
	// chunk must be split into multiple chunks.
	WriteSize() int
	// SetWriteSize changes the write size of this particular Writer.
	SetWriteSize(writeSize int)
}

// NewWriter returns a default implementation of the Writer interface.
func NewWriter(dest io.Writer, writeSize int) Writer {
	return &DefaultWriter{
		dest:      dest,
		writeSize: writeSize,
	}
}
