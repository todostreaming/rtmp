package stream

import "github.com/todostreaming/amf0"

// Type CommandHeader represents the command header belonging to commands shared
// by the NetStream and NetConnection.
//
// It holds data about the name, transaction ID, and arguments of a command sent
// over the NetStream from the client to the server. For all commands sent over
// the NetStream, the values for TransactionId and Arguments are 0 and nil,
// respectively.
//
// TODO(taylor): this sort of logic is shared between cmd/stream and cmd/conn,
// and should probably live in cmd, but that abstraction is tricky between
// reading and writing. Visit this later.
type CommandHeader struct {
	Name          string
	TransactionId float64
	Arguments     *amf0.Object
}

// Command is a tag-type for commands that may be received over the net stream.
type Command interface {
	IsCommand() bool
}
