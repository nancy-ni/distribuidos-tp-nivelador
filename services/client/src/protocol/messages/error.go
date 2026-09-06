package messages

import (
	"fmt"

	common "github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/common"
)

type ErrorMessage struct {
	Reason string
}

func (e *ErrorMessage) ToBytes() []byte {
	reasonLen := len(e.Reason)
	totalLen := 1 + reasonLen
	bytes := make([]byte, 0, totalLen)

	bytes = append(bytes, byte(reasonLen))
	bytes = append(bytes, e.Reason...)

	return bytes
}

func ErrorFromBytes(data []byte) (*ErrorMessage, error) {
	// Agregar validacion largo + Cambiar mensajes de error
	offset := 0

	reasonLen := int(data[offset])
	offset++
	if len(data) < offset+reasonLen {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	reason := string(data[offset : offset+reasonLen])
	offset += reasonLen

	return &ErrorMessage{Reason: reason}, nil
}
