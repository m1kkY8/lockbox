package teamodel

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// listenForMessages returns a command that listens for incoming messages
func (m *Model) listenForMessages() tea.Cmd {
	return func() tea.Msg {
		return <-m.channels.Messages
	}
}

// displayMessages handles messages displayed in the UI
func (m *Model) displayMessages(msg string) {
	messageList := m.state.MessageList
	messageList.messages = append(messageList.messages, msg)
	messageList.count++

	// If there are more messages than limit, remove the oldest
	if messageList.count > MessageLimit {
		messageList.messages = messageList.messages[1:]
		messageList.count--
	}

	content := strings.Join(messageList.messages, "\n")
	m.ui.Viewport.SetContent(content)
	m.ui.Viewport.GotoBottom()
}

// listenForOnlineUsers returns a command that listens for online user updates
func (m *Model) listenForOnlineUsers() tea.Cmd {
	return func() tea.Msg {
		return <-m.channels.OnlineUsers
	}
}

// displayOnlineUsers handles the display of online users in the UI
func (m *Model) displayOnlineUsers(users []string) {
	// Apply color styling to each username
	for i, user := range users {
		tokens := strings.Split(user, ":")
		if len(tokens) >= 2 {
			colorCode := tokens[0]
			username := tokens[1]
			users[i] = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorCode)).
				Render(username)
		}
	}

	title := m.ui.Styles.OnlineTitle.Render("Online:") + "\n"
	content := title + strings.Join(users, "\n")
	m.ui.OnlineUsers.SetContent(content)
}

// clear removes all messages from the chat display
func (m *Model) clear() {
	messageList := m.state.MessageList
	messageList.messages = nil
	messageList.count = 0

	content := strings.Join(messageList.messages, "\n")
	m.ui.Viewport.SetContent(content)
	m.ui.Viewport.GotoBottom()
}
