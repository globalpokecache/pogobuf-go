// +build windows

package pcrypt

/*
#cgo LDFLAGS: -l ws2_32
#include <pcrypt.h>
#include <stdlib.h>       // for free()
*/
import "C"
import "unsafe"

func Encrypt(input []byte, msSinceStart uint32) []byte {
	buf := C.CString("")
	d := C.encrypt((*C.char)(unsafe.Pointer(&input[0])), C.size_t(len(input)), C.uint32_t(msSinceStart), &buf)
	defer C.free(unsafe.Pointer(buf))
	return C.GoBytes(unsafe.Pointer(buf), d)
}
