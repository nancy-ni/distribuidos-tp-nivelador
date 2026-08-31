from protocol.common.utils import uint16_to_bytes, bytes_to_uint16
import safe_socket
from .packet import Packet

def send_packet(socket, packet):
    try:
        packet_bytes = packet.to_bytes()
        packet_length = len(packet_bytes)
        packet_length_bytes = uint16_to_bytes(packet_length)

        safe_socket.send_all(socket, packet_length_bytes)
        safe_socket.send_all(socket, packet_bytes)
    
        return None
    except Exception as e:
        return e


def receive_packet(socket) -> Packet:
    try:
        packet_length_bytes = safe_socket.recv_all(socket, 2)
        packet_length = bytes_to_uint16(packet_length_bytes)

        packet_bytes = safe_socket.recv_all(socket, packet_length)
        packet, err = Packet.from_bytes(packet_bytes)
        if err:
            return None, err

        return packet, None
    except Exception as e:
        return None, e