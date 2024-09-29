package message

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	ServerMessage  = 1
	ChatMessage    = 2
	CommandMessage = 3
)

type Message struct {
	Type      int    `msgpack:"type"`
	Author    string `msgpack:"author"`
	Content   string `msgpack:"content"`
	Room      string `msgpack:"room"`
	To        string `msgpack:"to"`
	Timestamp string `msgpack:"timestamp"`
	Color     string `msgpack:"color"`
}

// Parses Message into string format `timestamp user: message`
func Format(message Message) string {
	username := lipgloss.NewStyle().Foreground(lipgloss.Color(message.Color)).Render(message.Author)
	timestamp := lipgloss.NewStyle().Foreground(lipgloss.Color(message.Color)).Render(message.Timestamp)
	content := message.Content
	return fmt.Sprintf("%s %s: %s", timestamp, username, content)
}

func FormatWhisper(message Message) string {
	username := lipgloss.NewStyle().Foreground(lipgloss.Color(message.Color)).Render(message.Author)
	timestamp := lipgloss.NewStyle().Foreground(lipgloss.Color(message.Color)).Render(message.Timestamp)
	content := message.Content

	return fmt.Sprintf("%s (Whisper from %s): %s ", timestamp, username, content)
}

func EncodeMessage(message Message) ([]byte, error) {
	encodedMessage, err := msgpack.Marshal(message)
	if err != nil {
		return nil, err
	}

	return encodedMessage, nil
}

// Unpack []byte into Message
func DecodeMessage(byteMessage []byte) (Message, error) {
	var decodedMessage Message

	err := msgpack.Unmarshal(byteMessage, &decodedMessage)
	if err != nil {
		return Message{}, err
	}

	return decodedMessage, nil
}
