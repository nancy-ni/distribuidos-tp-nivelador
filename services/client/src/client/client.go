package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
)

const CONNECTION_ATTEMPTS_MAX = 3
const CONNECTION_ATTEMPS_DELAY_MS = 200

const ECHO_CLIENT_BUFFER_SIZE = 512
const ECHO_CLIENT_MESSAGE_AMOUNT = 3
const ECHO_CLIENT_MESSAGE_DELAY_MS = 1000

type ClientConfig struct {
	ServerHost string
	ServerPort string
	AgencyId   string
}

type Client struct {
	conn   net.Conn
	config ClientConfig
}

func NewClient(config ClientConfig) (*Client, error) {
	conn, err := connectToServer(config.ServerHost, config.ServerPort)
	if err != nil {
		logger.Warn("connect-to-server", logger.Fail)
		return nil, err
	}

	client := &Client{conn: conn, config: config}
	return client, nil
}

func connectToServer(host, port string) (net.Conn, error) {
	const action = "connect-to-server"
	var err error
	var conn net.Conn

	logger.Info(action, logger.InProgress)
	for i := range CONNECTION_ATTEMPTS_MAX {
		conn, err = net.Dial("tcp", host+":"+port)
		if err != nil {
			logger.Warn(action, logger.Fail, "attempt", i)
			time.Sleep(CONNECTION_ATTEMPS_DELAY_MS * time.Millisecond)
			continue
		}

		logger.Info(action, logger.Success)
		break
	}

	return conn, err
}

func (client *Client) Run() error {
	const mainAction = "test-echo-server"
	defer client.conn.Close()

	inputFile, err := os.Open(os.Getenv("INPUT_FILE"))
	if err != nil {
		fmt.Println("Error abriendo el archivo input")
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(os.Getenv("OUTPUT_FILE"))
	if err != nil {
		fmt.Println("Error creando el archivo output")
		return err
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	messageId := 0
	for scanner.Scan() {
		clientMessage := scanner.Text()
		messageBytes := []byte(clientMessage)

		messageArgs := []any{"agency-id", client.config.AgencyId, "message-id", messageId}
		logger.Info(mainAction, logger.InProgress, messageArgs...)

		if err := safe_socket.SendAll(client.conn, messageBytes); err != nil {
			logger.Error("send-message", logger.Fail, messageArgs...)
			return err
		}

		responseBuffer, err := safe_socket.RecvAll(client.conn, len(messageBytes))
		if err != nil {
			logger.Error("recv-response", logger.Fail, messageArgs...)
			return err
		}

		responseString := string(responseBuffer)
		_, err = outputFile.WriteString(responseString + "\n")
		if err != nil {
			logger.Error("write-response", logger.Fail, messageArgs...)
			return err
		}

		if responseString != clientMessage {
			logger.Error("check-response", logger.Fail, messageArgs...)
			return err
		}

		messageId++
		time.Sleep(ECHO_CLIENT_MESSAGE_DELAY_MS * time.Millisecond)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error escaneando linea del archivo")
		return err
	}

	logger.Info(mainAction, logger.Success, "agency-id", client.config.AgencyId)

	return nil
}
