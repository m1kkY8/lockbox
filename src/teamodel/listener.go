package teamodel

import (
	"context"
	"log"
	"strings"

	"github.com/m1kkY8/lockbox/src/encryption"
	"github.com/m1kkY8/lockbox/src/message"
	"github.com/vmihailenco/msgpack/v5"
)

// receiveMessages handles incoming WebSocket messages
// This method replaces the original recieveMessages (fixing the typo)
func (m *Model) receiveMessages(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in receiveMessages: %v", r)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Message receiving stopped due to context cancellation")
			return
		default:
			// Continue with message processing
		}

		// Read messages from socket with timeout/cancellation support
		messageType, byteMessage, err := m.client.Conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			return
		}

		_ = messageType // Acknowledge that we're not using this right now

		// Decode as regular message
		decodedMsg, err := message.DecodeMessage(byteMessage)
		if err != nil {
			log.Printf("Failed decoding message: %v", err)
			continue
		}

		if err := m.handleMessage(decodedMsg, byteMessage); err != nil {
			log.Printf("Failed to handle message: %v", err)
			continue
		}
	}
}

// handleMessage processes different types of incoming messages
func (m *Model) handleMessage(decodedMsg message.Message, byteMessage []byte) error {
	switch decodedMsg.Type {
	case message.ServerMessage:
		return m.handleServerMessage(decodedMsg)
	case message.KeyMessage:
		return m.handleKeyMessage(byteMessage)
	case message.ChatMessage:
		return m.handleChatMessage(decodedMsg)
	default:
		log.Printf("Unknown message type: %d", decodedMsg.Type)
		return nil
	}
}

// handleServerMessage processes server messages (online users)
func (m *Model) handleServerMessage(msg message.Message) error {
	users := strings.Split(msg.Content, " ")
	select {
	case m.channels.OnlineUsers <- users:
		return nil
	default:
		log.Println("Online users channel is full, dropping message")
		return nil
	}
}

// handleKeyMessage processes public key messages
func (m *Model) handleKeyMessage(byteMessage []byte) error {
	var keyMsg message.PublicKeys
	if err := msgpack.Unmarshal(byteMessage, &keyMsg); err != nil {
		return err
	}

	m.client.PublicKeys = keyMsg.PublicKeys
	return nil
}

// handleChatMessage processes encrypted chat messages
func (m *Model) handleChatMessage(decodedMsg message.Message) error {
	// Try to decrypt the message with available keys
	for _, key := range decodedMsg.AESKeys {
		aesKey, err := encryption.DecryptAesKey(key.Key, m.client.KeyPair.PrivateKey)
		if err != nil {
			continue // Try next key
		}

		// Decrypt message content with AES key
		decryptedContent, err := encryption.DecryptMessage([]byte(decodedMsg.Content), aesKey)
		if err != nil {
			continue // Try next key
		}

		// Successfully decrypted
		decodedMsg.Content = string(decryptedContent)
		formattedMessage := message.Format(decodedMsg)

		select {
		case m.channels.Messages <- formattedMessage:
			return nil
		default:
			log.Println("Message channel is full, dropping message")
			return nil
		}
	}

	log.Println("Failed to decrypt message with any available key")
	return nil
}

// Legacy method for backward compatibility
// TODO: Remove this once all references are updated
func (m *Model) recieveMessages() {
	ctx := context.Background()
	m.receiveMessages(ctx)
}
