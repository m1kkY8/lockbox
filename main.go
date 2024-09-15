package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"flag"
	"fmt"
	"log"
	_ "net/url"
	"strings"
	_ "time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

// Add a purple, rectangular border

type styles struct {
	border      lipgloss.Style
	senderStyle lipgloss.Style
}

type model struct {
	messages       []string
	input          textinput.Model
	view           viewport.Model
	styles         *styles
	width          int
	height         int
	conn           *websocket.Conn
	messageChannel chan string
	username       string
}

type message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	To       string `json:"to"`
}

func DefaultStyle() *styles {
	border := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))

	return &styles{
		border:      border,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
	}
}

func CreateWebSocketConnection(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func listenForMessages(m *model) tea.Cmd {
	return func() tea.Msg {
		return <-m.messageChannel // Block until a message is received
	}
}

func (m *model) HandleIncomingMessage() {
	for {
		_, message, err := m.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			continue
		}
		m.messageChannel <- string(message)
	}
}

func New() *model {
	styles := DefaultStyle()
	input := textinput.New()
	input.Prompt = ""
	input.Placeholder = "Message: "
	input.Width = 50
	input.Focus()

	vp := viewport.New(50, 20)
	vp.SetContent("Welcome, start messaging")

	return &model{
		messages:       []string{},
		input:          input,
		styles:         styles,
		view:           vp,
		messageChannel: make(chan string),
	}
}

func (m *model) Init() tea.Cmd {
	go m.HandleIncomingMessage()
	return nil
}

func (m *model) View() string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Bottom,
		lipgloss.JoinVertical(
			lipgloss.Bottom,
			m.styles.border.Render(m.view.View()),
			m.styles.border.Render(m.input.View()),
		),
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = msg.Width - 5
		m.view.Width = msg.Width - 4
		m.view.Height = msg.Height - 5

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			return m, tea.Quit

			// na enter se salju poruke
		case "enter":
			v := m.input.Value()
			if v == "" {
				return m, nil
			}
			if v == "/quit" {
				return m, tea.Quit
			}

			m.input.Reset()
			if m.conn != nil {
				formated := fmt.Sprintf("%s: %s", m.username, v)
				err := m.conn.WriteMessage(websocket.TextMessage, []byte(formated))
				if err != nil {
					log.Printf("Error writing message: %v", err)
				}
			}
			return m, listenForMessages(m)
		}

	case string:
		// Handle incoming messages from the channel
		m.messages = append(m.messages, msg)
		m.view.SetContent(strings.Join(m.messages, "\n"))
		m.view.GotoBottom()
		return m, listenForMessages(m)
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func main() {
	username := flag.String("u", "anon", "Username")
	flag.Parse()

	url := "ws://139.162.132.8:42069/ws"
	// url := "ws://localhost:42069/ws"
	conn, err := CreateWebSocketConnection(url)
	if err != nil {
		log.Printf("Error creating websocket connection: %v", err)
	}

	defer conn.Close()

	teaModel := New()
	teaModel.conn = conn
	teaModel.username = *username

	p := tea.NewProgram(teaModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}
