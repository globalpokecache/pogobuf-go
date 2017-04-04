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
	fmt.Println(input)

	result := Encrypt(input, uint32(time.Now().Unix()))
	fmt.Println(result)
}
