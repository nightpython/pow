package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/sirupsen/logrus"

	"pow/config"
	"pow/internal/hashcash"
	"pow/internal/protocol"
)

// Run is the main function that launches the client to connect and work with the server at the specified address.
func Run(ctx context.Context, address string, log *logrus.Logger) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Errorf("error connecting to %s: %v", address, err)
		return err
	}

	log.Infof("connected to %s", address)

	defer conn.Close()

	// The client will send a new request every 10 seconds endlessly.
	for {
		message, err := HandleConnection(ctx, conn, conn, log)
		if err != nil {
			log.Errorf("error handling connection: %v", err)
			return err
		}
		log.Infof("quote result: %s", message)
		time.Sleep(10 * time.Second)
	}
}

func HandleConnection(ctx context.Context, readerConn io.Reader, writerConn io.Writer, log *logrus.Logger) (string, error) {
	reader := bufio.NewReader(readerConn)

	// Step 1: Requesting challenge
	err := sendMsg(protocol.Message{
		Header: protocol.RequestChallenge,
	}, writerConn)
	if err != nil {
		log.Errorf("error sending request: %v", err)
		return "", fmt.Errorf("error send request: %w", err)
	}

	// Step 2: Reading and parsing response
	msgStr, err := readConnMsg(reader)
	if err != nil {
		log.Errorf("error reading message: %v", err)
		return "", fmt.Errorf("error read msg: %w", err)
	}
	msg, err := protocol.ParseMessage(msgStr)
	if err != nil {
		log.Errorf("error parsing message: %v", err)
		return "", fmt.Errorf("error parse msg: %w", err)
	}
	var hashcash hashcash.HashcashData
	err = json.Unmarshal([]byte(msg.Payload), &hashcash)
	if err != nil {
		log.Errorf("error parsing hashcash: %v", err)
		return "", fmt.Errorf("error parse hashcash: %w", err)
	}
	log.Infof("got hashcash: %+v", hashcash)

	// Step 3: Compute hashcash
	conf := ctx.Value("config").(*config.Config)
	hashcash, err = hashcash.BruteForceHashcash(conf.HashcashMaxIterations)
	if err != nil {
		log.Errorf("error computing hashcash: %v", err)
		return "", fmt.Errorf("error compute hashcash: %w", err)
	}
	log.Infof("hashcash computed: %+v", hashcash)
	// Marshal solution to JSON
	byteData, err := json.Marshal(hashcash)
	if err != nil {
		log.Errorf("error marshaling hashcash: %v", err)
		return "", fmt.Errorf("error marshal hashcash: %w", err)
	}

	// Step 4: Send challenge solution back to the server
	err = sendMsg(protocol.Message{
		Header:  protocol.RequestSolution,
		Payload: string(byteData),
	}, writerConn)
	if err != nil {
		log.Errorf("error sending request: %v", err)
		return "", fmt.Errorf("error send request: %w", err)
	}
	log.Infof("challenge sent to server")

	// Step 5: Get result quote from the server
	msgStr, err = readConnMsg(reader)
	if err != nil {
		log.Errorf("error reading message: %v", err)
		return "", fmt.Errorf("error read msg: %w", err)
	}
	msg, err = protocol.ParseMessage(msgStr)
	if err != nil {
		log.Errorf("error parsing message: %v", err)
		return "", fmt.Errorf("error parse msg: %w", err)
	}
	return msg.Payload, nil
}

// readConnMsg reads a string message from the connection.
func readConnMsg(reader *bufio.Reader) (string, error) {
	return reader.ReadString('\n')
}

// sendMsg sends a protocol message to the connection.
func sendMsg(msg protocol.Message, conn io.Writer) error {
	msgStr := fmt.Sprintf("%s\n", msg.ToString())
	_, err := conn.Write([]byte(msgStr))
	return err
}
