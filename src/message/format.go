package message

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

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
