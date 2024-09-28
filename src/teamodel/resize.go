package teamodel

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) resize(msg tea.WindowSizeMsg) {
	currWidth := msg.Width
	currHeight := msg.Height
	m.Input.Width = currWidth - 5
	m.Width = currWidth
	m.Height = currHeight
	m.Viewport.Height = currHeight - 5
	m.Viewport.Width = currWidth - (currWidth / 5) - 1
	m.OnlineUsers.Width = (currWidth / 5) - 5
	m.OnlineUsers.Height = currHeight - 5
}
