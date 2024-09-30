package connection

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func generateUUID() string {
	return uuid.New().String()
}

func ConnectToServer(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
