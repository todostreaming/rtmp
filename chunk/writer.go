package chunk

import "io"

type Writer interface {
	Write(*Chunk) error

	WriteSize() int
	SetWriteSize(writeSize int)
}

func NewWriter(dest io.Writer, writeSize int) Writer {
	return &DefaultWriter{
		dest:      dest,
		writeSize: writeSize,
	}
}
