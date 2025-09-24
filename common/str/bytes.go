package str

import (
	"unsafe"
)

// FromBytes converts byte slice to string.
func FromBytes(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ToBytes converts string to byte slice.
func ToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
