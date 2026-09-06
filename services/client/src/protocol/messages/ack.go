package messages

import (
	"fmt"

	common "github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/common"
)

type Ack struct {
	agencyId uint32
}

func (a *Ack) ToBytes() []byte {
	bytes := make([]byte, 0, AGENCY_ID_LEN_BYTES)

	agencyIdBytes := common.Uint32ToBytes(a.agencyId)
	bytes = append(bytes, agencyIdBytes...)

	return bytes
}

func AckFromBytes(data []byte) (*Ack, error) {
	if len(data) < AGENCY_ID_LEN_BYTES {
		return nil, fmt.Errorf(common.AckTooShort)
	}

	agencyId, err := common.BytesToUint32(data[:AGENCY_ID_LEN_BYTES])
	if err != nil {
		return nil, fmt.Errorf(common.DeserializeAckError)
	}

	return &Ack{agencyId: agencyId}, nil
}
