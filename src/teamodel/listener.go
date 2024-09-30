package teamodel

import (
	"log"
	"strings"

	// "github.com/m1kkY8/lockbox/src/encryption"
	"github.com/m1kkY8/lockbox/src/encryption"
	"github.com/m1kkY8/lockbox/src/message"
	"github.com/m1kkY8/lockbox/src/notification"
	"github.com/vmihailenco/msgpack/v5"
)

func (m *Model) recieveMessages() {
	for {

		// Read messages from socket
		_, byteMessage, err := m.conn.ReadMessage()
		if err != nil {
			return
		}

		// Decode as regular message
		decodedMsg, err := message.DecodeMessage(byteMessage)
		if err != nil {
			log.Println("Failed decoding")
			continue
		}

		switch decodedMsg.Type {

		case message.ServerMessage:
			// Handle online users Message
			m.onlineUsersChan <- strings.Split(decodedMsg.Content, " ")

		case message.KeyMessage:
			// Handle public keys
			var keyMsg message.PublicKeys

			err := msgpack.Unmarshal(byteMessage, &keyMsg)
			if err != nil {
				continue
			}
			m.PublicKeys = keyMsg.PublicKeys

		case message.ChatMessage:
			// Iterate over keys to find right one to unlock AES key
			for _, key := range decodedMsg.AESKeys {
				aesKey, err := encryption.DecryptAesKey(key.Key, m.keyPair.PrivateKey)
				if err != nil {
					continue
				}

				// ulock message with aes key
				decryptedContent, err := encryption.DecryptMessage([]byte(decodedMsg.Content), aesKey)
				if err != nil {
					continue
				}

				decodedMsg.Content = string(decryptedContent)
			}

			formattedMessage := message.Format(decodedMsg)
			m.messageChannel <- formattedMessage
			notification.Notify(decodedMsg, m.username)

		default:
			continue
		}
	}
}
