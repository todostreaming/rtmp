package conn

import "github.com/WatchBeam/amf0"

// Type Sendable is used to tag certain Responses as being able to be sent.
type Sendable interface {
	CanSend() bool

	// HACK: temporary fix until amf0 knows how to deserialize embedded types.
	OnPreSend()
}

type CreateStreamResponse struct {
	ResponseType  string
	TransactionId float64
	CmdObject     amf0.Object
	StreamID      float64
}

type ConnectResponse struct {
	ResponseType  string
	TransactionId float64
	Properties    amf0.Object
	Information   amf0.Object
}

func (_ *CreateStreamResponse) CanSend() bool { return true }
func (_ *ConnectResponse) CanSend() bool      { return true }

func (r *CreateStreamResponse) OnPreSend() { r.ResponseType = "_result" }
func (r *ConnectResponse) OnPreSend()      { r.ResponseType = "_result" }
