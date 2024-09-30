package message

import (
	"crypto/rsa"

	"github.com/m1kkY8/gochat/src/encryption"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	ServerMessage  = 1
	ChatMessage    = 2
	CommandMessage = 3
	KeyMessage     = 4
)

type Message struct {
	Type      int                 `msgpack:"type"`
	Author    string              `msgpack:"author"`
	Content   string              `msgpack:"content"`
	Room      string              `msgpack:"room"`
	To        string              `msgpack:"to"`
	Timestamp string              `msgpack:"timestamp"`
	Color     string              `msgpack:"color"`
	AESKeys   []encryption.AESKey `msgpack:"aes_keys"`
}

type PublicKeys struct {
	Type       int              `msgpack:"type"`
	PublicKeys []*rsa.PublicKey `msgpack:"public_keys"`
}

// Pack message into []byte
func EncodeMessage(message Message) ([]byte, error) {
	encodedMessage, err := msgpack.Marshal(message)
	if err != nil {
		return nil, err
	}

	return encodedMessage, nil
}

// Unpack []byte into Message
func DecodeMessage(byteMessage []byte) (Message, error) {
	var decodedMessage Message

	err := msgpack.Unmarshal(byteMessage, &decodedMessage)
	if err != nil {
		return Message{}, err
	}

	return decodedMessage, nil
}
