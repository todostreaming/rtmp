package stream_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/todostreaming/rtmp/cmd/stream"
	"github.com/stretchr/testify/assert"
)

func TestNewParserRetrunsNewParsers(t *testing.T) {
	p := stream.NewParser(map[string]stream.CommandFactory{})

	assert.IsType(t, new(stream.SimpleParser), p)
}

func TestParserParsesCommands(t *testing.T) {
	p := stream.DefaultParser

	cmd, err := p.Parse(bytes.NewReader([]byte{
		0x02, 0x00, 0x07, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68,
		0x00, 0x40, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05,
		0x02, 0x00, 0x03, 0x66, 0x6f, 0x6f,
		0x02, 0x00, 0x04, 0x6c, 0x69, 0x76, 0x65,
	}))

	assert.Nil(t, err)
	assert.Equal(t, &stream.CommandPublish{
		Name: "foo",
		Type: "live",
	}, cmd)
}

func TestParserReturnsParsingErrors(t *testing.T) {
	p := stream.DefaultParser

	cmd, err := p.Parse(bytes.NewReader([]byte{
	// io.EOF
	}))

	assert.Equal(t, err, io.EOF)
	assert.Nil(t, cmd)
}

func TestParserReturnsErrorsWhenNoMatchingTypeIsFound(t *testing.T) {
	p := stream.DefaultParser

	cmd, err := p.Parse(bytes.NewReader([]byte{
		0x02, 0x00, 0x03, 0x66, 0x6f, 0x6f, 0x00, 0x40, 0x14, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x05,
	}))

	assert.Nil(t, cmd)
	assert.Equal(t, "cmd/stream: unknown NetStream command foo", err.Error())
}

func TestParserReturnsBodyErrors(t *testing.T) {
	p := stream.DefaultParser

	cmd, err := p.Parse(bytes.NewReader([]byte{
		// <CorrectStatusHeader>
		0x02, 0x00, 0x07, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68,
		0x00, 0x40, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05,
		// </CorrectStatusHeader>
		// <EmptyBody />
	}))

	assert.Nil(t, cmd)
	assert.Equal(t, io.EOF, err)
}
