package chunk

import "sync"

// A Normalizer is responsible for filling in "missing" pieces of information
// from RTMP chunks. This "filling in" process includes merging partial type 1
// and type 2 message headers, as well as replacing "empty" type 3 message
// headers.
//
// A cache of headers, as well as the last full chunk that was received is
// always stored, and those are used in tandem to complete the process described
// above.
type Normalizer struct {
	// lmu guards last
	lmu sync.Mutex
	// last holds a pointer to the last full chunk that was received
	last *Chunk

	// hmu guards headers
	hmu sync.Mutex
	// headers maps chunk stream IDs (found in the basic header) -> complete
	// headers
	headers map[uint32]*Header
}

// NewNormalizer allocates and returns a pointer to a new instance of the
// Normalizer type.
func NewNormalizer() *Normalizer {
	return &Normalizer{
		headers: make(map[uint32]*Header),
	}
}

// Normalize can essentially be thought of as taking an "incomplete" chunk,
// missing some header-data, and returning a complete chunk, with the missing
// information filled in.
//
// For Type 1 and 2 headers, this means filling in the stream ID and length from
// the last chunk that was received on any chunk stream. For Type 3 headers,
// this means replacing the "missing" message header, with the last full message
// header sent over the matching chunk stream ID.
//
// Calling Normalize also updates the last received chunk to the one that was
// just normalized, eliminating the need to call the "Set<chunk|last>" methods.
func (n *Normalizer) Normalize(c *Chunk) *Chunk {
	last := n.Last()
	lastSameStream := n.Header(c.StreamId())

	if last != nil {
		n.fillPartialHeader(last.Header, c.Header)
	}
	if lastSameStream != nil {
		n.fillEmptyHeader(lastSameStream, c.Header)
	}

	n.SetLast(c)
	n.StoreHeader(c.Header)

	return c
}

// Last	returns the last full chunk that was received over any chunk stream, in
// a synchronous fashion.
func (n *Normalizer) Last() *Chunk {
	n.lmu.Lock()
	defer n.lmu.Unlock()

	return n.last
}

// SetLast sets the last received chunk received over any chunk stream, in a
// synchronous fashion.
func (n *Normalizer) SetLast(c *Chunk) {
	n.lmu.Lock()
	defer n.lmu.Unlock()

	n.last = c
}

// Header returns the last "full" header received over the given chunk stream,
// in a synchronous fashion.
func (n *Normalizer) Header(streamId uint32) *Header {
	n.hmu.Lock()
	defer n.hmu.Unlock()

	return n.headers[streamId]
}

// StoreHeader updates the last full header received over the given header's
// chunk stream in a synchronous fashion.
func (n *Normalizer) StoreHeader(h *Header) {
	n.hmu.Lock()
	defer n.hmu.Unlock()

	streamId := h.BasicHeader.StreamId
	n.headers[streamId] = h
}

// fillPartialHeader fills in partially empty chunk message headers, according
// to the RTMP spec.
func (n *Normalizer) fillPartialHeader(last *Header, h *Header) {
	typeId := h.MessageHeader.TypeId
	if typeId != 1 && typeId != 2 {
		return
	}

	h.MessageHeader.StreamId = last.MessageHeader.StreamId

	if typeId == 2 {
		h.MessageHeader.Length = last.MessageHeader.Length
	}
}

// fillEmptyHeaderp fills in completely empty chunk message headers, according
// to the RTMP spec.
func (n *Normalizer) fillEmptyHeader(last *Header, h *Header) {
	if h.MessageHeader.TypeId != 3 {
		return
	}

	h.MessageHeader = last.MessageHeader
}
