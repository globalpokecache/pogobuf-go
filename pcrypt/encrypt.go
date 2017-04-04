package pcrypt

import (
	"encoding/binary"

	"golang.org/x/crypto/twofish"
)

var (
	// 1.29.1 twofish key
	encKey = []byte{
		0x4F, 0xEB, 0x1C, 0xA5, 0xF6, 0x1A, 0x67, 0xCE,
		0x43, 0xF3, 0xF0, 0x0C, 0xB1, 0x23, 0x88, 0x35,
		0xE9, 0x8B, 0xE8, 0x39, 0xD8, 0x89, 0x8F, 0x5A,
		0x3B, 0x51, 0x2E, 0xA9, 0x47, 0x38, 0xC4, 0x14,
	}
)

type cRand struct {
	state uint32
}

func (rand *cRand) rand() byte {
	rand.state = (0x41C64E6D * rand.state) + 0x3039
	return byte((rand.state >> 16) & 0x7FFF)
}

func makeIv(rand *cRand) []byte {
	iv := make([]byte, twofish.BlockSize)
	for i := 0; i < len(iv); i++ {
		iv[i] = rand.rand()
	}
	return iv
}

func makeIntegrityByte(rand *cRand) byte {
	// hardcoded since 0.59.1
	return 0x21
}

func Encrypt(input []byte, msSinceStart uint32) []byte {
	tf, _ := twofish.NewCipher(encKey)

	rand := &cRand{msSinceStart}
	iv := makeIv(rand)

	inputlen := len(input)
	blockCount := (inputlen + 256) / 256

	outputSize := (blockCount * 256) + 5
	output := make([]byte, outputSize)

	binary.BigEndian.PutUint32(output, msSinceStart)

	copy(output[4:], input)
	output[4+inputlen] = byte(256 - inputlen%256)
	output[outputSize-2] = byte(256 - inputlen%256)

	for offset := 4; offset < blockCount*256; offset += twofish.BlockSize {
		for i := 0; i < twofish.BlockSize; i++ {
			output[offset+i] ^= iv[i]
		}
		tf.Encrypt(output[offset:], output[offset:])
		copy(iv, output[offset:offset+twofish.BlockSize])
	}

	output[outputSize-1] = makeIntegrityByte(rand)

	return output
}
