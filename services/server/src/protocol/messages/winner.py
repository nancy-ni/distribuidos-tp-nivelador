from lottery.bet import Bet
from protocol.messages.bet import BetWrapper
from protocol.common import errors

class Winner:
    def __init__(self, bet):
        self.winner = bet

    def to_bytes(self):
        buf = bytearray()

        buf.extend(self.winner.bet.agency_id.to_bytes(4, byteorder="big"))

        firstname_bytes = bytes(self.winner.bet.first_name, "utf-8")
        buf.append(len(firstname_bytes))
        buf.extend(firstname_bytes)

        lastname_bytes = bytes(self.winner.bet.last_name, "utf-8")
        buf.append(len(lastname_bytes))
        buf.extend(lastname_bytes)

        buf.extend(bytes(self.winner.bet.birthdate, "utf-8"))

        buf.extend(self.winner.bet.document.to_bytes(4, byteorder="big"))

        buf.extend(self.winner.bet.number.to_bytes(2, byteorder="big"))

        return bytes(buf)

    @classmethod
    def from_bytes(cls, data):
        offset = 0

        try:
            agency_id = int.from_bytes(data[offset : offset + 4], byteorder='big')
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

            dni = int.from_bytes(data[offset : offset + 4], byteorder='big')
            offset += 4

            betNumber = int.from_bytes(data[offset : offset + 2], byteorder='big')
            offset += 2

            bet = Bet(agency_id, firstname, lastname, dni, birthday, betNumber)
            bet_wrapper = BetWrapper(bet)

            return cls(bet_wrapper), None

        except Exception as e:
            return None, ValueError(f"{errors.DESERIALIZE_BET_ERR}: {e}")
