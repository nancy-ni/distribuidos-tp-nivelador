package messages

import (
	"fmt"

	common "github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/common"
)

type AskWinners struct {
	agencyId uint32
}

func NewAskWinners(agencyId uint32) AskWinners {
	return AskWinners{agencyId: agencyId}
}

func (a *AskWinners) ToBytes() []byte {
	bytes := make([]byte, 0, AGENCY_ID_LEN_BYTES)

	agencyIdBytes := common.Uint32ToBytes(a.agencyId)
	bytes = append(bytes, agencyIdBytes...)

	return bytes
}

func AskWinnersFromBytes(data []byte) (*AskWinners, error) {
	if len(data) < AGENCY_ID_LEN_BYTES {
		return nil, fmt.Errorf(common.AskWinnersTooShort)
	}

	agencyId, err := common.BytesToUint32(data[:AGENCY_ID_LEN_BYTES])
	if err != nil {
		return nil, fmt.Errorf(common.DeserializeAskWinnersError)
	}

	return &AskWinners{agencyId: agencyId}, nil
}
