package conn

import (
	"bytes"
	"fmt"

	"github.com/todostreaming/amf0"
	"github.com/todostreaming/rtmp/chunk"
)

const (
	// ChunkStreamId is the ID of the chunk stream that the net/conn package
	// listens to.
	ChunkStreamId uint32 = 3
)

// NetConn is an implementation of the NetConnection type as defined in the RTMP
// specification. It is capable of both reading and writing types that are
// able to be sent by the NetConnection section of the document (see RTMP
// specification).
//
// Both an In() and an Out() channel are exposed to read and write from the
// NetConnection.
type NetConn struct {
	// chunkStream is the incoming channel of Chunks.
	chunkStream <-chan *chunk.Chunk

	// parser is the Parser used to decode Receivables from the client,
	// before they are passed into the in channel.
	parser Parser
	// in is a channel of Receivables received from the client.
	in chan Receivable

	// writer is the chunk.Writer that is used to write data back into the
	// chunk stream.
	writer chunk.Writer
	// chunker is the chunker responsible for turning Marshallables into
	// chunks.
	chunker Chunker
	// out is a channel written to by owners of this type when they want to
	// send something over the channel.
	out chan Marshallable

	// errs is a channel which is written to when an error occurs.
	errs chan error
	// closer is a channel written to when the Listen operation should halt.
	closer chan struct{}
}

// NewNetConnection returns a new instance of the NetConn type initialized with
// the given input channel of *chunk.Chunks, as well as the chunk.Writer. All
// internal channels are initialzied, but the Listen method is not called.
func NewNetConnection(chunks <-chan *chunk.Chunk, writer chunk.Writer) *NetConn {
	return &NetConn{
		parser:      DefaultParser,
		chunkStream: chunks,
		writer:      writer,
		chunker:     NewChunker(ChunkStreamId),
		in:          make(chan Receivable),
		out:         make(chan Marshallable),
		errs:        make(chan error),
		closer:      make(chan struct{}),
	}
}

// In returns a read-only channel of Receivables. This channel is written to
// when a Receivable is read from the connected client.
func (n *NetConn) In() <-chan Receivable { return n.in }

// Out returns a write-only channel of Marshallables. It should be written to
// when one wants to send a Marshallable back to the client.
func (n *NetConn) Out() chan<- Marshallable { return n.out }

// Close halts the Listen operation after the current item has finished
// processing.
func (n *NetConn) Close() { n.closer <- struct{}{} }

// Errs returns a read-only channel of errors which were encountered during the
// Listen operation (see below).
func (n *NetConn) Errs() <-chan error { return n.errs }

// Listen monitors all of the ingoing and outgoing chnanels on the NetConn type
// and makes sure that things are in order.
//  - It decodes chunks when they are received into Receivables, passing them
//    along the In() channel, or writing an error to Errs() if a parse error was
//    encountered.
//
//  - It chunks outgoing messages written to the Out() channel, and sends them
//    over the chunk stream, writing an error to Errs() if one was encountered.
//
// Listen terminates when the closer channel can be read (accomplished by
// calling Close()).
//
// Listen runs within its own goroutine.
func (n *NetConn) Listen() {
	for {
		select {
		case c := <-n.chunkStream:
			buf := bytes.NewBuffer(c.Data)

			name, err := amf0.Decode(buf)
			if err != nil {
				n.errs <- err
				continue
			}

			nameStr, ok := name.(*amf0.String)
			if !ok {
				n.errs <- fmt.Errorf("rtmp/conn: wrong type for AMF header: %T (expected amf0.String)", name)
				continue
			}

			if r, err := n.parser.Parse(nameStr, buf); err != nil {
				n.errs <- err
			} else {
				n.in <- r
			}
		case out := <-n.out:
			c, err := n.chunker.Chunk(out)
			if err != nil {
				n.errs <- err
				continue
			}

			if err = n.writer.Write(c); err != nil {
				n.errs <- err
			}
		case <-n.closer:
			return
		}
	}
}
