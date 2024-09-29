package teamodel

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/m1kkY8/gochat/src/message"
)

func (m *Model) leave() {
	// Update the client's current room locally
	m.currentRoom = ""
	m.input.Placeholder = "Join any room to start typing"
	m.clear()
	var leaveMessage message.Message

	leaveMessage.Type = message.CommandMessage
	leaveMessage.Author = m.username
	leaveMessage.Content = "/leave"

	bytes, err := message.EncodeMessage(leaveMessage)
	if err != nil {
		return
	}

	if m.conn != nil {
		err := m.conn.WriteMessage(websocket.BinaryMessage, bytes)
		if err != nil {
			fmt.Println("Error sending leave message:", err)
		}
	}

	return
}
