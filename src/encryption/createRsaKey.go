package encryption

import (
	"crypto/rand"
	"crypto/rsa"
)

func CreateKey() (*RSAKeys, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	publicKey := &privateKey.PublicKey

	return &RSAKeys{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}
