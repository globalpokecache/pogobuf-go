package pcrypt

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"
)

func TestEncrypt(t *testing.T) {
	input := make([]byte, 32)
	rand.Read(input)

	ms := uint32(time.Now().Unix())
	// fmt.Println(input, ms)

	result := Encrypt(input, ms)
	fmt.Println(result)

	// dec := Decrypt(result)
	// fmt.Println(dec)

	// if bytes.Compare(input, dec) != 0 {
	// 	t.Fatal("Encrypted input is different from decrypted output")
	// }
}
