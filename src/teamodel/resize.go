package teamodel

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) resize(msg tea.WindowSizeMsg) {
	currWidth := msg.Width
	currHeight := msg.Height
	m.input.Width = currWidth - 5
	m.width = currWidth
	m.height = currHeight
	m.viewport.Height = currHeight - 5
	m.viewport.Width = currWidth - (currWidth / 5) - 1
	m.onlineUsers.Width = (currWidth / 5) - 5
	m.onlineUsers.Height = currHeight - 5
}
