package messages

import (
	"fmt"

	common "github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/common"
)

const BATCH_MIN_LEN = 2

type Batch struct {
	Bets []Bet
}

func NewBatch(bets []Bet) Batch {
	return Batch{Bets: bets}
}

func (batch *Batch) ToBytes() []byte {
	var bytes []byte

	batchSize := uint16(len(batch.Bets))
	batchSizeBytes := common.Uint16ToBytes(batchSize)
	bytes = append(bytes, batchSizeBytes...)

	for i := range batchSize {
		bet := batch.Bets[i]
		betBytes := bet.ToBytes()
		betLengthBytes := common.Uint16ToBytes(uint16(len(betBytes)))

		bytes = append(bytes, betLengthBytes...)
		bytes = append(bytes, betBytes...)
	}

	return bytes
}

func BatchFromBytes(data []byte) (*Batch, error) {
	if len(data) < BATCH_MIN_LEN {
		return nil, fmt.Errorf(common.BatchTooShort)
	}
	offset := 0

	batchSize, err := common.BytesToUint16(data[offset : offset+2])
	if err != nil {
		return nil, fmt.Errorf(common.DeserializeBatchError)
	}
	offset += 2

	batch := Batch{Bets: make([]Bet, 0, batchSize)}

	for _ = range batchSize {
		if len(data) < offset+2 {
			return nil, fmt.Errorf(common.DeserializeBatchError)
		}
		betLength, err := common.BytesToUint16(data[offset : offset+2])
		if err != nil {
			return nil, fmt.Errorf(common.DeserializeBatchError)
		}
		offset += 2

		if len(data) < offset+int(betLength) {
			return nil, fmt.Errorf(common.DeserializeBatchError)
		}
		bet, err := BetFromBytes(data[offset : offset+int(betLength)])
		if err != nil {
			return nil, fmt.Errorf(common.DeserializeBatchError)
		}
		offset += int(betLength)

		batch.Bets = append(batch.Bets, *bet)
	}

	return &batch, nil
}
