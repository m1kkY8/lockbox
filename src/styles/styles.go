package styles

import (
	"math/rand/v2"

	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Border      lipgloss.Style
	SenderStyle lipgloss.Style
}

func GenerateRandomANSIColor() int {
	// Seed the random number generator to ensure different results each time
	// ANSI 8-bit colors range from 0 to 255
	return rand.IntN(256)
}

func DefaultStyle(userColor string) *Styles {
	border := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))

	senderStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(userColor))

	return &Styles{
		Border:      border,
		SenderStyle: senderStyle,
	}
}
