package communication

import (
	"fmt"

	errors "github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/common"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/messages"
)

const PACKET_MIN_LEN = 1

type Packet struct {
	MessageCode uint8
	Message     messages.Message
}

func NewPacket(messageCode uint8, message messages.Message) Packet {
	return Packet{
		MessageCode: messageCode,
		Message:     message,
	}
}

func (p *Packet) ToBytes() []byte {
	messageBytes := p.Message.ToBytes()
	messageLen := uint16(len(messageBytes))

	totalLen := PACKET_MIN_LEN + messageLen
	bytes := make([]byte, 0, totalLen)

	bytes = append(bytes, byte(p.MessageCode))
	bytes = append(bytes, messageBytes...)

	return bytes
}

func PacketFromBytes(data []byte) (Packet, error) {
	if len(data) < PACKET_MIN_LEN {
		return Packet{}, fmt.Errorf(errors.PacketTooShortError)
	}
	offset := 0

	messageCode := int(data[offset])
	offset++

	var message messages.Message
	var err error
	switch messageCode {
	case messages.BET_CODE:
		message, err = messages.BetFromBytes(data[offset:])
	case messages.ASK_WINNERS_CODE:
		message, err = messages.AskWinnersFromBytes(data[offset:])
	case messages.WINNER_CODE:
		message, err = messages.BetFromBytes(data[offset:])
	case messages.FINISH_CODE:
		message, err = messages.FinishFromBytes(data[offset:])
	default:
		return Packet{}, fmt.Errorf(errors.UnexpectedMessage)
	}
	if err != nil {
		return Packet{}, err
	}

	return Packet{MessageCode: uint8(messageCode), Message: message}, nil
}
