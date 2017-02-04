package pcrypt

import (
	"crypto/rand"
	"reflect"
	"unsafe"
)

const SIZEOF_INT32 = 4 // bytes
func RandomBytes(n int) []byte {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return bytes
}
func AsDwordSlice(bytes []byte) []uint32 {
	// Get the slice header
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	// The length and capacity of the slice are different.
	header.Len /= SIZEOF_INT32
	header.Cap /= SIZEOF_INT32
	// Convert slice header to an []uint32
	return *(*[]uint32)(unsafe.Pointer(&header))
}
