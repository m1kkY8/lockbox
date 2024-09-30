package encryption

import (
	"crypto/rand"
)

func GenerateAES() (AESKey, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return AESKey{}, err
	}
	return AESKey{Key: key}, nil
}
