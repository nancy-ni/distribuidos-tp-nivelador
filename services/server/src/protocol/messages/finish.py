from protocol.common.utils import bytes_to_uint32
from protocol.common import errors

class Finish:
    def __init__(self, agency_id):
        self.agency_id = agency_id

    def to_bytes(self):
        buf = bytearray()
        buf.extend(self.agency_id.to_bytes(4, byteorder="big"))
        return bytes(buf)

    @classmethod
    def from_bytes(cls, data):
        try:
            agency_id = bytes_to_uint32(data[:4])
            return cls(agency_id), None
        except Exception as e:
            return None, ValueError(f"{errors.DESERIALIZE_FINISH_ERR}: {e}")