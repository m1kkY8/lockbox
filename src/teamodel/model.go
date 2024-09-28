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
	MessageList     *MessageList
}

type MessageList struct {
	Messages []string
	Count    int
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
		Conn:            conn,
		UserColor:       conf.Color,
		Username:        conf.Username,
		Input:           input,
		Styles:          styles,
		Viewport:        vp,
		OnlineUsers:     onlineList,
		MessageList:     &MessageList{},
		MessageChannel:  make(chan string),
		OnlineUsersChan: make(chan []string),
	}
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

func (m *Model) Init() tea.Cmd {
	go m.RecieveMessages()
	return tea.Batch(
		m.listenForMessages(),
		m.listenForOnlineUsers(),
	)
}
