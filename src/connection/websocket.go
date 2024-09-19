package connection

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/m1kkY8/gochat/src/config"
	"github.com/vmihailenco/msgpack/v5"
)

type Handshake struct {
	Username  string `msgpack:"username"`
	Color     string `msgpack:"color"`
	ClientId  string `msgpack:"client_id"`
	PublicKey string `msgpack:"pubkey"`
}

func ConnectToServer(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func GenerateUUID() string {
	return uuid.New().String()
}

func SendHandshake(conn *websocket.Conn, c config.Config) error {
	var handshake Handshake
	handshake.Username = c.Username
	handshake.ClientId = GenerateUUID()
	handshake.PublicKey = "kljuc"
	handshake.Color = c.Color

	bytesHandshake, err := msgpack.Marshal(handshake)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, bytesHandshake)
	if err != nil {
		return err
	}

	return nil
}
