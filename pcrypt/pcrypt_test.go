package pcrypt

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	enc := Encrypt([]byte{0, 1, 2, 3, 6, 7, 8, 9}, 100)
	fmt.Println(enc)
	dec := Decrypt(enc)
	fmt.Println(dec)
}
