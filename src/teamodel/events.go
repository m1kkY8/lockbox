package teamodel

import (
	"fmt"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	"github.com/m1kkY8/gochat/src/message"
)

var limit = 100

func (m Model) Init() tea.Cmd {
	go m.RecieveMessages()
	return tea.Batch(listenForMessages(m), listenForOnline(m))
}

func (m Model) View() string {
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	// resizing the terminal
	case tea.WindowSizeMsg:
		currWidth := msg.Width
		currHeight := msg.Height
		m.Input.Width = currWidth - 5
		m.Width = currWidth
		m.Height = currHeight
		m.Viewport.Height = currHeight - 5
		m.Viewport.Width = currWidth - (currWidth / 5) - 1
		m.OnlineUsers.Width = (currWidth / 5) - 5
		m.OnlineUsers.Height = currHeight - 5

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q", "ctrl+c":
			return m, tea.Quit

		case "ctrl+l":
			m.MessageList.Messages = nil
			m.MessageList.Count = 0
			m.Viewport.SetContent(strings.Join(m.MessageList.Messages, "\n"))
			m.Viewport.GotoBottom()
			// Enter to send messages
		case "enter":
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
					if m.MessageList.Count > limit {
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

				byteMessage, err := message.EncodeMessage(userMessage)
				if err != nil {
					log.Println("Failed encoding message")
					break
				}

				err = m.Conn.WriteMessage(websocket.BinaryMessage, byteMessage)
				if err != nil {
					log.Printf("Error writing message: %v", err)
					break
				}
			}
			return m, nil
		}

		// Handle displaying user messages
	case string:
		// add new message and increase count
		m.MessageList.Messages = append(m.MessageList.Messages, msg)
		m.MessageList.Count++

		// if there are more messages than limit pop the oldest from array
		if m.MessageList.Count > limit {
			m.MessageList.Messages = m.MessageList.Messages[1:]
			m.MessageList.Count--
		}

		m.Viewport.SetContent(strings.Join(m.MessageList.Messages, "\n"))
		m.Viewport.GotoBottom()
		return m, listenForMessages(m)

		// Handling meesage from servr containing online users
		// Parse ever user and get color of that user
	case []string:
		// message format is color:username
		for i, name := range msg {
			tokens := strings.Split(name, ":")
			msg[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(tokens[0])).Render(tokens[1])
		}

		title := m.Styles.OnlineTitle.Render("Online:") + "\n"

		m.OnlineUsers.SetContent(title + strings.Join(msg, "\n"))
		return m, listenForOnline(m)
	}

	// Every other unhandled keystroke or mouse movement is sent to Viewport and Input
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.Input, cmd = m.Input.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
