package server

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"

	"github.com/sirupsen/logrus"

	"pow/internal/pkg/config"
	"pow/internal/pkg/pow"
	"pow/internal/pkg/protocol"
)

// Quotes - const array of quotes to respond on client's request
var Quotes = []string{
	"Stay hungry, stay foolish.",

	"Innovation distinguishes between a leader and a follower.",

	"Your work is going to fill a large part of your life," +
		" and the only way to be truly satisfied is to do what you believe is great work." +
		" And the only way to do great work is to love what you do.",
}

var ErrQuit = errors.New("client requests to close connection")

// Run - mainRun - main function, launches server to listen on given address
// and handle new connections function, launches server to listen on given address
// and handle new connections
func Run(ctx context.Context, address string, log *logrus.Logger) error {
	// Create a TCP listener.
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("error listening: %v", err)
		return err
	}
	defer listener.Close()

	log.Infof("listening %s", listener.Addr())

	for {
		// Accept incoming connections.
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("error accepting connection: %v", err)
			continue
		}

		// Handle the connection in a separate goroutine.
		go func() {
			defer conn.Close()

			err := handleConnection(ctx, conn, log)
			if err != nil {
				log.Errorf("error handling connection: %v", err)
			}
		}()
	}
}

func handleConnection(ctx context.Context, conn net.Conn, log *logrus.Logger) error {
	log.Infof("new client: %s", conn.RemoteAddr())
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		req, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading connection: %w", err)
		}
		msg, err := ProcessRequest(ctx, req, conn.RemoteAddr().String(), log)
		if err != nil {
			return fmt.Errorf("error processing request: %w", err)
		}
		if msg != nil {
			err := sendMsg(*msg, conn)
			if err != nil {
				return fmt.Errorf("error sending message: %w", err)
			}
		}
	}
}

func ProcessRequest(ctx context.Context, msgStr string, clientInfo string, log *logrus.Logger) (*protocol.Message, error) {
	msg, err := protocol.ParseMessage(msgStr)
	if err != nil {
		return nil, err
	}

	var responseMsg protocol.Message

	switch msg.Header {
	case protocol.Quit:
		return nil, ErrQuit
	case protocol.RequestChallenge:
		log.Infof("client %s requests challenge", clientInfo)
		responseMsg, err = createChallengeResponse(ctx, clientInfo)
		if err != nil {
			return nil, err
		}
	case protocol.RequestResource:
		log.Infof("client %s requests resource with payload", clientInfo)
		responseMsg, err = processResourceRequest(ctx, msg.Payload, clientInfo, log)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown header")
	}

	return &responseMsg, nil
}

func createChallengeResponse(ctx context.Context, clientInfo string) (protocol.Message, error) {
	conf := ctx.Value("config").(*config.Config)

	hashcash := pow.HashcashData{
		ZerosCount: conf.HashcashZerosCount,
		Resource:   clientInfo,
		Counter:    0,
	}

	hashcashMarshaled, err := json.Marshal(hashcash)
	if err != nil {
		return protocol.Message{}, fmt.Errorf("error marshaling hashcash: %v", err)
	}

	responseMsg := protocol.Message{
		Header:  protocol.ResponseChallenge,
		Payload: string(hashcashMarshaled),
	}

	return responseMsg, nil
}

func processResourceRequest(ctx context.Context, payload string, clientInfo string, log *logrus.Logger) (protocol.Message, error) {
	var hashcash pow.HashcashData
	err := json.Unmarshal([]byte(payload), &hashcash)
	if err != nil {
		return protocol.Message{}, fmt.Errorf("error unmarshaling hashcash: %w", err)
	}

	if hashcash.Resource != clientInfo {
		return protocol.Message{}, fmt.Errorf("invalid hashcash resource")
	}

	maxIter := hashcash.Counter
	if maxIter == 0 {
		maxIter = 1
	}

	_, err = hashcash.BruteForceHashcash(maxIter)
	if err != nil {
		return protocol.Message{}, fmt.Errorf("invalid hashcash")
	}

	log.Infof("client %s successfully computed hashcash", clientInfo)

	responseMsg := protocol.Message{
		Header:  protocol.ResponseResource,
		Payload: Quotes[rand.Intn(3)],
	}

	return responseMsg, nil
}

// sendMsg - send protocol message to connection
func sendMsg(msg protocol.Message, conn net.Conn) error {
	msgStr := fmt.Sprintf("%s\n", msg.ToString())
	_, err := conn.Write([]byte(msgStr))
	return err
}
