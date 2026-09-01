from protocol.common import errors
from protocol.common import utils
from protocol.messages.bet import BetWrapper

BATCH_MIN_LEN = 2

class Batch:
    def __init__(self, bets):
        self.bets = bets

    def get_bets(self):
        return [wrapper.bet for wrapper in self.bets]

    def to_bytes(self):
        buf = bytearray()
        batch_size = len(self.bets)
        buf.extend(batch_size.to_bytes(2, byteorder="big"))

        for i in range(batch_size):
            bet = self.bets[i]
            bet_bytes = bet.bet.to_bytes()
            bet_length = len(bet_bytes)

            buf.extend(bet_length.to_bytes(2, byteorder="big"))
            buf.extend(bet_bytes)

        return bytes(buf)

    @classmethod
    def from_bytes(cls, data):
        if len(data) < BATCH_MIN_LEN:
            return None, ValueError(errors.BATCH_TOO_SHORT_ERR)
        offset = 0

        batch_size = utils.bytes_to_uint16(data[offset:offset+2])
        offset += 2

        batch = cls([])

        for _ in range(batch_size):
            if len(data) < offset + 2:
                return None, ValueError(errors.DESERIALIZE_BATCH_ERR)
            bet_length = utils.bytes_to_uint16(data[offset:offset+2])
            offset += 2

            if len(data) < offset + bet_length:
                return None, ValueError(errors.DESERIALIZE_BATCH_ERR)
            bet, err = BetWrapper.from_bytes(data[offset:offset+bet_length])
            if err is not None:
                return None, ValueError(errors.DESERIALIZE_BATCH_ERR)
            offset += bet_length

            batch.bets.append(bet)

        return batch, None




