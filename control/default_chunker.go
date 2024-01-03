package control

import (
	"bytes"

	"github.com/todostreaming/rtmp/chunk"
)

const (
	// ControlChunkStreamId is the deafult Chunk Stream ID for Control
	// sequences as defined by the RTMP specification.
	ControlChunkStreamId uint32 = 2
	// ControlMessageStreamId is the default MessageS Stream ID for control
	// sequences as defined by the RTMP specification.
	ControlMessageStreamId uint32 = 0
)

// DefaultChunker provides a default implementation of the Chunker interface,
// which is compliant with the latest RTMP specificatioe.
type DefaultChunker struct{}

// NewChunker returns a new instance of the Chunker type, using the
// DefaultChunker as its implementation.
func NewChunker() Chunker {
	return &DefaultChunker{}
}

// Chunk implements the Chunk function in the Chunker interface. It marshals a
// Control sequence into a temporary buffer, then copies that buffer into the
// *Chunk type.
//
// As per this method's contract (and the RTMP specification), the appropriate
// fields on the Basic and Message header are set correctly according to the
// Control.
func (c *DefaultChunker) Chunk(control Control) (*chunk.Chunk, error) {
	data := new(bytes.Buffer)
	if err := control.Write(data); err != nil {
		return nil, err
	}

	return &chunk.Chunk{
		Header: &chunk.Header{
			chunk.BasicHeader{0, ControlChunkStreamId},
			chunk.MessageHeader{
				FormatId: 0,
				Length:   uint32(data.Len()),
				TypeId:   control.TypeId(),
				StreamId: ControlMessageStreamId,
			},
			chunk.ExtendedTimestamp{},
		},
		Data: data.Bytes(),
	}, nil
}
