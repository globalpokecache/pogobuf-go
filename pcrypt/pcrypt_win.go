// +build windows

package pcrypt

/*
#cgo LDFLAGS: -l ws2_32
#include <pcrypt.h>
#include <stdlib.h>       // for free()
*/
import "C"
