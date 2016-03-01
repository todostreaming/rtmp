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
