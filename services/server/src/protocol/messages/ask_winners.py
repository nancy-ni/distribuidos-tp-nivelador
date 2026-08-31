from protocol.common.utils import bytes_to_uint32
from protocol.common import errors

class AskWinners:
    def __init__(self, agency_id):
        self.agency_id = agency_id

    @classmethod
    def from_bytes(cls, data):
        try:
            agency_id = bytes_to_uint32(data[:4])
            return cls(agency_id), None
        except Exception as e:
            return None, ValueError(f"{errors.DESERIALIZE_ASK_WINNERS_ERR}: {e}")