package teamodel

import (
	"log"
	"strings"

	"github.com/m1kkY8/gochat/src/message"
	"github.com/m1kkY8/gochat/src/notification"
)

func (m *Model) recieveMessages() {
	for {
		_, byteMessage, err := m.conn.ReadMessage()
		if err != nil {
			return
		}

		decodedMsg, err := message.DecodeMessage(byteMessage)
		if err != nil {
			log.Println("Failed decoding")
			continue
		}

		if decodedMsg.Type == message.ServerMessage {
			m.onlineUsersChan <- strings.Split(decodedMsg.Content, " ")
		} else {
			if decodedMsg.Author == decodedMsg.To {
				continue
			}

			if decodedMsg.To != "" && decodedMsg.To == m.username {
				// This is a whisper message intended for this client
				formattedMessage := message.FormatWhisper(decodedMsg)
				m.messageChannel <- formattedMessage
				continue
			}

			// Handle Regular Message (this includes whispers sent to other users)
			if decodedMsg.To == "" || decodedMsg.To == "all" {
				formattedMessage := message.Format(decodedMsg)
				m.messageChannel <- formattedMessage
				notification.Notify(decodedMsg, m.username)
			}
		}
	}
}
