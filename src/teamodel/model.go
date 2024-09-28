package teamodel

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"

	"github.com/m1kkY8/gochat/src/comps"
	"github.com/m1kkY8/gochat/src/config"
	"github.com/m1kkY8/gochat/src/styles"
)

type Model struct {
	input           textinput.Model
	viewport        comps.Model
	onlineUsers     comps.Model
	styles          *styles.Styles
	width           int
	height          int
	conn            *websocket.Conn
	username        string
	userColor       string
	messageChannel  chan string
	onlineUsersChan chan []string
	messageList     *MessageList
	currentRoom     string
}

type MessageList struct {
	messages []string
	count    int
}

var messageLimit = 100

func New(conf config.Config, conn *websocket.Conn) *Model {
	styles := styles.DefaultStyle(conf.Color)

	input := textinput.New()
	input.Prompt = ""
	input.Placeholder = "Message: "
	input.Width = 50
	input.Focus()

	vp := comps.New(50, 20)
	vp.SetContent("Welcome, start messaging")

	onlineList := comps.New(20, 20)
	onlineList.SetContent("Online")

	return &Model{
		conn:            conn,
		userColor:       conf.Color,
		username:        conf.Username,
		input:           input,
		styles:          styles,
		viewport:        vp,
		onlineUsers:     onlineList,
		messageList:     &MessageList{},
		messageChannel:  make(chan string),
		onlineUsersChan: make(chan []string),
		currentRoom:     "all",
	}
}

func (m *Model) View() string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.styles.Border.Render(m.viewport.View()),
				m.styles.Border.Render(m.onlineUsers.View()),
			),
			m.styles.Border.Render(m.input.View()),
		),
	)
}

func (m *Model) Init() tea.Cmd {
	go m.recieveMessages()
	return tea.Batch(
		m.listenForMessages(),
		m.listenForOnlineUsers(),
	)
}
