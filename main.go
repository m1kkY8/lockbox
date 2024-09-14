package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

		case "enter":
			v := m.input.Value()
			if v == "" {
				return m, nil
			}

			m.messages = append(m.messages, m.styles.senderStyle.Render("You: ")+v)
			m.view.SetContent(strings.Join(m.messages, "\n"))
			m.input.Reset()
			m.view.GotoBottom()
			return m, nil
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	p := tea.NewProgram(New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}
