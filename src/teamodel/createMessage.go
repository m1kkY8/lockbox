package teamodel

import (
	"time"

	"github.com/m1kkY8/gochat/src/message"
)

func (m *Model) createMessage(content string) message.Message {
	timestamp := time.Now().Format(time.TimeOnly)

	var userMessage message.Message
	userMessage.Type = message.ChatMessage
	userMessage.Author = m.username
	userMessage.Timestamp = timestamp
	userMessage.To = m.currentRoom
	userMessage.Content = content
	userMessage.Color = m.userColor

	return userMessage
}
