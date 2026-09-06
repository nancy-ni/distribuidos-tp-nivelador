from protocol.messages import message_codes, bet, ask_winners, batch
from protocol.common import errors

PACKET_MIN_LEN = 2
PACKET_LENGTH_BYTES = 2

class Packet:
    def __init__(self, message_code, message):
        self.message_code = message_code
        self.message = message

    def to_bytes(self):
        buf = bytearray()
        buf.append(self.message_code)
        buf.extend(self.message.to_bytes())
        return bytes(buf)


    @classmethod
    def from_bytes(cls, data):
        if len(data) < PACKET_MIN_LEN:
            raise ValueError(errors.PACKET_TOO_SHORT_ERR)
        offset = 0
        message_code = int(data[offset])
        offset += 1

        try:
            match message_code:
                case message_codes.BATCH_CODE:
                    message = batch.Batch.from_bytes(data[offset:])
                case message_codes.ASK_WINNERS_CODE:
                    message = ask_winners.AskWinners.from_bytes(data[offset:])
                case _:
                    raise ValueError(errors.UNEXPECTED_MESSAGE_ERR)

        except Exception as e:
            raise ValueError(f"{errors.DESERIALIZE_PACKET_ERR}: {e}")

        return cls(message_code, message)
