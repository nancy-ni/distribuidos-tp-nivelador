
def uint32_to_bytes(number):
    return bytes([
        (number >> 24) & 0xFF,
        (number >> 16) & 0xFF,
        (number >> 8) & 0xFF,
        number & 0xFF
    ])

def bytes_to_uint32(data):
    offset = 0
    return (
        data[offset] << 24 | 
        data[offset+1] << 16 |
        data[offset+2] << 8 | 
        data[offset+3]
    )

def uint16_to_bytes(number):
    return bytes([
        (number >> 8) & 0xFF,
        number & 0xFF
    ])

def bytes_to_uint16(data):
    offset = 0
    return (
        data[offset] << 8 |
        data[offset+1]
    )