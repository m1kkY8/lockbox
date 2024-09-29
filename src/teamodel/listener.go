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

		switch decodedMsg.Type {
		case message.ServerMessage:
			m.onlineUsersChan <- strings.Split(decodedMsg.Content, " ")

		case message.ChatMessage:
			formattedMessage := message.Format(decodedMsg)
			m.messageChannel <- formattedMessage

			notification.Notify(decodedMsg, m.username)
		case message.CommandMessage:
			continue

		default:
			continue
		}
	}
}
