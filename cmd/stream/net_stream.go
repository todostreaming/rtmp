package stream

import (
	"bytes"

	"github.com/todostreaming/rtmp/chunk"
)

// Type NetStream is an implementation of the NetStream type as described in the
// RTMP specification as published by Macromedia/Adobe.
//
// The NetStream provides a mechanism for the server to receive commands sent
// over by the client, with the optional ability to occasionally send the
// `onStatus` packet back.
//
// NetStream works in the typical Golang way, and it provides a <-chan of
// commands sent by the client, as well as a channel that is writeable to when
// the server wants to send an "onStatus" command back to the client command
// back to the client
type NetStream struct {
	// chunks is the incoming channel of chunks queued for processing by
	// this NetStream.
	chunks <-chan *chunk.Chunk
	// parser is the parser that is used to parse incoming commands.
	parser Parser
	// in is the outgoing channel written to when an incoming command has
	// been completely read and parsed, and is available to callers.
	in chan Command

	// statuses is a channel which is written to when a caller wants to
	// write an `onStatus` packet out to the connected client.
	statuses chan *Status
	// writer is the chunk.Writer where `onStatus` commands are written to.
	writer chunk.Writer

	// closer is a channel written to when the Listen operation should be
	// closed.
	closer chan struct{}
	// errs is a chnanel written to whenever an error is encountered during
	// the Listen goroutine.
	errs chan error
}

// New returns a new instance of the NetStream type, initialized with the given
// channel of chunks, and the specified chunk.Writer.
//
// Calling `New()` also instantiates the internal channels, but does not spawn
// the Listen operation.
func New(chunks <-chan *chunk.Chunk, writer chunk.Writer) *NetStream {
	return &NetStream{
		chunks: chunks,
		writer: writer,

		parser: DefaultParser,

		in:       make(chan Command),
		statuses: make(chan *Status),
		closer:   make(chan struct{}),
		errs:     make(chan error),
	}
}

// In returns a read-only channel of Commands which have been received from the
// client.
func (n *NetStream) In() <-chan Command { return n.in }

// Status returns a channel which, when written to, will write out the given
// Status command.
func (n *NetStream) Status() chan<- *Status { return n.statuses }

// Errs returns a read-only channel of errors encountered during the Listen
// operation.
func (n *NetStream) Errs() <-chan error { return n.errs }

// Close closes the Listen routine. Calling this function blocks until the
// Listen routine has entered a closing state. Should this function be called
// while a parse or send operation is taking place, then that operation will
// finish before the close operation takes place immediately afterwords.
func (n *NetStream) Close() { n.closer <- struct{}{} }

// Listen loops infinitely, managing the incoming and outgoing channel of chunks
// on the chunk stream shared between the server and client.
//
// Listen has three main goals:
//  - Parse incoming chunks, returning errors when they are unparsable.
//  - Serialize outgoing `onStatus` commands, returning an error when they are
//    either unserializable, or unwriteable.
//  - Respond to the `Close()` operation by closing all output channels.
//
// Listen runs within its own goroutine, and any errors encountered while
// running are sent over the internal errs channel, accessible from the `Errs()`
// function.
func (n *NetStream) Listen() {
	defer func() {
		close(n.statuses)
		close(n.in)
		close(n.errs)
		close(n.closer)
	}()

L:
	for {
		select {
		case chunk := <-n.chunks:
			cmd, err := n.parser.Parse(bytes.NewReader(chunk.Data))
			if err != nil {
				n.errs <- err
				continue
			}

			n.in <- cmd
		case st := <-n.statuses:
			c, err := st.AsChunk()
			if err != nil {
				n.errs <- err
				continue
			}

			if err = n.writer.Write(c); err != nil {
				n.errs <- err
			}
		case <-n.closer:
			break L
		}
	}
}
