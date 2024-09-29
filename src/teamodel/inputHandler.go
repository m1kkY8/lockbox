package teamodel

import "strings"

func (m *Model) commandHandler(commandString string) {
	commandParts := strings.Split(commandString, " ")

	if len(commandParts) < 2 {
		return
	}

	command := commandParts[0]

	switch command {
	case "/join":
		m.joinRoom(commandParts[1])
	}
}
