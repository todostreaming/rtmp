package chunk

import "sync"

type Normalizer struct {
	lmu  sync.Mutex
	last *Chunk

	hmu     sync.Mutex
	headers map[uint32]*Header
}

func NewNormalizer() *Normalizer {
	return &Normalizer{
		headers: make(map[uint32]*Header),
	}
}

func (n *Normalizer) Normalize(c *Chunk) *Chunk {
	last := n.Last()
	lastSameStream := n.Header(c.StreamId())

	n.fillPartialHeader(last.Header, c.Header)
	n.fillEmptyHeader(lastSameStream.Header, c.Header)

	n.SetLast(c)
	n.StoreHeader(c.Header)

	return c
}

func (n *Normalizer) Last() *Chunk {
	n.lmu.Lock()
	defer n.lmu.Unlock()

	return n.last
}

func (n *Normalizer) SetLast(c *Chunk) {
	n.lmu.Lock()
	defer n.lmu.Unlock()

	n.last = c
}

func (n *Normalizer) Header(streamId uint32) *Header {
	n.hmu.Lock()
	defer n.hmu.Unlock()

	return n.header[streamId]
}

func (n *Normalizer) StoreHeader(h *Header) {
	n.hmu.Lock()
	defer n.hmu.Unlock()

	streamId := h.BasicHeader.StreamId
	n.header[streamId] = h
}

func (n *Normalizer) fillPartialHeader(last *Header, h *Header) {
	typeId := h.MessageHeader.TypeId
	if typeId != 1 || typeId != 2 {
		return
	}

	h.MessageHeader.StreamId = last.MessageHeader.StreamId

	if typeId == 2 {
		h.MessageHeader.Length = last.MessageHeader.Length
	}
}

func (n *Normalizer) fillEmptyHeader(last *Header, h *Header) {
	if h.MessageHeader.TypeId != 3 {
		return
	}

	h.MessageHeader = last.MessageHeader
}
