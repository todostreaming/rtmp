package data

import (
	"bytes"

	"github.com/todostreaming/amf0"
	"github.com/todostreaming/amf0/encoding"
	"github.com/todostreaming/rtmp/chunk"
)

// DataFrame encapsulates the "@setDataFrame" type sent over the Data stream.
type DataFrame struct {
	// Header contains the "@setDataFrame" keyword.
	Header string
	// Type is the sub-type of the data frame packet.
	Type string
	// Arguments are the arguments that were sent in the packet.
	Arguments *amf0.Array
}

var _ Data = new(DataFrame)

// Id implements Data.Id.
func (d *DataFrame) Id() byte { return 0x12 }

// Read implements Data.Read. It uses the standard amf0-style procedure to
// unmarshal the amf0 encoded data.
func (d *DataFrame) Read(c *chunk.Chunk) error {
	return encoding.Unmarshal(bytes.NewReader(c.Data), d)
}

// Marshal implements the Data.Marshal function.
func (d *DataFrame) Marshal() (*chunk.Chunk, error) {
	m, err := encoding.Marshal(d)
	if err != nil {
		return nil, err
	}

	return &chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader: chunk.BasicHeader{0, 4},
			MessageHeader: chunk.MessageHeader{
				Length:   uint32(len(m)),
				TypeId:   0x12,
				StreamId: 1,
			},
		},
		Data: m,
	}, nil
}
