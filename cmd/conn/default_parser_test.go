package conn_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/todostreaming/amf0"
	"github.com/todostreaming/rtmp/cmd/conn"
	"github.com/stretchr/testify/assert"
)

var (
	CreatePayload = []byte{
		0x00, 0x40, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05,
	}
)

func TestNewParserMakesNewParser(t *testing.T) {
	p := conn.NewParser(map[string]conn.ReceviableFactory{})

	assert.IsType(t, new(conn.SimpleParser), p)
}

func TestParseParsesReadersIntoReceivables(t *testing.T) {
	p := conn.NewParser(map[string]conn.ReceviableFactory{
		"createStream": func() conn.Receivable {
			return new(conn.CreateStreamCommand)
		},
	})

	r, err := p.Parse(
		amf0.NewString("createStream"),
		bytes.NewReader(CreatePayload),
	)

	assert.Nil(t, err)

	switch typ := r.(type) {
	case *conn.CreateStreamCommand:
		assert.Equal(t, float64(4), typ.TransactionId)
		assert.Nil(t, typ.Metadata)
	default:
		t.Fatalf("rtmp/cmd/conn: unknown type %T", typ)
	}
}

func TestParseFailsWhenNoMatchingCommandIsPresent(t *testing.T) {
	p := conn.NewParser(map[string]conn.ReceviableFactory{
	// No types
	})

	r, err := p.Parse(amf0.NewString("not-a-type"), new(bytes.Buffer))

	assert.Nil(t, r)
	assert.Equal(t, "rtmp/cmd/conn: unknown command name: not-a-type",
		err.Error())
}

func TestParseFailsWhenDataIsUnmarshalable(t *testing.T) {
	p := conn.NewParser(map[string]conn.ReceviableFactory{
		"createStream": func() conn.Receivable {
			return new(conn.CreateStreamCommand)
		},
	})

	// Error: use an empty buffer to ensure that an EOF is thrown
	r, err := p.Parse(amf0.NewString("createStream"), new(bytes.Buffer))

	assert.Nil(t, r)
	assert.Equal(t, io.EOF, err)
}
