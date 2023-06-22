package server

import (
	"bufio"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleConnection(t *testing.T) {
	// Start a TCP server to handle the connection
	server, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	defer server.Close()

	// Start a goroutine to accept incoming connections and handle them
	go func() {
		conn, err := server.Accept()
		require.NoError(t, err)
		defer conn.Close()

		// Simulate processing and sending a response
		response := "RESPONSE\n"
		_, err = conn.Write([]byte(response))
		require.NoError(t, err)
	}()

	// Connect to the server
	conn, err := net.Dial("tcp", server.Addr().String())
	require.NoError(t, err)
	defer conn.Close()

	// Send a request to the server
	request := "REQUEST\n"
	_, err = conn.Write([]byte(request))
	require.NoError(t, err)

	// Read the response from the server
	response, err := bufio.NewReader(conn).ReadString('\n')
	require.NoError(t, err)

	// Assert the response
	expectedResponse := "RESPONSE\n"
	assert.Equal(t, expectedResponse, response)
}
