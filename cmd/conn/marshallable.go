package conn

import (
	"github.com/todostreaming/amf0"
	"github.com/todostreaming/amf0/encoding"
)

const (
	// SuccessfulResponseType is the respnse type string attached to
	// successful responses.
	SuccessfulResponseType = "_result"
)

// Type Marshallable is used to tag certain Responses as being able to be sent.
type Marshallable interface {
	Marshal() ([]byte, error)
}

type CreateStreamResponse struct {
	ResponseType  string
	TransactionId float64
	_             *amf0.Null
	StreamID      float64
}

type ConnectResponse struct {
	ResponseType  string
	TransactionId float64
	Properties    amf0.Object
	Information   amf0.Object
}

// Marshal implements Marshallable.Marshal.
func (r *CreateStreamResponse) Marshal() ([]byte, error) {
	r.ResponseType = SuccessfulResponseType
	return encoding.Marshal(r)
}

// Marshal implements Marshallable.Marshal.
func (r *ConnectResponse) Marshal() ([]byte, error) {
	r.ResponseType = SuccessfulResponseType
	return encoding.Marshal(r)
}
