package conn

import "github.com/todostreaming/amf0"

// Type Receviable is used to tag certain Commands as being able to be received.
type Receivable interface {
	CanReceive() bool
}

type ConnectCommand struct {
	TransactionId float64
	Metadata      *amf0.Object
}

type CreateStreamCommand struct {
	TransactionId float64
	Metadata      *amf0.Object
}

type ReleaseCommand struct {
	TransactionId float64
	Nil           *amf0.Null
	StreamKey     string
}

type FCPublishCommand struct {
	TransactionId float64
	Nil           *amf0.Null
	StreamKey     string
}

type GetStreamLength struct {
	StreamId float64
	Nil      *amf0.Null
	PlayPath string
}

func (_ *ConnectCommand) CanReceive() bool      { return true }
func (_ *CreateStreamCommand) CanReceive() bool { return true }
func (_ *ReleaseCommand) CanReceive() bool      { return true }
func (_ *FCPublishCommand) CanReceive() bool    { return true }
func (_ *GetStreamLength) CanReceive() bool     { return true }
