package cryptography

import "crypto/sha256"

func Hash(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

func HashString(data string) []byte {
	bData := []byte(data)
	h := sha256.New()
	h.Write(bData)
	return h.Sum(nil)
}
