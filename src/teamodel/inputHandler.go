package teamodel

import "strings"

func (m *Model) commandHandler(commandString string) {
	commandParts := strings.Split(commandString, " ")

	command := commandParts[0]

	switch command {
	case "/join":
		if len(commandParts) < 2 {
			return
		}
		m.joinRoom(commandParts[1])
	case "/leave":
		m.leave()
	}
}
