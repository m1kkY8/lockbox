package teamodel

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m1kkY8/gochat/src/message"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// resizing the terminal
		m.resize(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q", "ctrl+c":
			return m, tea.Quit

		case "ctrl+l":
			m.clear()
		case "enter":
			// Enter to send messages
			v := m.Input.Value()
			if v == "" {
				return m, nil
			}
			if v == ":q" {
				return m, tea.Quit
			}

			m.Input.Reset()
			if m.Conn != nil {
				// Create Message object and send it to the server
				timestamp := time.Now().Format(time.TimeOnly)

				var userMessage message.Message
				userMessage.Author = m.Username
				userMessage.Timestamp = m.Styles.SenderStyle.Render(timestamp)
				userMessage.To = "all"
				userMessage.Content = v

				err := m.sendMessage(userMessage)
				if err != nil {
					break
				}

			}
			return m, nil
		}

	case string:
		// Handle displaying user messages
		m.receiver(msg)
		return m, m.listenForMessages()

	case []string:
		// Handling meesage from servr containing online users
		m.onlineReceiver(msg)
		return m, m.listenForOnlineUsers()
	}

	// Every other unhandled keystroke or mouse movement is sent to Viewport and Input
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.Input, cmd = m.Input.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
