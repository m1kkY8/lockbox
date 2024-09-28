package teamodel

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// returns a string containg the message
func (m *Model) listenForMessages() tea.Cmd {
	return func() tea.Msg {
		return <-m.MessageChannel
	}
}

// Handles messages displayed in model
func (m *Model) receiver(msg string) {
	m.MessageList.Messages = append(m.MessageList.Messages, msg)
	m.MessageList.Count++

	// if there are more messages than limit pop the oldest from array
	if m.MessageList.Count > messageLimit {
		m.MessageList.Messages = m.MessageList.Messages[1:]
		m.MessageList.Count--
	}

	m.Viewport.SetContent(strings.Join(m.MessageList.Messages, "\n"))
	m.Viewport.GotoBottom()
}

// returns a []string of online users
func (m *Model) listenForOnlineUsers() tea.Cmd {
	return func() tea.Msg {
		return <-m.OnlineUsersChan
	}
}

// Handles array of online users displayed in model
func (m *Model) onlineReceiver(msg []string) {
	for i, name := range msg {
		tokens := strings.Split(name, ":")
		msg[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(tokens[0])).Render(tokens[1])
	}

	title := m.Styles.OnlineTitle.Render("Online:") + "\n"

	m.OnlineUsers.SetContent(title + strings.Join(msg, "\n"))
}

// Clears all messages
func (m *Model) clear() {
	m.MessageList.Messages = nil
	m.MessageList.Count = 0
	m.Viewport.SetContent(strings.Join(m.MessageList.Messages, "\n"))
	m.Viewport.GotoBottom()
}
