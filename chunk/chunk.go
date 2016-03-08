package chunk

// Chunk represents an RTMP chunk as defined by the RTMP specification. It
// contains a Header, and some Data.
type Chunk struct {
	// Header is the chunk's Header (as defined by the RTMP specification).
	Header *Header
	// Data is the chunk's payload, and has a length equal to `Length` field
	// given in the Header.MessageHeader.
	Data []byte
}

// New returns a new Chunk initialized with the given Header and Data fields.
func New(header *Header, data []byte) *Chunk {
	return &Chunk{
		Header: header,
		Data:   data,
	}
}

// StreamId returns the ID of the RTMP chunk stream that this Chunk belongs to.
func (c *Chunk) StreamId() uint32 { return c.Header.BasicHeader.StreamId }

// TypeId returns the TypeId of this chunk's MessageHeader.
func (c *Chunk) TypeId() byte { return c.Header.MessageHeader.TypeId }
