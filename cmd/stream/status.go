package stream

import (
	"github.com/todostreaming/amf0"
	"github.com/todostreaming/amf0/encoding"
	"github.com/todostreaming/rtmp/chunk"
)

const (
	// OnStatusChunkStreamId is the chunk stream ID that all OnStatus
	// commmands are to be sent over.
	OnStatusChunkStreamId uint32 = 5
	// OnStatusMessageStreamId is the message stream ID that all OnStatus
	// commands are to be sent over.
	OnStatusMessageStreamId uint32 = 1
	// OnStatusName is the name of the OnStatus command, as used in the
	// CommandHeader type.
	OnStatusName string = "onStatus"

	// Amf0CmdTypeId is the message type ID used to send the OnStatus
	// command in the chunk.
	Amf0CmdTypeId byte = 0x14
)

var (
	// OnStatusCommandHeader is a []byte containing a marshalled version of
	// the CommandHeader attached to all outgoing onStatus commands.
	OnStatusCommandHeader, _ = encoding.Marshal(&CommandHeader{
		Name: OnStatusName,
	})
)

// Status encapsulates the data contained in the body of an OnStatus command.
type Status struct {
	// Arguments correspond to the "arguments" field in the body of an
	// OnStatus command (as defined by the RTMP specification).
	Arguments amf0.Object
}

// NewStatus returns a new instance of the *Status type.
func NewStatus() *Status {
	return &Status{
		Arguments: *amf0.NewObject(),
	}
}

// Data marshals the data contained in the *Status type, returning either a
// []byte containing that data, or an error if it was unmarshallable.
func (s *Status) Data() ([]byte, error) {
	return encoding.Marshal(s)
}

// AsChunk formats the data contained in the entirety of the OnStatus command
// into a chunk, including the header. If the data was unable to be marshalled,
// then an error will be returned, otherwise a chunk will be returned in the
// happy case.
func (s *Status) AsChunk() (*chunk.Chunk, error) {
	body, err := s.Data()
	if err != nil {
		return nil, err
	}

	payload := append(OnStatusCommandHeader, body...)

	return &chunk.Chunk{
		Header: &chunk.Header{
			BasicHeader: chunk.BasicHeader{
				StreamId: OnStatusChunkStreamId,
			},
			MessageHeader: chunk.MessageHeader{
				Length:   uint32(len(payload)),
				TypeId:   Amf0CmdTypeId,
				StreamId: OnStatusMessageStreamId,
			},
		},
		Data: payload,
	}, nil
}
