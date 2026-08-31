from protocol.messages import message_codes, bet, ask_winners
from protocol.common import errors

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
        offset = 0

        if len(data) < offset + 1:
            return None, ValueError(errors.PACKET_TOO_SHORT_ERR)
        message_code = int(data[offset])
        offset += 1

        try:
            match message_code:
                case message_codes.BET_CODE:
                    message, err = bet.BetWrapper.from_bytes(data[offset:])
                case message_codes.ASK_WINNERS_CODE:
                    message, err = ask_winners.AskWinners.from_bytes(data[offset:])
                case _:
                    return None, ValueError(errors.UNEXPECTED_MESSAGE_ERR)
            if err is not None:
                return None, err
        except Exception as e:
            return None, ValueError(f"{errors.DESERIALIZE_PACKET_ERR}: {e}")

        return cls(message_code, message), None
