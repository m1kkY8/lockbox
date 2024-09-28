package teamodel

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// returns a string containg the message
func (m *Model) listenForMessages() tea.Cmd {
	return func() tea.Msg {
		return <-m.messageChannel
	}
}

// Handles messages displayed in model
func (m *Model) displayMessages(msg string) {
	m.messageList.messages = append(m.messageList.messages, msg)
	m.messageList.count++

	// if there are more messages than limit pop the oldest from array
	if m.messageList.count > messageLimit {
		m.messageList.messages = m.messageList.messages[1:]
		m.messageList.count--
	}

	m.viewport.SetContent(strings.Join(m.messageList.messages, "\n"))
	m.viewport.GotoBottom()
}

// returns a []string of online users
func (m *Model) listenForOnlineUsers() tea.Cmd {
	return func() tea.Msg {
		return <-m.onlineUsersChan
	}
}

// Handles array of online users displayed in model
func (m *Model) displayOnlineUsers(msg []string) {
	for i, name := range msg {
		tokens := strings.Split(name, ":")
		msg[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(tokens[0])).Render(tokens[1])
	}

	title := m.styles.OnlineTitle.Render("Online:") + "\n"

	m.onlineUsers.SetContent(title + strings.Join(msg, "\n"))
}

// Clears all messages
func (m *Model) clear() {
	m.messageList.messages = nil
	m.messageList.count = 0
	m.viewport.SetContent(strings.Join(m.messageList.messages, "\n"))
	m.viewport.GotoBottom()
}
