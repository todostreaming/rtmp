package control

import "github.com/todostreaming/rtmp/chunk"

// Chunker is an interface representing a type responsible for tunring a RTMP
// control message into an RTMP chunk, capable of being sent over the network.
type Chunker interface {
	// Chunk marshals a Control sequence into an RTMP chunk, returning any
	// errors encountered along thw way as they come up.
	//
	// By specification, the RMTP chunks must be sent over ChunkStreamId
	// 0x2, and MessageStreamId 0x0. The TypeID of the MessageHeader must be
	// equivelant to the TypeId of the Control sequence.
	Chunk(Control) (*chunk.Chunk, error)
}
