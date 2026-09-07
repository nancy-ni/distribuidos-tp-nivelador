from protocol.common.errors import DESERIALIZE_UINT32_ERR, DESERIALIZE_UINT16_ERR

def uint32_to_bytes(number):
    return bytes([
        (number >> 24) & 0xFF,
        (number >> 16) & 0xFF,
        (number >> 8) & 0xFF,
        number & 0xFF
    ])


def bytes_to_uint32(data):
    if len(data) < 4:
        raise ValueError(DESERIALIZE_UINT32_ERR)
    return (
        data[0] << 24 | 
        data[1] << 16 |
        data[2] << 8 | 
        data[3]
    )


def uint16_to_bytes(number):
    return bytes([
        (number >> 8) & 0xFF,
        number & 0xFF
    ])


def bytes_to_uint16(data):
    if len(data) < 2:
        raise ValueError(DESERIALIZE_UINT16_ERR)
    return (
        data[0] << 8 |
        data[1]
    )