package conn

import "github.com/WatchBeam/amf0"

// Type Sendable is used to tag certain Responses as being able to be sent.
type Sendable interface {
	CanSend() bool
}

type CreateStreamResponse struct {
	TransactionId float64
	CmdObject     amf0.Object
	StreamID      float64
}

type ConnectResponse struct {
	TransactionId float64
	Properties    amf0.Object
	Information   amf0.Object
}

func (_ *CreateStreamResponse) CanSend() bool { return true }
func (_ *ConnectResponse) CanSend() bool      { return true }
