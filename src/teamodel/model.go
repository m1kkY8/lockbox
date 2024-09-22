package teamodel

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/gorilla/websocket"

	"github.com/m1kkY8/gochat/src/comps"
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
	onlineList.SetContent("Online")

	return &Model{
		Conn:            conn,
		UserColor:       color,
		Username:        username,
		Input:           input,
		Styles:          styles,
		Viewport:        vp,
		OnlineUsers:     onlineList,
		MessageList:     &MessageList{},
		MessageChannel:  make(chan string),
		OnlineUsersChan: make(chan []string),
	}
}
