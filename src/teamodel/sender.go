package teamodel

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/m1kkY8/lockbox/src/message"
)

// sendMessage sends a message through the WebSocket connection
func (m *Model) sendMessage(userMessage message.Message) error {
	byteMessage, err := message.EncodeMessage(userMessage)
	if err != nil {
		log.Printf("Failed encoding message: %v", err)
		return err
	}

	err = m.client.Conn.WriteMessage(websocket.BinaryMessage, byteMessage)
	if err != nil {
		log.Printf("Error writing message: %v", err)
		return err
	}

	return nil
}
