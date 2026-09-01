package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
	errors "github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/common"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/communication"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/messages"
)

const CONNECTION_ATTEMPTS_MAX = 3
const CONNECTION_ATTEMPS_DELAY_MS = 200

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
	defer client.conn.Close()

	if err := client.sendBets(); err != nil {
		return err
	}

	if err := client.sendWinnersRequest(); err != nil {
		return err
	}

	if err := client.receiveWinners(); err != nil {
		return err
	}

	logger.Info("test-echo-server", logger.Success, "agency-id", client.config.AgencyId)

	return nil
}

func (client *Client) sendBets() error {
	agencyIdNumber, err := strconv.ParseUint(client.config.AgencyId, 10, 32)
	if err != nil {
		return fmt.Errorf(errors.InvalidAgencyIdError)
	}

	inputFile, err := os.Open(os.Getenv("INPUT_FILE"))
	if err != nil {
		return fmt.Errorf(errors.OpenInputFileError)
	}
	defer inputFile.Close()

	batchSizeStr := os.Getenv("BATCH_SIZE")
	if batchSizeStr == "" {
		return fmt.Errorf(errors.InvalidBatchSizeError)
	}
	batchSize, err := strconv.Atoi(batchSizeStr)
	if err != nil {
		return fmt.Errorf(errors.InvalidBatchSizeError)
	}

	scanner := bufio.NewScanner(inputFile)
	batch := messages.NewBatch([]messages.Bet{})
	messageId := 0
	for scanner.Scan() {
		betString := scanner.Text()
		messageArgs := []any{"agency-id", client.config.AgencyId, "message-id", messageId}

		bet, err := assembleBet(betString, uint32(agencyIdNumber))
		if err != nil {
			logger.Error("assemble-bet", logger.Fail, messageArgs...)
			continue
		}

		batch.Bets = append(batch.Bets, bet)
		if len(batch.Bets) == batchSize {
			batchPacket := communication.NewPacket(messages.BATCH_CODE, &batch)
			logger.Info("test-echo-server", logger.InProgress, messageArgs...)

			if err := communication.SendPacket(client.conn, batchPacket); err != nil {
				logger.Error("send-message", logger.Fail, messageArgs...)
				return err
			}
			batch.Bets = []messages.Bet{}
			messageId++
		}

	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf(errors.ScanFileError)
	}

	messageArgs := []any{"agency-id", client.config.AgencyId, "message-id", messageId}
	if len(batch.Bets) > 0 {
		batchPacket := communication.NewPacket(messages.BATCH_CODE, &batch)
		logger.Info("test-echo-server", logger.InProgress, messageArgs...)

		if err := communication.SendPacket(client.conn, batchPacket); err != nil {
			logger.Error("send-message", logger.Fail, messageArgs...)
			return err
		}

		batch.Bets = []messages.Bet{}
		messageId++
	}

	return nil
}

func (client *Client) sendWinnersRequest() error {
	agencyIdNumber, err := strconv.ParseUint(client.config.AgencyId, 10, 32)
	if err != nil {
		return fmt.Errorf(errors.InvalidAgencyIdError)
	}

	askWinners := messages.NewInquirie(uint32(agencyIdNumber))
	askWinnersPacket := communication.NewPacket(messages.ASK_WINNERS_CODE, &askWinners)
	if err := communication.SendPacket(client.conn, askWinnersPacket); err != nil {
		logger.Error("ask-winners", logger.Fail)
		return err
	}

	return nil
}

func (client *Client) receiveWinners() error {
	outputFile, err := os.Create(os.Getenv("OUTPUT_FILE"))
	if err != nil {
		return fmt.Errorf(errors.OpenOutputFileError)
	}
	defer outputFile.Close()

	for {
		// TODO: agregar timeout socket
		packet, err := communication.ReceivePacket(client.conn)
		if err != nil {
			logger.Error("recv-response", logger.Fail)
			return err
		}
		if packet.MessageCode == messages.FINISH_CODE {
			break
		}
		if packet.MessageCode != messages.WINNER_CODE {
			logger.Error("recv-response", logger.Fail)
			continue
		}

		switch message := packet.Message.(type) {
		case *messages.Bet:
			_, err := outputFile.WriteString(message.ToString() + "\n")
			fmt.Println("RECIBI WINNER CORRECTAMENTE")
			if err != nil {
				logger.Error("write-response", logger.Fail)
				continue
			}
		default:
			logger.Error("recv-response", logger.Fail)
			continue
		}
	}

	return nil
}

func assembleBet(betString string, agencyId uint32) (messages.Bet, error) {
	betData := strings.Split(betString, ",")
	if len(betData) != 5 {
		return messages.Bet{}, fmt.Errorf(errors.AssembleBetError, betData)
	}
	firstName, lastName, dniString, birthday, betNumberString := betData[0], betData[1], betData[2], betData[3], betData[4]
	dni, err := strconv.ParseUint(dniString, 10, 32)
	if err != nil {
		return messages.Bet{}, fmt.Errorf(errors.AssembleBetError, betData)
	}
	betNumber, err := strconv.ParseUint(betNumberString, 10, 16)
	if err != nil {
		return messages.Bet{}, fmt.Errorf(errors.AssembleBetError, betData)
	}

	bet := messages.NewBet(agencyId, firstName, lastName, uint32(dni), birthday, uint16(betNumber))
	return bet, nil
}
