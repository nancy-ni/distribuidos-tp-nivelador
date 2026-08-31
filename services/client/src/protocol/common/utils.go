package common

import "fmt"

func Uint32ToBytes(num uint32) []byte {
	bytes := make([]byte, 4)
	bytes[0] = byte(num >> 24)
	bytes[1] = byte(num >> 16)
	bytes[2] = byte(num >> 8)
	bytes[3] = byte(num)
	return bytes
}

func Uint16ToBytes(num uint16) []byte {
	bytes := make([]byte, 2)
	bytes[0] = byte(num >> 8)
	bytes[1] = byte(num)
	return bytes
}

func BytesToUint32(bytes []byte) (uint32, error) {
	if len(bytes) < 4 {
		return 0, fmt.Errorf(DeserializeUint32Error)
	}
	return uint32(bytes[0])<<24 |
		uint32(bytes[1])<<16 |
		uint32(bytes[2])<<8 |
		uint32(bytes[3]), nil
}

func BytesToUint16(bytes []byte) (uint16, error) {
	if len(bytes) < 2 {
		return 0, fmt.Errorf(DeserializeUint16Error)
	}
	return uint16(bytes[0])<<8 |
		uint16(bytes[1]), nil
}
