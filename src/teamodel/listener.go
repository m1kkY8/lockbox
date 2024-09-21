package teamodel

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m1kkY8/gochat/src/message"
)

func (m Model) RecieveMessages() {
	for {
		_, byteMessage, err := m.Conn.ReadMessage()
		if err != nil {
			return
		}

		decodedMsg, err := message.DecodeMessage(byteMessage)
		if err != nil {
			log.Println("Failed decoding")
			continue
		}

		if decodedMsg.Type == message.ServerMessage {
			m.OnlineUsersChan <- strings.Split(decodedMsg.Content, " ")
		} else {
			formattedMessage := message.Format(decodedMsg)
			m.MessageChannel <- formattedMessage

		}

	}
}

// returns a string containg the message
func listenForMessages(m Model) tea.Cmd {
	return func() tea.Msg {
		return <-m.MessageChannel
	}
}

// returns a []string of online users
func listenForOnline(m Model) tea.Cmd {
	return func() tea.Msg {
		return <-m.OnlineUsersChan
	}
}
