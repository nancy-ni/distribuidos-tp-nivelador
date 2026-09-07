package messages

import (
	"fmt"
	"strconv"

	common "github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/common"
)

const BET_MIN_LEN = 22
const AGENCY_ID_LEN_BYTES = 4
const FIRSTNAME_MAX_LEN = 255
const LASTNAME_MAX_LEN = 255
const BIRTHDAY_LEN_BYTES = 10
const DNI_LEN_BYTES = 4
const BET_NUMBER_LEN_BYTES = 2

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

// Serializa un Bet. El orden de los campos serializados es: agencyId, firstnameLen, firstname,
// lastnameLen, lastname, birthday, dni, betNumber.
func (bet *Bet) ToBytes() []byte {
	firstnameLen := len(bet.firstName)
	lastnameLen := len(bet.lastName)

	totalLen := BET_MIN_LEN + firstnameLen + lastnameLen
	bytes := make([]byte, 0, totalLen)

	// agencyId
	agencyIdBytes := common.Uint32ToBytes(bet.agencyId)
	bytes = append(bytes, agencyIdBytes...)

	// firstnameLen + firstname
	bytes = append(bytes, byte(firstnameLen))
	bytes = append(bytes, bet.firstName...)

	// lastnameLen + lastname
	bytes = append(bytes, byte(lastnameLen))
	bytes = append(bytes, bet.lastName...)

	// birthday
	bytes = append(bytes, bet.birthday...)

	// dni
	dniBytes := common.Uint32ToBytes(bet.dni)
	bytes = append(bytes, dniBytes...)

	// betNumber
	betNumberBytes := common.Uint16ToBytes(bet.betNumber)
	bytes = append(bytes, betNumberBytes...)

	return bytes
}

// Deserializa un Bet, leyendo los campos en el mismo orden que la serializacion.
func BetFromBytes(data []byte) (*Bet, error) {
	if len(data) < BET_MIN_LEN {
		return nil, fmt.Errorf(common.BetTooShort)
	}
	offset := 0

	// agencyId
	if len(data) < offset+AGENCY_ID_LEN_BYTES {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	agencyId, err := common.BytesToUint32(data[offset : offset+AGENCY_ID_LEN_BYTES])
	if err != nil {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	offset += AGENCY_ID_LEN_BYTES

	// firstname
	firstNameLen := int(data[offset])
	offset++
	if len(data) < offset+firstNameLen {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	firstName := string(data[offset : offset+firstNameLen])
	offset += firstNameLen

	// lastname
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
	if len(data) < offset+BIRTHDAY_LEN_BYTES {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	birthday := string(data[offset : offset+BIRTHDAY_LEN_BYTES])
	offset += BIRTHDAY_LEN_BYTES

	// dni
	if len(data) < offset+DNI_LEN_BYTES {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	dni, err := common.BytesToUint32(data[offset : offset+DNI_LEN_BYTES])
	if err != nil {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	offset += DNI_LEN_BYTES

	// betNumber
	if len(data) < offset+BET_NUMBER_LEN_BYTES {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	betNumber, err := common.BytesToUint16(data[offset : offset+BET_NUMBER_LEN_BYTES])
	if err != nil {
		return nil, fmt.Errorf(common.DeserializeBetError)
	}
	offset += BET_NUMBER_LEN_BYTES

	return &Bet{agencyId: agencyId, firstName: firstName, lastName: lastName, birthday: birthday, dni: dni, betNumber: betNumber}, nil
}
