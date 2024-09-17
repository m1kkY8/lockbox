package util

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func Format(content string, username string, style lipgloss.Style) string {
	timestamp := time.Now().Format(time.TimeOnly)
	usr := style.Render(timestamp + " " + username)
	formatted := fmt.Sprintf("%s: %s", usr, content)

	return formatted
}
