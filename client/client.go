package client

import "io"

// Client represents a client connected to a RTMP server (see
// github.com/WatchBeam/rtmp/server for more). Clients are able to be written to
// and read from, and may have additional metadata attached to them in the
// future.
type Client struct {
	// Conn represents the readable and writeable connection that links to
	// the client. This may be a net.Conn, or even just a bytes.Buffer.
	Conn io.ReadWriter
}

// New instantiates and returns a pointer to a new instance of type Client. The
// client is initialized with the given connection.
func New(conn io.ReadWriter) *Client {
	return &Client{
		Conn: conn,
	}
}
