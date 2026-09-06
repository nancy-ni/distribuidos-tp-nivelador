package messages

import (
	"fmt"

	common "github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/common"
)

type Finish struct {
	agencyId uint32
}

func (f *Finish) ToBytes() []byte {
	bytes := make([]byte, 0, AGENCY_ID_LEN_BYTES)

	agencyIdBytes := common.Uint32ToBytes(f.agencyId)
	bytes = append(bytes, agencyIdBytes...)

	return bytes
}

func FinishFromBytes(data []byte) (*Finish, error) {
	if len(data) < AGENCY_ID_LEN_BYTES {
		return nil, fmt.Errorf(common.FinishTooShort)
	}

	agencyId, err := common.BytesToUint32(data[:AGENCY_ID_LEN_BYTES])
	if err != nil {
		return nil, fmt.Errorf(common.DeserializeFinishError)
	}

	return &Finish{agencyId: agencyId}, nil
}
