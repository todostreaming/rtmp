package client_test

import (
	"bytes"
	"testing"

	"github.com/WatchBeam/rtmp/client"
	"github.com/stretchr/testify/assert"
)

func TestNewConstructsNewClients(t *testing.T) {
	b := new(bytes.Buffer)
	c, err := client.New(b)

	assert.IsType(t, &client.Client{}, c)
	assert.Nil(t, err)
	assert.Equal(t, b, c.Conn)
}
