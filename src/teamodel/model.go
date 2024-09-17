package teamodel

import (
	"log"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"

	"github.com/m1kkY8/gochat/src/notification"
	"github.com/m1kkY8/gochat/src/styles"
	"github.com/m1kkY8/gochat/src/util"
)

type Model struct {
	Messages       []string
	Input          textinput.Model
	Viewport       viewport.Model
	Styles         *styles.Styles
	Width          int
	Height         int
	Conn           *websocket.Conn
	MessageChannel chan string
	Username       string
	UserColor      string
}

var mutex sync.Mutex

func listenForMessages(m Model) tea.Cmd {
	return func() tea.Msg {
		return <-m.MessageChannel // Block until a message is received
	}
}

func (m Model) HandleIncomingMessage() {
	for {
		_, message, err := m.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("Unexpected error reading message: %v", err)
			} else {
				log.Println("Connection closed gracefully, stopping message handling.")
			}
			return
		}
		m.MessageChannel <- string(message)

		notification.Notify(string(message), m.Username)
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

	return &Model{
		Conn:           conn,
		UserColor:      color,
		Username:       username,
		Input:          input,
		Styles:         styles,
		Viewport:       vp,
		Messages:       []string{},
		MessageChannel: make(chan string),
	}
}

func (m Model) Init() tea.Cmd {
	go m.HandleIncomingMessage()
	return listenForMessages(m)
}

func (m Model) View() string {
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Bottom,
		lipgloss.JoinVertical(
			lipgloss.Bottom,
			m.Styles.Border.Render(m.Viewport.View()),
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
		m.Input.Width = msg.Width - 5
		m.Viewport.Width = msg.Width - 4
		m.Viewport.Height = msg.Height - 5

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

				formatedMessage := util.Format(v, m.Username, m.Styles.SenderStyle)

				err := m.Conn.WriteMessage(websocket.TextMessage, []byte(formatedMessage))
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
	}

	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}
