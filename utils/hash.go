package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"hash"
)

func GenerateRandomSalt(saltSize int) []byte {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt[:])

	if err != nil {
		panic(err)
	}

	return salt
}

func ComputeHMAC(hg func() hash.Hash, key, data []byte) []byte {
	mac := hmac.New(hg, key)
	mac.Write(data)
	return mac.Sum(nil)
}

func ComputeHash(hg func() hash.Hash, b []byte) []byte {
	h := hg()
	h.Write(b)
	return h.Sum(nil)
}
