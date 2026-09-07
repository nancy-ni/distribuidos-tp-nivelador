from protocol.common.utils import bytes_to_uint32
from protocol.common import errors
from protocol.messages.bet import AGENCY_ID_LEN_BYTES

class AskWinners:
    def __init__(self, agency_id):
        self.agency_id = agency_id


    @classmethod
    def from_bytes(cls, data):
        if len(data) < AGENCY_ID_LEN_BYTES:
            raise ValueError(errors.ASK_WINNERS_TOO_SHORT_ERR)
        try:
            agency_id = bytes_to_uint32(data[:AGENCY_ID_LEN_BYTES])
            return cls(agency_id)
        except Exception as e:
            raise ValueError(f"{errors.DESERIALIZE_ASK_WINNERS_ERR}: {e}")