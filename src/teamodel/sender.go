package teamodel

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/m1kkY8/lockbox/src/message"
)

func (m *Model) sendMessage(userMessage message.Message) error {
	byteMessage, err := message.EncodeMessage(userMessage)
	if err != nil {
		log.Println("Failed encoding message")
		return err
	}

	err = m.conn.WriteMessage(websocket.BinaryMessage, byteMessage)
	if err != nil {
		log.Printf("Error writing message: %v", err)
		return err
	}

	return nil
}
