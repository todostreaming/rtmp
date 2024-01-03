package cmd

import (
	"github.com/todostreaming/rtmp/chunk"
	"github.com/todostreaming/rtmp/cmd/conn"
	"github.com/todostreaming/rtmp/cmd/data"
	"github.com/todostreaming/rtmp/cmd/stream"
)

// Manager sits in front of all sub-packages of `cmd` and cleans up incoming
// chunks coming over streams 3, 4, and 5 into their appropriate spots, and
// uses the Gate mechanism to dispatch them appropriately to each sub-package.
type Manager struct {
	// chunks is the incoming chunk stream to feed from. In most normal
	// cases, this will be a *chunk.MultiStream, but either works.
	chunks chunk.Stream
	// closer is a channel which is written to when it is time to close the
	// Manager.
	closer chan struct{}

	// channels maps Gates to the channel which they are gating.
	channels map[Gate]chan<- *chunk.Chunk

	// netConn is the NetConnection which is engaged with the connecting
	// client.
	netConn *conn.NetConn
	// dataStream is the *data.Stream which is engaged with the connecting
	// client.
	dataStream *data.Stream
	// netStream is the NetStream which is engaged with the connecting
	// client.
	netStream *stream.NetStream
}

// New returns a new instance of the *Manager type. It takes in an incoming
// chunk stream, as well as a *chunk.Writer to read and write chunks from,
// respectively.
//
// During initialization, this function instantiates new instances of the
// NetConn, NetStream, and DataStream types, and attaches them the appropriate
// channels (also instantiated by this function).
//
// All internal channels (including the closer) are instantiated at this time,
// as well.
func New(chunks chunk.Stream, writer chunk.Writer) *Manager {
	netConnChunks := make(chan *chunk.Chunk)
	dataStreamChunks := make(chan *chunk.Chunk)
	netStreamChunks := make(chan *chunk.Chunk)

	return &Manager{
		chunks: chunks,
		closer: make(chan struct{}),

		channels: map[Gate]chan<- *chunk.Chunk{
			NetConnGate:    netConnChunks,
			DataStreamGate: dataStreamChunks,
			NetStreamGate:  netStreamChunks,
		},

		netConn:    conn.NewNetConnection(netConnChunks, writer),
		dataStream: data.NewStream(dataStreamChunks, writer),
		netStream:  stream.New(netStreamChunks, writer),
	}
}

// NetConn returns the NetConnection that is associated with this client.
func (m *Manager) NetConn() *conn.NetConn { return m.netConn }

//  NetStream returns the NetStream that is associated with this client.
func (m *Manager) NetStream() *stream.NetStream { return m.netStream }

// DataStream returns the DataStream that is associated with this client.
func (m *Manager) DataStream() *data.Stream { return m.dataStream }

func (m *Manager) Close() { m.closer <- struct{}{} }

// Dispatch handles the dispatch loop responsible for processing all incoming
// chunks that are received over the given chunk.Stream (see `New()`).
//
// Dispatch has two definite, and one optional responsibility:
//   1) Manage children. If `manageChildren` is passed as `true`, then the
//   lifecycles of the NetConn, NetStream, and DataStream will be managed.
//   "Managed," in this sense, means they are started and stopped at the
//   appropriate times. They are started at the same time this loop is started,
//   and they are stopped when the loop is terminated.
//
//   2) Respond to incoming chunks. To do this, each incoming chunk is read, and
//   then matched against all gates. If a gate is open for that particular
//   chunk, then it is dispatched over the corresponding channel. A chunk may be
//   distributed more than once, but in most cases, the set of channels given is
//   mutually exclusive.
//
//   3) Respond to the `Close()` operation. If close is passed, then the loop
//   will terminate and, if manageChildren is set to true, the children will be
//   closed as well.
//
// Dispatch runs within its own goroutine.
func (m *Manager) Dispatch(manageChildren bool) {
	if manageChildren {
		m.startChildren()
		defer m.cleanupChildren()
	}

	defer func() {
		close(m.closer)

		for _, ch := range m.channels {
			close(ch)
		}
	}()

	for {
		select {
		case c := <-m.chunks.In():
			for gate, chunks := range m.channels {
				if gate.Open(c) {
					chunks <- c
				}
			}
		case <-m.closer:
			break
		}
	}
}

// startChildren spawns all of the `Listen` subroutines for each managed child.
func (m *Manager) startChildren() {
	go m.netConn.Listen()
	go m.netStream.Listen()
	go m.dataStream.Recv()
}

// cleanupChildren stops all of the `Listen` subroutines for each managed child.
func (m *Manager) cleanupChildren() {
	m.netConn.Close()
	m.netStream.Close()
	m.dataStream.Close()
}
