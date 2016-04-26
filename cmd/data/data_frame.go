package data

import (
	"bytes"

	"github.com/WatchBeam/amf0"
	"github.com/WatchBeam/amf0/encoding"
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
func (d *DataFrame) Read(b []byte) error {
	return encoding.Unmarshal(bytes.NewReader(b), d)
}
