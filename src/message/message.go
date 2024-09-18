package message

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

var (
	ServerMessage = 1
	ChatMessage   = 2
)

type Message struct {
	Type      int    `msgpack:"type"`
	Author    string `msgpack:"author"`
	Content   string `msgpack:"content"`
	To        string `msgpack:"to"`
	Timestamp string `msgpack:"timestamp"`
}

func Format(message Message) string {
	timestamp := message.Timestamp
	username := message.Author
	content := message.Content
	return fmt.Sprintf("%s %s: %s", timestamp, username, content)
}

func EncodeMessage(message Message) ([]byte, error) {
	encodedMessage, err := msgpack.Marshal(message)
	if err != nil {
		return nil, err
	}

	return encodedMessage, nil
}

func DecodeMessage(byteMessage []byte) (Message, error) {
	var decodedMessage Message

	err := msgpack.Unmarshal(byteMessage, &decodedMessage)
	if err != nil {
		return Message{}, err
	}

	return decodedMessage, nil
}
