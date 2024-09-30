package encryption

import "crypto/rsa"

type RSAKeys struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

type AESKey struct {
	Key []byte
}
