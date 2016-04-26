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
	// For Type 1 and 2 basic headers, this means filling in the stream ID
	// and length from the last chunk that was received on any chunk stream.
	// For Type 3 headers, this means replacing the "missing" message
	// header, with the last full message header sent over the matching
	// chunk stream ID.
	//
	// Calling Normalize also updates the last received chunk to the one
	// that was just normalized, eliminating the need to call the
	// "Set<chunk|last>" methods.
	Normalize(*Header) *Header
}

var (
	// NoopNormalizer is a singleton implementation of the Normalizer type
	// that simply serves as a pass-through for all given headers.
	NoopNormalizer Normalizer = new(noopNormalizer)
)

// noopNormalizer is the internal implementation of the NoopNormalizer (see
// above).
type noopNormalizer struct{}

// Normalize implements Normalizer.Normalize.
func (n *noopNormalizer) Normalize(h *Header) *Header { return h }
