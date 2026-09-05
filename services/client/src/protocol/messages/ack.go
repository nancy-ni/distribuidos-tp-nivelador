package messages

import (
	"fmt"

	common "github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/common"
)

type Ack struct {
	agencyId uint32
}

func (a *Ack) ToBytes() []byte {
	bytes := make([]byte, 0, 4)

	agencyIdBytes := common.Uint32ToBytes(a.agencyId)
	bytes = append(bytes, agencyIdBytes...)

	return bytes
}

func AckFromBytes(data []byte) (*Ack, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf(common.FinishTooShort)
	}

	agencyId, err := common.BytesToUint32(data[:4])
	if err != nil {
		return nil, fmt.Errorf(common.DeserializeFinishError)
	}

	return &Ack{agencyId: agencyId}, nil
}
