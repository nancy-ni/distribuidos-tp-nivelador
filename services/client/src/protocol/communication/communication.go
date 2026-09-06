package communication

import (
	"net"
	"time"

	common "github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/common"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
)

func SendPacket(socket net.Conn, packet Packet) error {
	if err := socket.SetWriteDeadline(time.Now().Add(common.SOCKET_WRITE_TIMEOUT_SEC * time.Second)); err != nil {
		return err
	}

	packetBytes := packet.ToBytes()
	packetLength := uint16(len(packetBytes))
	packetLengthBytes := common.Uint16ToBytes(packetLength)

	fullPacket := append(packetLengthBytes, packetBytes...)
	if err := safe_socket.SendAll(socket, fullPacket); err != nil {
		return err
	}
	return nil
}

func ReceivePacket(socket net.Conn) (Packet, error) {
	if err := socket.SetReadDeadline(time.Now().Add(common.SOCKET_READ_TIMEOUT_SEC * time.Second)); err != nil {
		return Packet{}, err
	}

	packetLengthBytes, err := safe_socket.RecvAll(socket, PACKET_LENGTH_BYTES)
	if err != nil {
		return Packet{}, err
	}
	packetLength, err := common.BytesToUint16(packetLengthBytes)
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
