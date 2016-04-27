package control

import "io"

var (
	// Controls is a list of all Control types defined in the RTMP
	// specification.
	Controls []Control = []Control{
		&SetChunkSize{},
		&AbortMessage{},
		&Acknowledgement{},
		&Event{},
		&WindowAckSize{},
		&SetPeerBandwidth{},
	}
)

// Control represents an interface that encapsulates the various Control
// sequences as defined in the RTMP specification (see
// http://www.adobe.com/devnet/rtmp.html, section 5.4).
type Control interface {
	// Read reads the body of a given control sequence from the specified
	// io.Reader, returning any errors it encounters along the way.
	Read(io.Reader) error
	// Write marshals the body of the current control sequence into the
	// given io.Writer, returning any errors it encounters along the way.
	Write(io.Writer) error

	// TypeId returns the type ID for this given control sequence, as
	// defined by the RTMP specification.
	TypeId() byte
}
