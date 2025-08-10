package teamodel

import (
	"time"

	"github.com/m1kkY8/lockbox/src/encryption"
	"github.com/m1kkY8/lockbox/src/message"
)

// createMessage creates an encrypted message ready to be sent
func (m *Model) createMessage(content string) message.Message {
	timestamp := time.Now().Format(time.TimeOnly)

	aesKey, _ := encryption.GenerateAES()

	var userMessage message.Message
	userMessage.Type = message.ChatMessage
	userMessage.Author = m.state.Username
	userMessage.Timestamp = timestamp
	userMessage.Color = m.state.UserColor
	userMessage.Room = m.state.CurrentRoom

	// Encrypt the AES key for each public key
	for _, pubKey := range m.client.PublicKeys {
		encryptedKey, err := encryption.EncryptAesKey(aesKey, pubKey)
		if err != nil {
			// Log error but continue with other keys
			continue
		}
		userMessage.AESKeys = append(userMessage.AESKeys, encryption.AESKey{Key: encryptedKey})
	}

	// Encrypt the message content
	encryptedContent, _ := encryption.EncryptMessage([]byte(content), aesKey.Key)
	userMessage.Content = string(encryptedContent)

	return userMessage
}
