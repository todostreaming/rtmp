package data

import "errors"

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
	Read(b []byte) error
}

// data is a simple implementation of part of the Data interface.
type data struct {
	// Control represents the control sequence appended to the front of
	// each data frame.
	Control byte
	// Payload represents the actual data encoded in each Data frame.
	Payload []byte
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
func (d *data) Read(b []byte) error {
	if len(b) < 1 {
		return ErrControlMissing
	}

	d.Control = b[0]
	d.Payload = b[1:]

	return nil
}
