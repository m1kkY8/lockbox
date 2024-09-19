package connection

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack/v5"
)

type Handshake struct {
	Username  string `msgpack:"username"`
	Color     string `msgmsgpack:"color"`
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

func SendHandshake(conn *websocket.Conn, handshake Handshake) error {
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
