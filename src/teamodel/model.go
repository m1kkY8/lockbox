package teamodel

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"

	"github.com/m1kkY8/gochat/src/message"
	"github.com/m1kkY8/gochat/src/styles"
)

type Model struct {
	Input           textinput.Model
	Viewport        viewport.Model
	OnlineUsers     viewport.Model
	Styles          *styles.Styles
	Width           int
	Height          int
	Conn            *websocket.Conn
	Username        string
	UserColor       string
	MessageChannel  chan string
	Jebmti          []string
	OnlineUsersChan chan []string
	Messages        []string
}

var mutex sync.Mutex

type Users struct {
	Content []string `msgpack:"content"`
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

	vp := viewport.New(50, 20)
	vp.SetContent("Welcome, start messaging")

	onlineList := viewport.New(20, 20)
	onlineList.SetContent("online")

	return &Model{
		Conn:            conn,
		UserColor:       color,
		Username:        username,
		Input:           input,
		Styles:          styles,
		Viewport:        vp,
		OnlineUsers:     onlineList,
		Messages:        []string{},
		MessageChannel:  make(chan string),
		OnlineUsersChan: make(chan []string),
	}
}

func (m Model) Init() tea.Cmd {
	go m.RecieveMessages()
	return tea.Batch(listenForMessages(m), listenForOnline(m))
	// return nil
}

func (m Model) View() string {
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Top,
		lipgloss.Right,
		lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.Styles.Border.Render(m.Viewport.View()),
				m.Styles.Border.Render(m.OnlineUsers.View()),
			),
			// m.Styles.Border.Render(m.OnlineUsers.View()),
			m.Styles.Border.Render(m.Input.View()),
		),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

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
				kurac.Content = v
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
		m.Messages = append(m.Messages, msg)
		m.Viewport.SetContent(strings.Join(m.Messages, "\n"))
		m.Viewport.GotoBottom()
		return m, listenForMessages(m)

	case []string:

		m.OnlineUsers.SetContent(strings.Join(msg, "\n"))
		// m.Messages = append(m.Messages, amogus)
		// m.Viewport.SetContent(strings.Join(m.Messages, "\n"))
		// m.Viewport.GotoBottom()
		return m, listenForOnline(m)
	}

	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}
