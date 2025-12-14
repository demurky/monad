package symmetric

import (
	"crypto/sha3"
	"crypto/subtle"
)

func getChecksum(h *sha3.SHAKE) []byte {
	const shakeReadLen = 32
	sum := make([]byte, shakeReadLen)
	h.Read(sum)
	return sum
}

func equal(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}
