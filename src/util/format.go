package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func Format(content string, username string, style lipgloss.Style) string {
	timestamp := time.Now().Format(time.TimeOnly)
	usr := style.Render(timestamp + " " + username)
	formatted := fmt.Sprintf("%s: %s", usr, content)

	// Encode encrypted message in base64 to make it a readable string
	encryptedBase64 := base64.StdEncoding.EncodeToString([]byte(formatted))

	return encryptedBase64
}

// BUG: MRACNE SILE UNIVERZUMA NE DIRAJ
func encryptMessage(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	// Encrypt the message with the recipient's public key and random padding
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, message, nil)
	if err != nil {
		return nil, fmt.Errorf("error encrypting message: %v", err)
	}
	return ciphertext, nil
}

func loadPubKey(filepath string) (*rsa.PublicKey, error) {
	// Read public key from file
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Panicln("Error reading public key file:", err)
		return nil, err
	}

	// Decode PEM block
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing the public key")
	}

	// Parse the public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %v", err)
	}

	// Assert type to RSA public key
	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return publicKey, nil
}
