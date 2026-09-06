from protocol.common import errors

class ErrorMessage:
    def __init__(self, reason):
        self.reason = reason

    def to_bytes(self):
        buf = bytearray()
        reason_bytes = bytes(self.reason, "utf-8")
        buf.append(len(reason_bytes))
        buf.extend(reason_bytes)
        return bytes(buf)

    @classmethod
    def from_bytes(cls, data):
        offset = 0
        try:
            reason_len = data[offset]
            offset += 1
            reason = data[offset : offset + reason_len].decode("utf-8")
            offset += reason_len

            return cls(reason)
        except Exception as e:
            raise ValueError(f"{errors.DESERIALIZE_FINISH_ERR}: {e}")