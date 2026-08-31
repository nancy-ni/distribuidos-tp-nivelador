from lottery.bet import Bet
from protocol.common.utils import bytes_to_uint32, bytes_to_uint16
from protocol.common import errors

class BetWrapper:
    def __init__(self, bet):
        self.bet = bet

    def to_bytes(self):
        buf = bytearray()

        buf.extend(self.bet.agency_id.to_bytes(4, byteorder="big"))

        firstname_bytes = bytes(self.bet.first_name, "utf-8")
        buf.append(len(firstname_bytes))
        buf.extend(firstname_bytes)

        lastname_bytes = bytes(self.bet.last_name, "utf-8")
        buf.append(len(lastname_bytes))
        buf.extend(lastname_bytes)

        buf.extend(bytes(self.bet.birthdate, "utf-8"))

        buf.extend(self.bet.document.to_bytes(4, byteorder="big"))

        buf.extend(self.bet.number.to_bytes(2, byteorder="big"))

        return bytes(buf)

    @classmethod
    def from_bytes(cls, data):
        offset = 0

        try:
            agency_id = bytes_to_uint32(data[offset : offset + 4])
            offset += 4

            firstname_len = data[offset]
            offset += 1
            firstname = data[offset : offset + firstname_len].decode("utf-8")
            offset += firstname_len

            lastname_len = data[offset]
            offset += 1
            lastname = data[offset : offset + lastname_len].decode("utf-8")
            offset += lastname_len

            birthday = data[offset : offset + 10].decode("utf-8")
            offset += 10

            dni = bytes_to_uint32(data[offset : offset + 4])
            offset += 4

            betNumber = bytes_to_uint16(data[offset : offset + 2])
            offset += 2

            bet = Bet(agency_id, firstname, lastname, dni, birthday, betNumber)

            return cls(bet), None
        
        except Exception as e:
            return None, ValueError(f"{errors.DESERIALIZE_BET_ERR}: {e}")

