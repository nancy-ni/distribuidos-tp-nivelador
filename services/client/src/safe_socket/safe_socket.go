package safe_socket

import (
	"fmt"
	"io"
)

//TODO: Complete with a short-read/short-write tolerant implementation

func SendAll(socket io.Writer, bytes []byte) error {
	totalSent := 0

	for totalSent < len(bytes) {
		n, err := socket.Write(bytes[totalSent:])
		if err != nil {
			return err
		}
		if n == 0 {
			return fmt.Errorf("Escritura devuelve 0 bytes")
		}
		totalSent += n
	}

	return nil
}

func RecvAll(socket io.Reader, size int) ([]byte, error) {
	buff := make([]byte, size)
	totalRead := 0

	for totalRead < size {
		n, err := socket.Read(buff[totalRead:])
		if err != nil {
			return nil, err
		}
		if n == 0 {
			return nil, fmt.Errorf("Lectura devuelve 0 bytes")
		}
		totalRead += n
	}

	return buff[:totalRead], nil
}
