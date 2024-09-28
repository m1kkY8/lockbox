package teamodel

import (
	"time"

	"github.com/m1kkY8/gochat/src/message"
)

func (m *Model) createMessage(content string) message.Message {
	timestamp := time.Now().Format(time.TimeOnly)

	var userMessage message.Message
	userMessage.Author = m.username
	userMessage.Timestamp = m.styles.SenderStyle.Render(timestamp)
	userMessage.To = "all"
	userMessage.Content = content

	return userMessage
}
