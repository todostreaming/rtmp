package chunk

import "sync"

// DefaultNormalizer provides a default implementation of the Normalizer
// interface.
type DefaultNormalizer struct {
	// lmu guards last
	lmu sync.Mutex
	// last holds a pointer to the last full chunk that was received
	last *Header

	// hmu guards headers
	hmu sync.Mutex
	// headers maps chunk stream IDs (found in the basic header) -> complete
	// headers
	headers map[uint32]*Header
}

// NewNormalizer allocates and returns a pointer to a new instance of the
// Normalizer type.
func NewNormalizer() Normalizer {
	return &DefaultNormalizer{
		headers: make(map[uint32]*Header),
	}
}

// Normalize implements the `Normalize` func from the Normalizer interface.
func (n *DefaultNormalizer) Normalize(h *Header) *Header {
	last := n.Last()
	lastSameStream := n.Header(h.BasicHeader.StreamId)

	if last != nil {
		n.fillPartialHeader(last, h)
	}
	if lastSameStream != nil {
		n.fillEmptyHeader(lastSameStream, h)
	}

	n.SetLast(h)
	n.StoreHeader(h)

	return h
}

// Last implements the `Last` func from the Normalizer interface.
func (n *DefaultNormalizer) Last() *Header {
	n.lmu.Lock()
	defer n.lmu.Unlock()

	return n.last
}

// SetLast implements the `SetLast` func from the Normalizer interface.
func (n *DefaultNormalizer) SetLast(h *Header) {
	n.lmu.Lock()
	defer n.lmu.Unlock()

	n.last = h
}

// Header implements the `Header` func from the Normalizer interface.
func (n *DefaultNormalizer) Header(streamId uint32) *Header {
	n.hmu.Lock()
	defer n.hmu.Unlock()

	return n.headers[streamId]
}

// StoreHeader implements the `StoreHeader` func from the Normalizer interface.
func (n *DefaultNormalizer) StoreHeader(h *Header) {
	n.hmu.Lock()
	defer n.hmu.Unlock()

	streamId := h.BasicHeader.StreamId
	n.headers[streamId] = h
}

// fillPartialHeader fills in partially empty chunk message headers, according
// to the RTMP spec.
func (n *DefaultNormalizer) fillPartialHeader(last *Header, h *Header) {
	fmtId := h.BasicHeader.FormatId
	if fmtId != 1 && fmtId != 2 {
		return
	}

	h.MessageHeader.StreamId = last.MessageHeader.StreamId

	if fmtId == 2 {
		h.MessageHeader.Length = last.MessageHeader.Length
	}
}

// fillEmptyHeaderp fills in completely empty chunk message headers, according
// to the RTMP spec.
func (n *DefaultNormalizer) fillEmptyHeader(last *Header, h *Header) {
	if h.BasicHeader.FormatId != 3 {
		return
	}

	h.MessageHeader = last.MessageHeader
}
