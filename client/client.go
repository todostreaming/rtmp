package client

import (
	"io"

	"github.com/WatchBeam/rtmp/chunk"
	"github.com/WatchBeam/rtmp/control"
)

// Client represents a client connected to a RTMP server (see
// github.com/WatchBeam/rtmp/server for more). Clients are able to be written to
// and read from, and may have additional metadata attached to them in the
// future.
type Client struct {
	chunks      *chunk.Parser
	chunkWriter chunk.Writer

	controlStream *control.Stream

	// Conn represents the readable and writeable connection that links to
	// the client. This may be a net.Conn, or even just a bytes.Buffer.
	Conn io.ReadWriter
}

// New instantiates and returns a pointer to a new instance of type Client. The
// client is initialized with the given connection.
func New(conn io.ReadWriter) *Client {
	chunkWriter := chunk.NewWriter(conn, chunk.DefaultReadSize)
	chunks := chunk.NewParser(
		chunk.NewReader(conn, chunk.DefaultReadSize),
		chunk.NewNormalizer(),
	)

	return &Client{
		chunks:      chunks,
		chunkWriter: chunkWriter,

		controlStream: control.NewStream(
			chunks.Stream(2),
			chunkWriter,
			control.NewParser(),
			control.NewChunker(),
		),

		Conn: conn,
	}
}

// Controls returns the stream of control sequences that are being received
// from the connected client.
func (c *Client) Controls() *control.Stream { return c.controlStream }
