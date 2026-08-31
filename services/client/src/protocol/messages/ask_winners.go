package messages

import (
	"fmt"

	common "github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/common"
)

type AskWinners struct {
	agencyId uint32
}

func NewInquirie(agencyId uint32) AskWinners {
	return AskWinners{agencyId: agencyId}
}

func (inquirie *AskWinners) ToBytes() []byte {
	bytes := make([]byte, 0, 4)

	agencyIdBytes := common.Uint32ToBytes(inquirie.agencyId)
	bytes = append(bytes, agencyIdBytes...)

	return bytes
}

func AskWinnersFromBytes(data []byte) (*AskWinners, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf(common.AskWinnersTooShort)
	}

	agencyId, err := common.BytesToUint32(data[:4])
	if err != nil {
		return nil, fmt.Errorf(common.DeserializeAskWinnersError)
	}

	return &AskWinners{agencyId: agencyId}, nil
}
