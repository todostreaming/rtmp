package chunk

// A Normalizer is responsible for filling in "missing" pieces of information
// from RTMP chunks. This "filling in" process includes merging partial type 1
// and type 2 message headers, as well as replacing "empty" type 3 message
// headers.
//
// A cache of headers, as well as the last full chunk that was received is
// always stored, and those are used in tandem to complete the process described
// above.
type Normalizer interface {
	// Normalize can essentially be thought of as taking an "incomplete"
	// chunk, missing some header-data, and returning a complete chunk, with
	// the missing information filled in.
	//
	// For Type 1 and 2 headers, this means filling in the stream ID and
	// length from the last chunk that was received on any chunk stream. For
	// Type 3 headers, this means replacing the "missing" message header,
	// with the last full message header sent over the matching chunk stream
	// ID.
	//
	// Calling Normalize also updates the last received chunk to the one
	// that was just normalized, eliminating the need to call the
	// "Set<chunk|last>" methods.
	Normalize(chunk *Chunk)

	// Last	returns the last full chunk that was received over any chunk
	// stream, in a synchronous fashion.
	Last() *Chunk
	// SetLast sets the last received chunk received over any chunk stream,
	// in a synchronous fashion.
	SetLast(*Chunk)

	// Header returns the last "full" header received over the given chunk
	// stream, in a synchronous fashion.
	Header(streamId uint32) *Header

	// StoreHeader updates the last full header received over the given
	// header's chunk stream in a synchronous fashion.
	StoreHeader(*Header)
}
