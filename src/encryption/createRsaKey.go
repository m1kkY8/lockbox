package encryption

import (
	"crypto/rand"
	"crypto/rsa"
)

func CreateRsaKey() (*RSAKeys, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	publicKey := &privateKey.PublicKey

	return &RSAKeys{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}
