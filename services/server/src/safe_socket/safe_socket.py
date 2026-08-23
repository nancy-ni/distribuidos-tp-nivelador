import socket

# TODO: Complete with a short-read/short-write tolerant implementation


def recv_all(socket: socket.socket, size):
    buff = b""
    while len(buff) < size:
        received = socket.recv(size - len(buff))
        if not received:
            if not buff:
                return b""
            raise ConnectionError("Hubo un problema con la conexion")

        buff += received

    return buff


def send_all(socket: socket.socket, bytes):
    total_sent = 0
    while total_sent < len(bytes):
        n = socket.send(bytes[total_sent:])
        if n == 0:
            raise ConnectionError("Hubo un problema con la conexion")

        total_sent += n
    return total_sent
