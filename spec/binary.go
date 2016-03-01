package spec

import "encoding/binary"

var (
	DefaultEndianness = binary.BigEndian
)

func Uint16(b []byte) uint16 {
	return DefaultEndianness.Uint16(pad(b, 2))
}

func Uint32(b []byte) uint32 {
	return DefaultEndianness.Uint32(pad(b, 4))
}

func Uint64(b []byte) uint64 {
	return DefaultEndianness.Uint64(pad(b, 8))
}

func pad(b []byte, n int) []byte {
	return append(make([]byte, n-len(b)), b...)
}
