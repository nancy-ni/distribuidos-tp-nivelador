from protocol.common.utils import uint16_to_bytes, bytes_to_uint16
import safe_socket
from .packet import Packet, PACKET_LENGTH_BYTES

# Envia un Packet a traves de la conexion. Primero envia el largo del paquete en bytes,
# y posteriormente los bytes del paquete en si.
def send_packet(socket, packet):
    packet_bytes = packet.to_bytes()
    packet_length = len(packet_bytes)
    packet_length_bytes = uint16_to_bytes(packet_length)

    full_packet = packet_length_bytes + packet_bytes
    safe_socket.send_all(socket, full_packet)


# Recibe un Packet a traves de la conexion. Primero lee los bytes del largo del paquete (fijo),
# y posteriormente, utilizado el dato del largo, lee el paquete entero.
def receive_packet(socket) -> Packet:
    packet_length_bytes = safe_socket.recv_all(socket, PACKET_LENGTH_BYTES)
    packet_length = bytes_to_uint16(packet_length_bytes)
    packet_bytes = safe_socket.recv_all(socket, packet_length)

    return Packet.from_bytes(packet_bytes)