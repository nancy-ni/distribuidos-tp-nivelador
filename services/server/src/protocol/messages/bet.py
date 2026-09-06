from lottery.bet import Bet
from protocol.common.utils import bytes_to_uint32, bytes_to_uint16
from protocol.common import errors

BET_MIN_LEN = 22
AGENCY_ID_LEN_BYTES = 4
BIRTHDAY_LEN_BYTES = 10
DNI_LEN_BYTES = 4
BET_NUMBER_LEN_BYTES = 2

class BetWrapper:
    def __init__(self, bet):
        self.bet = bet

    def to_bytes(self):
        buf = bytearray()

        buf.extend(self.bet.agency_id.to_bytes(AGENCY_ID_LEN_BYTES, byteorder="big"))

        firstname_bytes = bytes(self.bet.first_name, "utf-8")
        buf.append(len(firstname_bytes))
        buf.extend(firstname_bytes)

        lastname_bytes = bytes(self.bet.last_name, "utf-8")
        buf.append(len(lastname_bytes))
        buf.extend(lastname_bytes)

        buf.extend(bytes(self.bet.birthdate, "utf-8"))

        buf.extend(self.bet.document.to_bytes(DNI_LEN_BYTES, byteorder="big"))

        buf.extend(self.bet.number.to_bytes(BET_NUMBER_LEN_BYTES, byteorder="big"))

        return bytes(buf)

    @classmethod
    def from_bytes(cls, data):
        if len(data) < BET_MIN_LEN:
            raise ValueError(errors.BET_TOO_SHORT_ERR)
        offset = 0

        try:
            agency_id = bytes_to_uint32(data[offset : offset + AGENCY_ID_LEN_BYTES])
            offset += AGENCY_ID_LEN_BYTES

            firstname_len = data[offset]
            offset += 1
            firstname = data[offset : offset + firstname_len].decode("utf-8")
            offset += firstname_len

            lastname_len = data[offset]
            offset += 1
            lastname = data[offset : offset + lastname_len].decode("utf-8")
            offset += lastname_len

            birthday = data[offset : offset + BIRTHDAY_LEN_BYTES].decode("utf-8")
            offset += BIRTHDAY_LEN_BYTES

            dni = bytes_to_uint32(data[offset : offset + DNI_LEN_BYTES])
            offset += DNI_LEN_BYTES

            betNumber = bytes_to_uint16(data[offset : offset + BET_NUMBER_LEN_BYTES])
            offset += BET_NUMBER_LEN_BYTES

            bet = Bet(agency_id, firstname, lastname, dni, birthday, betNumber)

            return cls(bet)
        
        except Exception as e:
            raise ValueError(f"{errors.DESERIALIZE_BET_ERR}: {e}")

