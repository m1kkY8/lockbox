package teamodel

import tea "github.com/charmbracelet/bubbletea"

// resize handles terminal window resize events
func (m *Model) resize(msg tea.WindowSizeMsg) {
	currWidth := msg.Width
	currHeight := msg.Height

	// Update model dimensions
	m.ui.Width = currWidth
	m.ui.Height = currHeight

	// Update input component
	m.ui.Input.Width = currWidth - 5

	// Update viewport dimensions
	m.ui.Viewport.Height = currHeight - 5
	m.ui.Viewport.Width = currWidth - (currWidth / 5) - 1

	// Update online users panel
	m.ui.OnlineUsers.Width = (currWidth / 5) - 5
	m.ui.OnlineUsers.Height = currHeight - 5
}
