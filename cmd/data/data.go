package data

import (
	"errors"

	"github.com/todostreaming/rtmp/chunk"
)

var (
	ErrControlMissing = errors.New("rtmp/data: missing control byte")
)

// Data represents a single frame of data coming over the RTMP chunk stream.
// A Data knows about its RTMP chunk's Type ID, as well as how to read itself
// from a slice of bytes.
type Data interface {
	// Id returns the type ID as a byte that is associated with this frame
	// of Data. This should be equivalent to the ID found in
	// Chunk.Header.MessageHeader.TypeId.
	Id() byte

	// Read is a destructive operation which reads data into this chunk
	// stream by feeding off of the given byte slice. In normal operation,
	// the payload of an RTMP chunk is used here.
	Read(c *chunk.Chunk) error

	// Marshal marshals the data that is encoded in the implementation of
	// Data that this method belongs to in a writeable *chunk.Chunk, or an
	// error if the data was unable to be marshaled.
	Marshal() (*chunk.Chunk, error)
}

// data is a simple implementation of part of the Data interface.
type data struct {
	// header is the *chunk.Header that was read during Read().
	header *chunk.Header
	// data is the data read from the *chunk.Chunk verbatim
	data []byte
}

// Read implements the Data.Read function. In the best case, it assigns the
// first byte in `b` to the Control sequence, and the last `n-1` bytes from `b`
// to be the encoded data.
//
// If less than one byte exists on the byte slice `b`, then an error is thrown
// since the control is missing.
//
// Otherwise a value of nil is returned, signifying that the Read operation
// succeeded without error.
func (d *data) Read(c *chunk.Chunk) error {
	if len(c.Data) < 1 {
		return ErrControlMissing
	}

	d.header = c.Header
	d.data = c.Data

	return nil
}

// Control represents the control sequence appended to the front of
// each data frame.
func (d *data) Control() byte { return d.data[0] }

// Payload represents the actual data encoded in each Data frame.
func (d *data) Payload() []byte { return d.data[1:] }

// Marshal implements the Data.Marshal, using the same header that was sent
// during the original read.
func (d *data) Marshal() (*chunk.Chunk, error) {
	return &chunk.Chunk{
		Header: d.header,
		Data:   d.data,
	}, nil
}
