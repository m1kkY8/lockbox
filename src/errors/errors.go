package errors

import (
	"errors"
	"fmt"
)

// Application error types
var (
	// Connection errors
	ErrConnectionFailed = errors.New("connection failed")
	ErrHandshakeFailed  = errors.New("handshake failed")
	ErrConnectionLost   = errors.New("connection lost")

	// Configuration errors
	ErrInvalidConfig = errors.New("invalid configuration")
	ErrMissingHost   = errors.New("missing host configuration")
	ErrMissingUser   = errors.New("missing username")

	// Encryption errors
	ErrKeyGeneration = errors.New("key generation failed")
	ErrEncryption    = errors.New("encryption failed")
	ErrDecryption    = errors.New("decryption failed")

	// Message errors
	ErrMessageEncode = errors.New("message encoding failed")
	ErrMessageDecode = errors.New("message decoding failed")
	ErrMessageSend   = errors.New("message send failed")
)

// ConnectionError represents connection-related errors
type ConnectionError struct {
	Host string
	Err  error
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("connection to %s failed: %v", e.Host, e.Err)
}

func (e *ConnectionError) Unwrap() error {
	return e.Err
}

// NewConnectionError creates a new connection error
func NewConnectionError(host string, err error) *ConnectionError {
	return &ConnectionError{
		Host: host,
		Err:  err,
	}
}

// ConfigError represents configuration-related errors
type ConfigError struct {
	Field string
	Value string
	Err   error
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("configuration error in field '%s' with value '%s': %v", e.Field, e.Value, e.Err)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}

// NewConfigError creates a new configuration error
func NewConfigError(field, value string, err error) *ConfigError {
	return &ConfigError{
		Field: field,
		Value: value,
		Err:   err,
	}
}

// EncryptionError represents encryption-related errors
type EncryptionError struct {
	Operation string
	Err       error
}

func (e *EncryptionError) Error() string {
	return fmt.Sprintf("encryption operation '%s' failed: %v", e.Operation, e.Err)
}

func (e *EncryptionError) Unwrap() error {
	return e.Err
}

// NewEncryptionError creates a new encryption error
func NewEncryptionError(operation string, err error) *EncryptionError {
	return &EncryptionError{
		Operation: operation,
		Err:       err,
	}
}
