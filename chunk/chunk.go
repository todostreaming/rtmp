package chunk

type Chunk struct {
	Header *Header
	Data   []byte
}

func New(header *Header, data []byte) *Chunk {
	return &Chunk{
		Header: header,
		Data:   data,
	}
}

func (c *Chunk) TypeId() uint32 {
	return c.Header.MessageHeader.TypeId
}
