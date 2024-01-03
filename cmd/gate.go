package cmd

import "github.com/todostreaming/rtmp/chunk"

// Gate is a single-function interfaces that provides infomration about whether
// a certain chan<- *chunk.Chunk is "open" to accept a Chunk. When wrapped over
// a channel, a "filtered channel" of sorts, is created.
//
// There are several simple implementations of this type below, as well as ways
// to combine multiple implementations.
type Gate interface {
	Open(*chunk.Chunk) bool
}

// StreamIdGate implements the Gate interface and filters chunks to be matching
// a certain ChunkStream ID.
type StreamIdGate struct {
	// StreamId is the ChunkStream ID that will be allowed through the
	// filter.
	StreamId uint32
}

var _ Gate = new(StreamIdGate)

// Open implements Gate.Open.
func (g *StreamIdGate) Open(c *chunk.Chunk) bool {
	return c.Header.BasicHeader.StreamId == g.StreamId
}

// MessageStreamGate implements the Gate interface and filters chunks to be
// matching a certain MessageStream ID.
type MessageStreamGate struct {
	// StreamId is the MessageStream ID that will be allowed through the
	// filter.
	StreamId uint32
}

var _ Gate = new(MessageStreamGate)

// Open implements Gate.Open.
func (g *MessageStreamGate) Open(c *chunk.Chunk) bool {
	return c.Header.MessageHeader.StreamId == g.StreamId
}

// TypeIdGate provides an implementation of the Gate interface, filtering chunks
// by their MessageHeader's TypeId.
type TypeIdGate struct {
	// TypeId is the TypeId of chunk's MessageHeaders that will be allowed
	// through the filter.
	TypeId byte
}

var _ Gate = new(TypeIdGate)

// Open implements Gate.Open.
func (g *TypeIdGate) Open(c *chunk.Chunk) bool {
	return c.Header.MessageHeader.TypeId == g.TypeId
}

// UnionGate is an implementation of the Gate type that essentially represents a
// logical AND. It is only open when all sub-gates are also open.
type UnionGate struct {
	gates []Gate
}

// NewUnionGate constructs a new instance of the UnionGate type, initializing
// all sub-gates to be the variadic gates argument passed.
func NewUnionGate(gates ...Gate) *UnionGate {
	return &UnionGate{gates: gates}
}

var _ Gate = new(UnionGate)

// Open implements Gate.Open. It is open noly when there is:
//
//  a) the amount of subgates `n` is at least 1
//  b) all of those `n` gates are open.
func (g *UnionGate) Open(c *chunk.Chunk) bool {
	for _, gate := range g.gates {
		if !gate.Open(c) {
			return false
		}
	}

	return len(g.gates) > 0
}

// AnyGate is an implementation of the Gate interface that essentially
// represents a logical OR. It is open when any of the sub-gates are open.
type AnyGate struct {
	gates []Gate
}

// NewAnyGate returns a new instance of the AnyGate type initialized with all
// sub-gates passed as variadic arguments.
func NewAnyGate(gates ...Gate) *AnyGate {
	return &AnyGate{gates: gates}
}

var _ Gate = new(AnyGate)

// Open implements Gate.Open, opening the AnyGate when:
//
//  a) there is at least one gate
//  b) any of those gates are open
func (g *AnyGate) Open(c *chunk.Chunk) bool {
	for _, gate := range g.gates {
		if gate.Open(c) {
			return true
		}
	}

	return false
}

var (
	// NetConnGate filters chunks to only those matching the NetConn type.
	NetConnGate = NewAnyGate(
		&StreamIdGate{3},
		NewUnionGate(&StreamIdGate{8}, &MessageStreamGate{0x0}),
	)

	// NetStreamGate filters chunks to only those matching the NetStream
	// type.
	NetStreamGate = NewUnionGate(
		NewAnyGate(
			&StreamIdGate{4},
			NewUnionGate(
				&StreamIdGate{8},
				&MessageStreamGate{0x1},
			),
		),
		&TypeIdGate{0x14},
	)

	// DataStreamGate filters chunks to only those matching the DataStream
	// type.
	DataStreamGate = NewUnionGate(&StreamIdGate{4}, NewAnyGate(
		&TypeIdGate{0x08}, &TypeIdGate{0x09}, &TypeIdGate{0x12},
	))
)
