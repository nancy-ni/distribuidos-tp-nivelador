package messages

import (
	"fmt"

	common "github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/common"
)

type Finish struct {
	agencyId uint32
}

func (f *Finish) ToBytes() []byte {
	bytes := make([]byte, 0, 4)

	agencyIdBytes := common.Uint32ToBytes(f.agencyId)
	bytes = append(bytes, agencyIdBytes...)

	return bytes
}

func FinishFromBytes(data []byte) (*Finish, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf(common.FinishTooShort)
	}

	agencyId, err := common.BytesToUint32(data[:4])
	if err != nil {
		return nil, fmt.Errorf(common.DeserializeFinishError)
	}

	return &Finish{agencyId: agencyId}, nil
}
