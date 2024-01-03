package client_test

import (
	"bytes"
	"testing"

	"github.com/todostreaming/rtmp/client"
	"github.com/stretchr/testify/assert"
)

func TestNewConstructsNewClients(t *testing.T) {
	b := new(bytes.Buffer)
	c := client.New(b)

	assert.IsType(t, &client.Client{}, c)
	assert.Equal(t, b, c.Conn)
}
