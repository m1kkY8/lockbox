package teamodel

import (
	"strings"

	"github.com/m1kkY8/lockbox/src/config"
)

// commandHandler processes chat commands
func (m *Model) commandHandler(commandString string) {
	commandParts := strings.Split(commandString, " ")
	if len(commandParts) == 0 {
		return
	}

	command := commandParts[0]

	switch command {
	case config.JoinCommand:
		if len(commandParts) < 2 {
			return
		}
		m.joinRoom(commandParts[1])
	case config.LeaveCommand:
		m.leave()
	default:
		// Unknown command, could show error message to user
	}
}
