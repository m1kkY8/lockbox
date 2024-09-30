package connection

import (
	"crypto/rsa"

	"github.com/gorilla/websocket"
	"github.com/m1kkY8/lockbox/src/config"
	"github.com/vmihailenco/msgpack/v5"
)

type Handshake struct {
	Username  string         `msgpack:"username"`
	Color     string         `msgpack:"color"`
	ClientId  string         `msgpack:"client_id"`
	PublicKey *rsa.PublicKey `msgpack:"pubkey"`
}

func SendHandshake(conn *websocket.Conn, c config.Config, pubKey *rsa.PublicKey) error {
	var handshake Handshake
	handshake.Username = c.Username
	handshake.ClientId = generateUUID()
	handshake.PublicKey = pubKey
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
