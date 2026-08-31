package communication

import (
	"io"

	protocol "github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/common"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
)

func SendPacket(socket io.Writer, packet Packet) error {
	packetBytes := packet.ToBytes()
	packetLength := uint16(len(packetBytes))
	packetLengthBytes := protocol.Uint16ToBytes(packetLength)

	if err := safe_socket.SendAll(socket, packetLengthBytes); err != nil {
		return err
	}
	if err := safe_socket.SendAll(socket, packetBytes); err != nil {
		return err
	}
	return nil
}

func ReceivePacket(socket io.Reader) (Packet, error) {
	packetLengthBytes, err := safe_socket.RecvAll(socket, 2)
	if err != nil {
		return Packet{}, err
	}
	packetLength, err := protocol.BytesToUint16(packetLengthBytes)
	if err != nil {
		return Packet{}, err
	}

	packetBytes, err := safe_socket.RecvAll(socket, int(packetLength))
	if err != nil {
		return Packet{}, err
	}
	packet, err := PacketFromBytes(packetBytes)
	if err != nil {
		return Packet{}, err
	}

	return packet, nil
}
