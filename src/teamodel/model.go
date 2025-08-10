package teamodel

import (
	"crypto/rsa"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"

	"github.com/m1kkY8/lockbox/src/comps"
	"github.com/m1kkY8/lockbox/src/config"
	"github.com/m1kkY8/lockbox/src/encryption"
	"github.com/m1kkY8/lockbox/src/styles"
)

// Constants for UI configuration
const (
	MessageLimit        = 100
	DefaultInputWidth   = 50
	DefaultViewportSize = 20
	DefaultOnlineWidth  = 20
)

// Default messages
const (
	WelcomeMessage      = "Welcome, start messaging"
	OnlineUsersTitle    = "Online"
	JoinRoomPlaceholder = "Join any room to start typing"
)

// UIComponents groups all UI-related components
type UIComponents struct {
	Width       int
	Height      int
	Input       textinput.Model
	Viewport    comps.Model
	OnlineUsers comps.Model
	Styles      *styles.Styles
}

// ChatState holds the current state of the chat session
type ChatState struct {
	Username    string
	UserColor   string
	CurrentRoom string
	MessageList *MessageList
}

// ChatClient handles network and cryptographic operations
type ChatClient struct {
	Conn       *websocket.Conn
	KeyPair    *encryption.RSAKeys
	PublicKeys []*rsa.PublicKey
}

// MessageChannels groups all communication channels
type MessageChannels struct {
	Messages    chan string
	OnlineUsers chan []string
	PublicKeys  chan []*rsa.PublicKey
}

// Model represents the main TUI model for the chat application
type Model struct {
	// UI components
	ui UIComponents

	// Chat state
	state ChatState

	// Network and crypto
	client ChatClient

	// Communication channels
	channels MessageChannels
}

// MessageList represents a collection of chat messages
type MessageList struct {
	messages []string
	count    int
}

// New creates a new Model instance with the provided configuration
func New(conf config.Config, conn *websocket.Conn, keyPair *encryption.RSAKeys) *Model {
	appStyles := styles.DefaultStyle(conf.Color)

	input := textinput.New()
	input.Prompt = ""
	input.Placeholder = JoinRoomPlaceholder
	input.Width = DefaultInputWidth
	input.Focus()

	vp := comps.New(DefaultInputWidth, DefaultViewportSize)
	vp.SetContent(WelcomeMessage)

	onlineList := comps.New(DefaultOnlineWidth, DefaultViewportSize)
	onlineList.SetContent(OnlineUsersTitle)

	return &Model{
		ui: UIComponents{
			Input:       input,
			Styles:      appStyles,
			Viewport:    vp,
			OnlineUsers: onlineList,
		},
		state: ChatState{
			UserColor:   conf.Color,
			Username:    conf.Username,
			CurrentRoom: "",
			MessageList: &MessageList{},
		},
		client: ChatClient{
			Conn:       conn,
			KeyPair:    keyPair,
			PublicKeys: make([]*rsa.PublicKey, 0),
		},
		channels: MessageChannels{
			Messages:    make(chan string),
			OnlineUsers: make(chan []string),
			PublicKeys:  make(chan []*rsa.PublicKey),
		},
	}
}

// View renders the TUI layout
func (m *Model) View() string {
	return lipgloss.Place(
		m.ui.Width,
		m.ui.Height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.ui.Styles.Border.Render(m.ui.Viewport.View()),
				m.ui.Styles.Border.Render(m.ui.OnlineUsers.View()),
			),
			m.ui.Styles.Border.Render(m.ui.Input.View()),
		),
	)
}

// Init initializes the model and starts background processes
func (m *Model) Init() tea.Cmd {
	go m.recieveMessages() // Keep the original method name for now
	return tea.Batch(
		m.listenForMessages(),
		m.listenForOnlineUsers(),
	)
}

// Getter methods for backward compatibility and clean access
func (m *Model) Username() string {
	return m.state.Username
}

func (m *Model) CurrentRoom() string {
	return m.state.CurrentRoom
}

func (m *Model) SetCurrentRoom(room string) {
	m.state.CurrentRoom = room
}

func (m *Model) Connection() *websocket.Conn {
	return m.client.Conn
}

func (m *Model) KeyPair() *encryption.RSAKeys {
	return m.client.KeyPair
}

func (m *Model) PublicKeys() []*rsa.PublicKey {
	return m.client.PublicKeys
}

func (m *Model) SetPublicKeys(keys []*rsa.PublicKey) {
	m.client.PublicKeys = keys
}

func (m *Model) MessageChannel() chan string {
	return m.channels.Messages
}

func (m *Model) OnlineUsersChannel() chan []string {
	return m.channels.OnlineUsers
}

func (m *Model) Input() *textinput.Model {
	return &m.ui.Input
}

func (m *Model) Viewport() *comps.Model {
	return &m.ui.Viewport
}

func (m *Model) OnlineUsers() *comps.Model {
	return &m.ui.OnlineUsers
}

func (m *Model) Styles() *styles.Styles {
	return m.ui.Styles
}

func (m *Model) SetDimensions(width, height int) {
	m.ui.Width = width
	m.ui.Height = height
}

func (m *Model) MessageList() *MessageList {
	return m.state.MessageList
}
