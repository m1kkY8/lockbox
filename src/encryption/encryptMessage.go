package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// EncryptMessage encrypts the plaintext using the AES key.
func EncryptMessage(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext = pad(plaintext) // Pad the plaintext
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// Generate a random IV (initialization vector)
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// Create a new CBC cipher
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

// Pad adds padding to the plaintext to ensure it's a multiple of the block size.
func pad(src []byte) []byte {
	padLen := aes.BlockSize - len(src)%aes.BlockSize
	pad := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, pad...)
}
