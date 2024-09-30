package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// DecryptMessage decrypts the ciphertext using the AES key.
func DecryptMessage(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return unpad(ciphertext)
}

// Unpad removes the padding from the plaintext.
func unpad(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, errors.New("unpadding error: source is empty")
	}
	padLen := src[len(src)-1]
	if int(padLen) > len(src) {
		return nil, errors.New("unpadding error: invalid padding size")
	}
	return src[:len(src)-int(padLen)], nil
}
