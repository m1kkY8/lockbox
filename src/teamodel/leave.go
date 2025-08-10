package teamodel

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/m1kkY8/lockbox/src/message"
)

// leave handles leaving the current chat room
func (m *Model) leave() {
	// Update the client's current room locally
	m.state.CurrentRoom = ""
	m.ui.Input.Placeholder = "Join any room to start typing"
	m.clear()

	var leaveMessage message.Message
	leaveMessage.Type = message.CommandMessage
	leaveMessage.Author = m.state.Username
	leaveMessage.Content = "/leave"

	bytes, err := message.EncodeMessage(leaveMessage)
	if err != nil {
		return
	}

	if m.client.Conn != nil {
		err := m.client.Conn.WriteMessage(websocket.BinaryMessage, bytes)
		if err != nil {
			fmt.Println("Error sending leave message:", err)
		}
	}
}
