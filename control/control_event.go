package control

import (
	"encoding/binary"
	"io"
	"io/ioutil"
)

// EventType is a const type that wraps the uint16 base type to represent a
// particular sub-type of the Event message sent over the control sequence chunk
// stream.
type EventType uint16

const (
	SetBufferLength EventType = 3
)

// Event encapsulates any event that is sent over the control stream.
//
// NOTE: this is a temporary type, eventually it will be replaced with
// individual types implementing both Control and Event interfaces. The current
// structure of the parser/identifier does not easily allow for this, so some
// more heavy lifting will have to occur before this gets further attention.
type Event struct {
	// Type is the event type of the event.
	Type EventType
	// Body is the body payload of the event.
	Body []byte
}

var _ Control = new(Event)

// Read implements the Event.Read function, returning any errors that it
// encounters, or nil if the read was successful.
func (e *Event) Read(r io.Reader) error {
	var tid uint16
	if err := binary.Read(r, binary.BigEndian, &tid); err != nil {
		return err
	}

	e.Type = EventType(tid)

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	e.Body = body

	return nil
}

// Write implements the Event.Write function, returning any errors that it
// encounters, or nil if the write was successful.
func (e *Event) Write(w io.Writer) error {
	if err := binary.Write(
		w, binary.BigEndian, uint16(e.Type),
	); err != nil {
		return err
	}

	if _, err := w.Write(e.Body); err != nil {
		return err
	}

	return nil
}

// TypeId implements Event.TypeId and returns the TypeId of an Event control
// sequence.
func (e *Event) TypeId() byte { return 0x04 }
