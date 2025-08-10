package teamodel

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/m1kkY8/lockbox/src/message"
)

// joinRoom handles joining a chat room
func (m *Model) joinRoom(room string) {
	// Update the client's current room locally
	m.state.CurrentRoom = room
	m.ui.Input.Placeholder = "Message " + "[" + m.state.CurrentRoom + "]:"
	m.clear()

	var joinMessage message.Message
	joinMessage.Type = message.CommandMessage
	joinMessage.Author = m.state.Username
	joinMessage.Content = fmt.Sprintf("/join %s", m.state.CurrentRoom)

	bytes, err := message.EncodeMessage(joinMessage)
	if err != nil {
		return
	}

	// Send the join message to the server over WebSocket
	if m.client.Conn != nil {
		err := m.client.Conn.WriteMessage(websocket.BinaryMessage, bytes)
		if err != nil {
			fmt.Println("Error sending join message:", err)
		}
	}
}
