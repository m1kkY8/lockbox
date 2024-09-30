package teamodel

import (
	"time"

	"github.com/m1kkY8/gochat/src/encryption"
	"github.com/m1kkY8/gochat/src/message"
)

func (m *Model) createMessage(content string) message.Message {
	timestamp := time.Now().Format(time.TimeOnly)

	aesKey, _ := encryption.GenerateAES()

	var userMessage message.Message
	userMessage.Type = message.ChatMessage
	userMessage.Author = m.username
	userMessage.Timestamp = timestamp
	userMessage.Color = m.userColor
	userMessage.Room = m.currentRoom

	for _, pubKey := range m.PublicKeys {
		encryptedKey, err := encryption.EncryptAesKey(aesKey, pubKey)
		if err != nil {
			return message.Message{}
		}
		userMessage.AESKeys = append(userMessage.AESKeys, encryption.AESKey{Key: encryptedKey})
	}

	encryptedContent, _ := encryption.EncryptMessage([]byte(content), aesKey.Key)
	userMessage.Content = string(encryptedContent)

	return userMessage
}
