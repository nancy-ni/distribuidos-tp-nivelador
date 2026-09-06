from lottery.bet import Bet
from protocol.messages.bet import BetWrapper
from protocol.common import errors
from protocol.messages.bet import BET_MIN_LEN, AGENCY_ID_LEN_BYTES, BIRTHDAY_LEN_BYTES, DNI_LEN_BYTES, BET_NUMBER_LEN_BYTES

class Winner:
    def __init__(self, bet):
        self.winner = bet

    def to_bytes(self):
        buf = bytearray()

        buf.extend(self.winner.bet.agency_id.to_bytes(AGENCY_ID_LEN_BYTES, byteorder="big"))

        firstname_bytes = bytes(self.winner.bet.first_name, "utf-8")
        buf.append(len(firstname_bytes))
        buf.extend(firstname_bytes)

        lastname_bytes = bytes(self.winner.bet.last_name, "utf-8")
        buf.append(len(lastname_bytes))
        buf.extend(lastname_bytes)

        buf.extend(bytes(self.winner.bet.birthdate, "utf-8"))

        buf.extend(self.winner.bet.document.to_bytes(DNI_LEN_BYTES, byteorder="big"))

        buf.extend(self.winner.bet.number.to_bytes(BET_NUMBER_LEN_BYTES, byteorder="big"))

        return bytes(buf)

    @classmethod
    def from_bytes(cls, data):
        if len(data) < BET_MIN_LEN:
            raise ValueError(errors.BET_TOO_SHORT_ERR)
        offset = 0

        try:
            agency_id = int.from_bytes(data[offset : offset + AGENCY_ID_LEN_BYTES], byteorder='big')
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

            dni = int.from_bytes(data[offset : offset + DNI_LEN_BYTES], byteorder='big')
            offset += DNI_LEN_BYTES

            betNumber = int.from_bytes(data[offset : offset + BET_NUMBER_LEN_BYTES], byteorder='big')
            offset += BET_NUMBER_LEN_BYTES

            bet = Bet(agency_id, firstname, lastname, dni, birthday, betNumber)
            bet_wrapper = BetWrapper(bet)

            return cls(bet_wrapper)

        except Exception as e:
            raise ValueError(f"{errors.DESERIALIZE_BET_ERR}: {e}")
