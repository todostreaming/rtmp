package client

import (
	"io"

	"github.com/todostreaming/rtmp/chunk"
	"github.com/todostreaming/rtmp/cmd"
	"github.com/todostreaming/rtmp/control"
	"github.com/todostreaming/rtmp/handshake"
)

// Client represents a client connected to a RTMP server (see
// github.com/todostreaming/rtmp/server for more). Clients are able to be written to
// and read from, and may have additional metadata attached to them in the
// future.
type Client struct {
	chunks *chunk.Parser

	controlStream *control.Stream
	cmdManager    *cmd.Manager

	// Conn represents the readable and writeable connection that links to
	// the client. This may be a net.Conn, or even just a bytes.Buffer.
	Conn io.ReadWriter
}

// New instantiates and returns a pointer to a new instance of type Client. The
// client is initialized with the given connection.
func New(conn io.ReadWriter) *Client {
	chunkWriter := chunk.NewWriter(conn, 4096)
	chunks := chunk.NewParser(chunk.NewReader(
		conn, chunk.DefaultReadSize, chunk.NewNormalizer(),
	))

	controlChunks, _ := chunks.Stream(2)
	netChunks, _ := chunks.Stream(3, 4, 5, 8)

	return &Client{
		chunks: chunks,

		controlStream: control.NewStream(
			controlChunks,
			chunkWriter,
			control.NewParser(),
			control.NewChunker(),
		),

		cmdManager: cmd.New(netChunks, chunkWriter),

		Conn: conn,
	}
}

// Handshake preforms the handshake operation against the connecting client. If
// an error is encountered during any point of the handshake process, it will be
// returned immediately.
//
// If no error is encounterd while handshaking, the chunk reading process will
// begin.
//
// See github.com/todostreaming/RTMP/handshake for details.
func (c *Client) Handshake() error {
	if err := handshake.With(&handshake.Param{
		Conn: c.Conn,
	}).Handshake(); err != nil {
		return err
	}

	go c.chunks.Recv()

	return nil
}

// Controls returns the stream of control sequences that are being received
// from the connected client.
func (c *Client) Controls() *control.Stream { return c.controlStream }

// Net returns the *cmd.Manager responsible for handling the NetConnection,
// NetStrema, and DataStream exchanged with this client.
func (c *Client) Net() *cmd.Manager { return c.cmdManager }
