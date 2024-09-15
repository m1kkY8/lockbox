package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"log"
	_ "net/url"
	"strings"
	"time"

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
	messages []string
	input    textinput.Model
	view     viewport.Model
	styles   *styles
	width    int
	height   int
	conn     *websocket.Conn
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
		messages: []string{},
		input:    input,
		styles:   styles,
		view:     vp,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			// TODO: Ne apdejtovati veiwport ovde
			// Napraviti go rutinu koja ce da slusa za poruke i sama apdejtuje viewport nezavisno od svega
			m.messages = append(m.messages, m.styles.senderStyle.Render(time.Now().Format(time.TimeOnly)+" You: ")+v)
			m.view.SetContent(strings.Join(m.messages, "\n"))
			m.input.Reset()
			m.view.GotoBottom()

			if m.conn != nil {
				// wsMsg := message{
				// 	Username: "test",
				// 	Message:  v,
				// 	To:       "",
				// }
				err := m.conn.WriteMessage(websocket.TextMessage, []byte(v))
				if err != nil {
					log.Printf("Error writing message: %v", err)
				}
			}
			return m, nil
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func main() {
	url := "ws://localhost:42069/ws"
	conn, err := CreateWebSocketConnection(url)
	if err != nil {
		log.Printf("Error creating websocket connection: %v", err)
	}

	defer conn.Close()

	teaModel := New()
	teaModel.conn = conn

	p := tea.NewProgram(teaModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}
