package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
)

// Return encrypted AES key
func EncryptAesKey(key AESKey, pubKey *rsa.PublicKey) ([]byte, error) {
	encryptedKey, err := rsa.EncryptOAEP(
		sha512.New(),
		rand.Reader,
		pubKey,
		key.Key,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return encryptedKey, nil
}
