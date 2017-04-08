package pcrypt

import (
	"encoding/binary"

	"golang.org/x/crypto/twofish"
)

func Decrypt(buffer []byte) []byte {
	tf, _ := twofish.NewCipher(encKey)

	rand := &cRand{
		binary.BigEndian.Uint32(buffer[0:4]),
	}
	iv := makeIv(rand)

	size := len(buffer)

	blockCount := (size - 5) / twofish.BlockSize

	for offset := 4; offset < blockCount*twofish.BlockSize; offset += twofish.BlockSize {
		tarbyte := make([]byte, twofish.BlockSize)

		tf.Decrypt(tarbyte, buffer[offset:])

		for i := 0; i < twofish.BlockSize; i++ {
			tarbyte[i] ^= iv[i]
			iv[i] = buffer[offset+i]
			buffer[offset+i] = tarbyte[i]
		}
	}

	zeroSize := int(buffer[size-1-1])
	tarSize := size - 1 - zeroSize

	return buffer[4:tarSize]
}
