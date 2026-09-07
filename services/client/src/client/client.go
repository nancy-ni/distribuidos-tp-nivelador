package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
	errors "github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/common"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/communication"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/messages"
)

const CONNECTION_ATTEMPTS_MAX = 3
const CONNECTION_ATTEMPS_DELAY_MS = 500

const BET_ELEMS_LEN = 5

type ClientConfig struct {
	ServerHost string
	ServerPort string
	AgencyId   string
}

type Client struct {
	conn           net.Conn
	config         ClientConfig
	isShuttingDown bool
	shutdownMutex  sync.Mutex
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

func (client *Client) setShutdown() {
	client.shutdownMutex.Lock()
	defer client.shutdownMutex.Unlock()
	client.isShuttingDown = true
}

func (client *Client) checkShutdown() bool {
	client.shutdownMutex.Lock()
	defer client.shutdownMutex.Unlock()
	return client.isShuttingDown
}

func (client *Client) handleShutdown(done chan struct{}) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer signal.Stop(signalChannel)
		select {
		case <-signalChannel:
			client.setShutdown()
			client.conn.Close()
		case <-done:
			return
		}
	}()
}

func (client *Client) Run() error {
	done := make(chan struct{})
	defer client.conn.Close()
	defer close(done)
	client.handleShutdown(done)

	if err := client.sendBets(); err != nil {
		logger.Error("send-bets", logger.Fail)
		return err
	}

	if err := client.sendWinnersRequest(); err != nil {
		logger.Error("ask-winners", logger.Fail)
		return err
	}

	if err := client.receiveWinners(); err != nil {
		logger.Error("recv-winners", logger.Fail)
		return err
	}

	logger.Info("test-echo-server", logger.Success, "agency-id", client.config.AgencyId)

	return nil
}

// Lee el archivo de input linea por linea, ensambla la apuesta de cada linea, y cuando
// se acumulan suficientes apuestas para formar un batch (o cuando se termina el archivo),
// se envia este ultimo al servidor. En caso de error, se chequea si el cliente estaba en
// un caso de shutdown, en lugar de un error de conexion.
func (client *Client) sendBets() error {
	inputFile, err := os.Open(os.Getenv("INPUT_FILE"))
	if err != nil {
		return fmt.Errorf(errors.OpenInputFileError)
	}
	defer inputFile.Close()

	agencyIdNumber, err := client.getAgencyIdNumber()
	if err != nil {
		return err
	}
	batchSize, err := client.getBatchSizeNumber()
	if err != nil {
		return err
	}

	batch := messages.NewBatch([]messages.Bet{})
	scanner := bufio.NewScanner(inputFile)
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
			if err := client.sendBatch(batch, messageId); err != nil {
				if client.checkShutdown() {
					return nil
				}
				return err
			}
			batch.Bets = []messages.Bet{}
			messageId++
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf(errors.ScanFileError)
	}

	if len(batch.Bets) > 0 {
		if err := client.sendBatch(batch, messageId); err != nil {
			if client.checkShutdown() {
				return nil
			}
			return err
		}
		batch.Bets = []messages.Bet{}
		messageId++
	}

	return nil
}

// Envia un batch y espera recibir un Ack de respuesta por parte del servidor.
func (client *Client) sendBatch(batch messages.Batch, messageId int) error {
	messageArgs := []any{"agency-id", client.config.AgencyId, "message-id", messageId}
	logger.Info("test-echo-server", logger.InProgress, messageArgs...)

	batchPacket := communication.NewPacket(messages.BATCH_CODE, &batch)
	if err := communication.SendPacket(client.conn, batchPacket); err != nil {
		logger.Error("send-batch", logger.Fail, messageArgs...)
		return err
	}
	if err := client.receiveAck(); err != nil {
		logger.Error("recv-ack", logger.Fail, messageArgs...)
		return err
	}
	return nil
}

func (client *Client) receiveAck() error {
	response, err := communication.ReceivePacket(client.conn)
	if err != nil {
		return err
	}
	if response.MessageCode != messages.ACK_CODE && response.MessageCode != messages.ERROR_CODE {
		return fmt.Errorf(errors.UnexpectedMessage)
	}
	if response.MessageCode == messages.ERROR_CODE {
		if errorMsg, ok := response.Message.(*messages.ErrorMessage); ok {
			return fmt.Errorf("%s", errorMsg.Reason)
		}
	}
	return nil
}

// Envia un mensaje AskWinners al servidor, indicando que la agencia esta lista para
// recibir los resultados del sorteo.
func (client *Client) sendWinnersRequest() error {
	agencyIdNumber, err := client.getAgencyIdNumber()
	if err != nil {
		return err
	}

	askWinners := messages.NewAskWinners(uint32(agencyIdNumber))
	askWinnersPacket := communication.NewPacket(messages.ASK_WINNERS_CODE, &askWinners)
	if err := communication.SendPacket(client.conn, askWinnersPacket); err != nil {
		return err
	}

	return nil
}

// Recibe las apuestas ganadoras del servidor, hasta que reciba un mensaje Finish que indica
// el fin del intercambio (o si ocurre un timeout). A los ganadores los registra en
// el path de output especificado.
func (client *Client) receiveWinners() error {
	outputFile, err := os.Create(os.Getenv("OUTPUT_FILE"))
	if err != nil {
		return fmt.Errorf(errors.OpenOutputFileError)
	}
	defer outputFile.Close()

	for {
		packet, err := communication.ReceivePacket(client.conn)
		if err != nil {
			if client.checkShutdown() {
				return nil
			}
			return err
		}
		if packet.MessageCode == messages.FINISH_CODE {
			return nil
		}
		if err := client.processPacket(packet, outputFile); err != nil {
			return err
		}
	}
}

func (client *Client) processPacket(packet communication.Packet, outputFile *os.File) error {
	switch packet.MessageCode {
	case messages.ERROR_CODE:
		if errorMsg, ok := packet.Message.(*messages.ErrorMessage); ok {
			return fmt.Errorf("%s", errorMsg.Reason)
		}
	case messages.WINNER_CODE:
		if bet, ok := packet.Message.(*messages.Bet); ok {
			_, err := outputFile.WriteString(bet.ToString() + "\n")
			if err != nil {
				logger.Error("save-winner", logger.Fail)
				return err
			}
		}
	}
	return nil
}

func assembleBet(betString string, agencyId uint32) (messages.Bet, error) {
	betData := strings.Split(betString, ",")
	if len(betData) != BET_ELEMS_LEN {
		return messages.Bet{}, fmt.Errorf(errors.AssembleBetError, betData)
	}
	firstName, lastName, dniString, birthday, betNumberString := betData[0], betData[1], betData[2], betData[3], betData[4]
	if len(firstName) > messages.FIRSTNAME_MAX_LEN || len(lastName) > messages.LASTNAME_MAX_LEN {
		return messages.Bet{}, fmt.Errorf(errors.AssembleBetError, betData)
	}
	dni, err := strconv.ParseUint(dniString, 10, messages.DNI_LEN_BYTES*8)
	if err != nil {
		return messages.Bet{}, fmt.Errorf(errors.AssembleBetError, betData)
	}
	betNumber, err := strconv.ParseUint(betNumberString, 10, messages.BET_NUMBER_LEN_BYTES*8)
	if err != nil {
		return messages.Bet{}, fmt.Errorf(errors.AssembleBetError, betData)
	}

	bet := messages.NewBet(agencyId, firstName, lastName, uint32(dni), birthday, uint16(betNumber))
	return bet, nil
}

func (client *Client) getAgencyIdNumber() (uint64, error) {
	agencyIdNumber, err := strconv.ParseUint(client.config.AgencyId, 10, messages.AGENCY_ID_LEN_BYTES*8)
	if err != nil {
		return 0, fmt.Errorf(errors.InvalidAgencyIdError)
	}
	return agencyIdNumber, nil
}

func (client *Client) getBatchSizeNumber() (int, error) {
	batchSizeStr := os.Getenv("BATCH_SIZE")
	if batchSizeStr == "" {
		return 0, fmt.Errorf(errors.InvalidBatchSizeError)
	}
	batchSize, err := strconv.Atoi(batchSizeStr)
	if err != nil {
		return 0, fmt.Errorf(errors.InvalidBatchSizeError)
	}
	return batchSize, nil
}
