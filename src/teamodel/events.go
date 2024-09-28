package teamodel

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/m1kkY8/gochat/src/message"
)

var messageLimit = 100

func (m *Model) Init() tea.Cmd {
	go m.RecieveMessages()
	return tea.Batch(
		m.listenForMessages(),
		m.listenForOnlineUsers(),
	)
}

func (m *Model) View() string {
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.Styles.Border.Render(m.Viewport.View()),
				m.Styles.Border.Render(m.OnlineUsers.View()),
			),
			m.Styles.Border.Render(m.Input.View()),
		),
	)
}

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
				// userMessage.Author = m.Styles.SenderStyle.Render(m.Username)
				userMessage.Author = m.Username
				userMessage.Timestamp = m.Styles.SenderStyle.Render(timestamp)

				// Check if it's a whisper command
				if strings.HasPrefix(v, "/whisper") {
					whisper := strings.TrimPrefix(v, "/whisper ")
					parts := strings.SplitN(whisper, " ", 2)

					if len(parts) < 2 {
						break
					}

					// Set the target user and content for whisper
					userMessage.To = parts[0]
					userMessage.Content = parts[1]

					newMessage := fmt.Sprintf("%s (Whisper to %s): %s ", userMessage.Timestamp, userMessage.To, parts[1])

					m.MessageList.Messages = append(m.MessageList.Messages, newMessage)
					m.MessageList.Count++

					// if there are more messages than limit pop the oldest from array
					if m.MessageList.Count > messageLimit {
						m.MessageList.Messages = m.MessageList.Messages[1:]
						m.MessageList.Count--
					}

					m.Viewport.SetContent(strings.Join(m.MessageList.Messages, "\n"))
					m.Viewport.GotoBottom()

					//	m.placeMessage(newMessage)
				} else {
					// Normal message
					userMessage.To = "all"
					userMessage.Content = v
				}

				err := m.sendMessage(userMessage)
				if err != nil {
					break
				}

			}
			return m, nil
		}

		// Handle displaying user messages
	case string:
		m.receiver(msg)
		return m, m.listenForMessages()

		// Handling meesage from servr containing online users
		// Parse ever user and get color of that user
	case []string:
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
