package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
)

func DecryptAesKey(encryptedKey []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	decryptedKey, err := rsa.DecryptOAEP(
		sha512.New(),
		rand.Reader,
		privateKey,
		encryptedKey,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return decryptedKey, nil
}
