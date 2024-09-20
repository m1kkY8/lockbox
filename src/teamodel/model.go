package teamodel

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"

	"github.com/m1kkY8/gochat/src/comps"
	"github.com/m1kkY8/gochat/src/message"
	"github.com/m1kkY8/gochat/src/styles"
)

const useHighPerformanceRenderer = false

type Model struct {
	Input           textinput.Model
	Viewport        comps.Model
	OnlineUsers     comps.Model
	Styles          *styles.Styles
	Width           int
	Height          int
	Conn            *websocket.Conn
	Username        string
	UserColor       string
	MessageChannel  chan string
	OnlineUsersChan chan []string
	Messages        string
}

func listenForMessages(m Model) tea.Cmd {
	return func() tea.Msg {
		return <-m.MessageChannel
	}
}

func listenForOnline(m Model) tea.Cmd {
	return func() tea.Msg {
		return <-m.OnlineUsersChan
	}
}

func (m Model) RecieveMessages() {
	for {
		_, byteMessage, err := m.Conn.ReadMessage()
		if err != nil {
			return
		}

		decodedMsg, err := message.DecodeMessage(byteMessage)
		if err != nil {
			log.Println("Failed decoding")
			continue
		}

		if decodedMsg.Type == message.ServerMessage {
			m.OnlineUsersChan <- strings.Split(decodedMsg.Content, " ")
			continue
		} else {
			formattedMessage := message.Format(decodedMsg)
			m.MessageChannel <- formattedMessage

		}

	}
}

func New(color string, username string, conn *websocket.Conn) *Model {
	styles := styles.DefaultStyle(color)
	input := textinput.New()
	input.Prompt = ""
	input.Placeholder = "Message: "
	input.Width = 50
	input.Focus()

	vp := comps.New(50, 20)
	vp.SetContent("Welcome, start messaging")

	onlineList := comps.New(20, 20)
	onlineList.SetContent("online")

	return &Model{
		Conn:            conn,
		UserColor:       color,
		Username:        username,
		Input:           input,
		Styles:          styles,
		Viewport:        vp,
		OnlineUsers:     onlineList,
		Messages:        "",
		MessageChannel:  make(chan string),
		OnlineUsersChan: make(chan []string),
	}
}

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

			// na enter se salju poruke
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
				timestamp := time.Now().Format(time.TimeOnly)

				var kurac message.Message
				kurac.Author = m.Styles.SenderStyle.Render(m.Username)
				kurac.Timestamp = m.Styles.SenderStyle.Render(timestamp)
				kurac.Content = v + "\r\n"
				kurac.To = ""

				byteMessage, err := message.EncodeMessage(kurac)
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

	case string:
		m.Messages = fmt.Sprintf("%s%s", m.Messages, msg)
		m.Viewport.SetContent(m.Messages)
		m.Viewport.GotoBottom()
		return m, listenForMessages(m)

	case []string:

		var lista []string
		for _, name := range msg {
			tokens := strings.Split(name, ":")
			lista = append(lista, lipgloss.NewStyle().
				Foreground(lipgloss.Color(tokens[0])).
				Render(tokens[1]))
		}

		title := lipgloss.
			NewStyle().
			Bold(true).
			Italic(true).
			Foreground(lipgloss.Color("40")).
			Render("Online:") + "\n"

		m.OnlineUsers.SetContent(title + strings.Join(lista, "\n"))
		return m, listenForOnline(m)
	}

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.Input, cmd = m.Input.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
