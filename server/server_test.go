package server_test

import (
	"net"
	"testing"

	"github.com/todostreaming/rtmp/client"
	"github.com/todostreaming/rtmp/server"
	"github.com/stretchr/testify/assert"
)

func TestNewServerConstructsServerWithValidBind(t *testing.T) {
	s, err := server.New("127.0.0.1:1234")
	defer s.Close()

	assert.IsType(t, &server.Server{}, s)
	assert.Nil(t, err)
}

func TestServerFailsWithInvalidBind(t *testing.T) {
	s, err := server.New("256.256.256.256:1234")

	assert.IsType(t, &server.Server{}, s)
	assert.NotNil(t, err)
}

func TestListenGetsNewClients(t *testing.T) {
	s, err := server.New("127.0.0.1:1935")
	assert.Nil(t, err)

	go s.Accept()
	defer s.Close()

	_, err = net.Dial("tcp", "127.0.0.1:1935")
	assert.Nil(t, err)

	assert.IsType(t, &client.Client{}, <-s.Clients())
}
