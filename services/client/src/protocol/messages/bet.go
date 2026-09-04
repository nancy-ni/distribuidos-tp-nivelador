package messages

import (
	"fmt"
	"strconv"

	common "github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/common"
)

const FIRSTNAME_MAX_LEN = 255
const LASTNAME_MAX_LEN = 255
const BET_MIN_LEN = 22

type Bet struct {
	agencyId  uint32
	firstName string
	lastName  string
	dni       uint32
	birthday  string
	betNumber uint16
}

func NewBet(agencyId uint32, firstName string, lastName string, dni uint32, birthday string, betNumber uint16) Bet {
	return Bet{
		agencyId:  agencyId,
		firstName: firstName,
		lastName:  lastName,
		dni:       dni,
		birthday:  birthday,
		betNumber: betNumber,
	}
}

func (bet *Bet) ToString() string {
	return bet.firstName + "," + bet.lastName + "," + strconv.Itoa(int(bet.dni)) + "," + bet.birthday + "," + strconv.Itoa(int(bet.betNumber))
}

func (bet *Bet) ToBytes() []byte {
	firstnameLen := len(bet.firstName)
	lastnameLen := len(bet.lastName)

	totalLen := BET_MIN_LEN + firstnameLen + lastnameLen
	bytes := make([]byte, 0, totalLen)

	agencyIdBytes := common.Uint32ToBytes(bet.agencyId)
	bytes = append(bytes, agencyIdBytes...)

	bytes = append(bytes, byte(firstnameLen))
	bytes = append(bytes, bet.firstName...)

	bytes = append(bytes, byte(lastnameLen))
	bytes = append(bytes, bet.lastName...)

	bytes = append(bytes, bet.birthday...)

	dniBytes := common.Uint32ToBytes(bet.dni)
	bytes = append(bytes, dniBytes...)

	betNumberBytes := common.Uint16ToBytes(bet.betNumber)
	bytes = append(bytes, betNumberBytes...)

	return bytes
}

func BetFromBytes(data []byte) (*Bet, error) {
	if len(data) < BET_MIN_LEN {
		return nil, fmt.Errorf(common.BetTooShort)
	}
	offset := 0

	// agency id
	if len(data) < offset+4 {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	agencyId, err := common.BytesToUint32(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	offset += 4

	// first name
	firstNameLen := int(data[offset])
	offset++
	if len(data) < offset+firstNameLen {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	firstName := string(data[offset : offset+firstNameLen])
	offset += firstNameLen

	// last name
	if len(data) < offset+1 {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	lastNameLen := int(data[offset])
	offset++
	if len(data) < offset+lastNameLen {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	lastName := string(data[offset : offset+lastNameLen])
	offset += lastNameLen

	// birthday
	if len(data) < offset+10 {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	birthday := string(data[offset : offset+10])
	offset += 10

	// dni
	if len(data) < offset+4 {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	dni, err := common.BytesToUint32(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	offset += 4

	// betNumber
	if len(data) < offset+2 {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	betNumber, err := common.BytesToUint16(data[offset : offset+2])
	if err != nil {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	offset += 2

	return &Bet{agencyId: agencyId, firstName: firstName, lastName: lastName, birthday: birthday, dni: dni, betNumber: betNumber}, nil
}
