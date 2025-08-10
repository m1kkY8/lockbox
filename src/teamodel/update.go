package teamodel

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m1kkY8/lockbox/src/config"
)

// Update handles Bubble Tea messages and updates the model state
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleResize(msg)

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case string:
		// Handle displaying user messages
		m.displayMessages(msg)
		return m, m.listenForMessages()

	case []string:
		// Handle message from server containing online users
		m.displayOnlineUsers(msg)
		return m, m.listenForOnlineUsers()
	}

	// Update UI components
	m.ui.Viewport, cmd = m.ui.Viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.ui.Input, cmd = m.ui.Input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// handleKeyPress processes keyboard input
func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+q", "ctrl+c":
		return m, tea.Quit

	case "ctrl+l":
		m.clear()
		return m, nil

	case "enter":
		return m.handleEnterKey()

	default:
		// Let the input component handle other keys
		var cmd tea.Cmd
		m.ui.Input, cmd = m.ui.Input.Update(msg)
		return m, cmd
	}
}

// handleEnterKey processes the enter key for sending messages or commands
func (m *Model) handleEnterKey() (tea.Model, tea.Cmd) {
	content := m.ui.Input.Value()
	if content == "" {
		return m, nil
	}

	// Reset input first
	m.ui.Input.Reset()

	// Handle exit command
	if content == config.ExitCommand {
		return m, tea.Quit
	}

	// Handle other commands
	if strings.HasPrefix(content, "/") {
		m.commandHandler(content)
		return m, nil
	}

	// Send regular message
	return m.sendRegularMessage(content)
}

// sendRegularMessage handles sending a regular chat message
func (m *Model) sendRegularMessage(content string) (tea.Model, tea.Cmd) {
	// Check if connected
	if m.client.Conn == nil {
		return m, nil
	}

	// Check if in a room
	if m.state.CurrentRoom == "" {
		return m, nil
	}

	userMessage := m.createMessage(content)
	if err := m.sendMessage(userMessage); err != nil {
		// TODO: Show error to user
		return m, nil
	}

	return m, nil
}

// handleResize updates the model dimensions when the terminal is resized
func (m *Model) handleResize(msg tea.WindowSizeMsg) {
	m.resize(msg)
}
