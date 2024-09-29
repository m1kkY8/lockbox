package teamodel

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/m1kkY8/gochat/src/message"
)

func (m *Model) joinRoom(room string) {
	// Update the client's current room locally
	m.currentRoom = room
	m.input.Placeholder = "Message " + "[" + m.currentRoom + "] :"
	var joinMessage message.Message

	joinMessage.Author = m.username
	joinMessage.Content = fmt.Sprintf("/join %s", m.currentRoom)
	joinMessage.Type = message.CommandMessage

	bytes, err := message.EncodeMessage(joinMessage)
	if err != nil {
		return
	}

	// Create the /join command to send to the server

	// Send the /join message to the server over WebSocket
	if m.conn != nil {
		err := m.conn.WriteMessage(websocket.BinaryMessage, bytes)
		if err != nil {
			fmt.Println("Error sending join message:", err)
		}
	}

	return
}
