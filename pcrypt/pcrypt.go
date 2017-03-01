package pcrypt

/*
#cgo CFLAGS: -std=c99 -shared -Wno-implicit-function-declaration -fPIC -O3
#include <pcrypt.h>
#include <stdlib.h>       // for free()
*/
import "C"
import "unsafe"

func Encrypt(input []byte, msSinceStart uint32) []byte {
	buf := C.CString("")
	d := C.encrypt((*C.char)(unsafe.Pointer(&input[0])), C.size_t(len(input)), C.uint32_t(msSinceStart), &buf, 3)
	defer C.free(unsafe.Pointer(buf))
	return C.GoBytes(unsafe.Pointer(buf), d)
}

func Decrypt(input []byte) []byte {
	buf := C.CString("")
	d := C.decrypt((*C.char)(unsafe.Pointer(&input[0])), C.size_t(len(input)), &buf)
	defer C.free(unsafe.Pointer(buf))
	var r []byte
	if d > 0 {
		r = C.GoBytes(unsafe.Pointer(buf), d)
	}
	return r
}
