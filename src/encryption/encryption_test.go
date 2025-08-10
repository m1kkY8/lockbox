package encryption

import (
	"crypto/rand"
	"testing"
)

func TestCreateRsaKey(t *testing.T) {
	keys, err := CreateRsaKey()
	if err != nil {
		t.Fatalf("CreateRsaKey() error = %v", err)
	}

	if keys == nil {
		t.Fatal("CreateRsaKey() returned nil keys")
	}

	if keys.PrivateKey == nil {
		t.Error("CreateRsaKey() returned nil private key")
	}

	if keys.PublicKey == nil {
		t.Error("CreateRsaKey() returned nil public key")
	}

	// Test key size
	if keys.PrivateKey.Size() != 256 { // 2048 bits = 256 bytes
		t.Errorf("CreateRsaKey() key size = %v, want 256", keys.PrivateKey.Size())
	}
}

func TestGenerateAES(t *testing.T) {
	key, err := GenerateAES()
	if err != nil {
		t.Fatalf("GenerateAES() error = %v", err)
	}

	if len(key.Key) != 32 {
		t.Errorf("GenerateAES() key length = %v, want 32", len(key.Key))
	}

	// Test that multiple calls generate different keys
	key2, err := GenerateAES()
	if err != nil {
		t.Fatalf("GenerateAES() second call error = %v", err)
	}

	if string(key.Key) == string(key2.Key) {
		t.Error("GenerateAES() should generate different keys")
	}
}

func TestEncryptDecryptMessage(t *testing.T) {
	// Generate AES key
	aesKey, err := GenerateAES()
	if err != nil {
		t.Fatalf("GenerateAES() error = %v", err)
	}

	plaintext := []byte("Hello, World! This is a test message.")

	// Encrypt
	ciphertext, err := EncryptMessage(plaintext, aesKey.Key)
	if err != nil {
		t.Fatalf("EncryptMessage() error = %v", err)
	}

	if len(ciphertext) == 0 {
		t.Error("EncryptMessage() returned empty ciphertext")
	}

	// Decrypt
	decrypted, err := DecryptMessage(ciphertext, aesKey.Key)
	if err != nil {
		t.Fatalf("DecryptMessage() error = %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("DecryptMessage() = %v, want %v", string(decrypted), string(plaintext))
	}
}

func TestEncryptDecryptAesKey(t *testing.T) {
	// Generate RSA keys
	rsaKeys, err := CreateRsaKey()
	if err != nil {
		t.Fatalf("CreateRsaKey() error = %v", err)
	}

	// Generate AES key
	aesKey, err := GenerateAES()
	if err != nil {
		t.Fatalf("GenerateAES() error = %v", err)
	}

	// Encrypt AES key with RSA public key
	encrypted, err := EncryptAesKey(aesKey, rsaKeys.PublicKey)
	if err != nil {
		t.Fatalf("EncryptAesKey() error = %v", err)
	}

	// Decrypt AES key with RSA private key
	decrypted, err := DecryptAesKey(encrypted, rsaKeys.PrivateKey)
	if err != nil {
		t.Fatalf("DecryptAesKey() error = %v", err)
	}

	if string(decrypted) != string(aesKey.Key) {
		t.Error("DecryptAesKey() did not return original AES key")
	}
}

func TestEncryptMessageWithInvalidKey(t *testing.T) {
	plaintext := []byte("test message")
	invalidKey := []byte("invalid") // Too short

	_, err := EncryptMessage(plaintext, invalidKey)
	if err == nil {
		t.Error("EncryptMessage() should fail with invalid key")
	}
}

func TestDecryptMessageWithInvalidKey(t *testing.T) {
	// Create a valid ciphertext first
	aesKey, _ := GenerateAES()
	plaintext := []byte("test message")
	ciphertext, _ := EncryptMessage(plaintext, aesKey.Key)

	// Try to decrypt with wrong key
	wrongKey := make([]byte, 32)
	rand.Read(wrongKey)

	_, err := DecryptMessage(ciphertext, wrongKey)
	if err == nil {
		t.Error("DecryptMessage() should fail with wrong key")
	}
}

// Benchmark encryption performance
func BenchmarkEncryptMessage(b *testing.B) {
	aesKey, _ := GenerateAES()
	plaintext := []byte("This is a benchmark test message for encryption performance.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := EncryptMessage(plaintext, aesKey.Key)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecryptMessage(b *testing.B) {
	aesKey, _ := GenerateAES()
	plaintext := []byte("This is a benchmark test message for decryption performance.")
	ciphertext, _ := EncryptMessage(plaintext, aesKey.Key)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DecryptMessage(ciphertext, aesKey.Key)
		if err != nil {
			b.Fatal(err)
		}
	}
}
