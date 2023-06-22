package protocol

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Quit              = iota // Quit indicates that either the server or client should close the connection.
	RequestChallenge         // RequestChallenge is a message from the client to the server, requesting a new challenge.
	ResponseChallenge        // ResponseChallenge is a message from the server to the client, containing a challenge for the client.
	RequestResource          // RequestResource is a message from the client to the server, containing a solved challenge.
	ResponseResource         // ResponseResource is a message from the server to the client, containing useful information if the solution is correct, or an error if it is not.
)

// Message represents a message structure used by both the server and client.
type Message struct {
	Header  int    // Header represents the type of the message, indicating its purpose or action.
	Payload string // Payload contains the data associated with the message, which can be in JSON format, a quote, or empty.
}

func (m *Message) ToString() string {
	return fmt.Sprintf("%d|%s", m.Header, m.Payload)
}

// ParseMessage parses the message from a string representation, checking the header and payload.
// It returns a pointer to the parsed Message struct or an error if the message doesn't match the protocol.
func ParseMessage(str string) (*Message, error) {
	str = strings.TrimSpace(str)
	parts := strings.SplitN(str, "|", 2)
	if len(parts) < 1 || len(parts) > 2 || (len(parts) == 2 && parts[0] == "") {
		return nil, fmt.Errorf("message doesn't match protocol")
	}
	msgType, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("cannot parse header")
	}
	msg := &Message{
		Header:  msgType,
		Payload: "",
	}
	if len(parts) == 2 {
		msg.Payload = parts[1]
	}
	return msg, nil
}
