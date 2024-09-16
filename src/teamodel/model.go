package teamodel

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
	"regexp"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	"github.com/m1kkY8/gochat/src/styles"
	"github.com/gen2brain/beeep"
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

func listenForMessages(m *Model) tea.Cmd {
	return func() tea.Msg {
		return <-m.MessageChannel // Block until a message is received
	}
}

func (m *Model) HandleIncomingMessage() {
	for {
		_, message, err := m.Conn.ReadMessage()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok {
				log.Printf("Network operation error: %v (Op: %v, Net: %v, Addr: %v)", opErr.Err, opErr.Op, opErr.Net, opErr.Addr)
			}

			// Log detailed error information
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Connection closed unexpectedly: %v", err)
				return
			} else if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("Connection closed normally: %v", err)
				return
			} else {
				log.Printf("Read error: %v (type: %T)", err, err)
			}
			return
		}
		// Successfully received a message
		m.MessageChannel <- string(message)
		Notify(m,string(message))
	}
}

func Notify(m *Model, msg string){
	reg := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	cleanedMessage := reg.ReplaceAllString(msg, "")
	partMsg := strings.SplitN(cleanedMessage, ": ",2)
	getPartUser := strings.Split(partMsg[0]," ")
	fromUser := getPartUser[1]

	if(fromUser == m.Username){
		return
	}

	formatMsg := partMsg[1]
	err := beeep.Notify(fromUser ,formatMsg , "assets/amogus.png")
	if err != nil {
    		panic(err)
	}
}

func New(color string) *Model {
	styles := styles.DefaultStyle(color)
	input := textinput.New()
	input.Prompt = ""
	input.Placeholder = "Message: "
	input.Width = 50
	input.Focus()

	vp := viewport.New(50, 20)
	vp.SetContent("Welcome, start messaging")

	return &Model{
		Messages:       []string{},
		Input:          input,
		Styles:         styles,
		Viewport:       vp,
		MessageChannel: make(chan string),
	}
}

func (m *Model) Init() tea.Cmd {
	go m.HandleIncomingMessage()
	return listenForMessages(m)
}

func (m *Model) View() string {
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

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "ctrl+q":
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
				usr := m.Styles.SenderStyle.Render(timestamp + " " + m.Username)
				formated := fmt.Sprintf("%s: %s", usr, v)

				err := m.Conn.WriteMessage(websocket.TextMessage, []byte(formated))
				if err != nil {
					log.Printf("Error writing message: %v", err)
					break
				}
			}
			return m, listenForMessages(m)
		}

	case string:
		// Handle incoming messages from the channel
		m.Messages = append(m.Messages, msg)
		m.Viewport.SetContent(strings.Join(m.Messages, "\n"))
		m.Viewport.GotoBottom()
		return m, listenForMessages(m)
	}

	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m *Model) CloseConnection() {
	if m.Conn != nil {
		err := m.Conn.Close()
		if err != nil {
			log.Printf("Error closing WebSocket connection: %v", err)
		}
	}
}
