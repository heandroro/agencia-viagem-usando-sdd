package crypto

import (
	"testing"
)

func TestAES256Service_EncryptDecrypt(t *testing.T) {
	service, err := NewAES256Service("test-key-that-is-32-bytes-long!")
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	tests := []struct {
		name      string
		plaintext string
	}{
		{"CPF", "12345678901"},
		{"Passport", "AB123456"},
		{"Empty", ""},
		{"Long text", "this is a very long text that should be encrypted and decrypted correctly"},
		{"Special chars", "!@#$%^&*()_+-=[]{}|;':\",./<>?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := service.Encrypt(tt.plaintext)
			if err != nil {
				t.Errorf("Encrypt() error = %v", err)
				return
			}

			// Ciphertext should be different from plaintext
			if encrypted == tt.plaintext && tt.plaintext != "" {
				t.Error("Encrypted text should be different from plaintext")
			}

			decrypted, err := service.Decrypt(encrypted)
			if err != nil {
				t.Errorf("Decrypt() error = %v", err)
				return
			}

			if decrypted != tt.plaintext {
				t.Errorf("Decrypt() = %v, want %v", decrypted, tt.plaintext)
			}
		})
	}
}

func TestAES256Service_Hash(t *testing.T) {
	service, _ := NewAES256Service("test-key-that-is-32-bytes-long!")

	hash1 := service.Hash("12345678901")
	hash2 := service.Hash("12345678901")
	hash3 := service.Hash("different-value")

	// Same input should produce same hash
	if hash1 != hash2 {
		t.Error("Same input should produce same hash")
	}

	// Different input should produce different hash
	if hash1 == hash3 {
		t.Error("Different input should produce different hash")
	}

	// Hash should be 64 characters (SHA-256 hex)
	if len(hash1) != 64 {
		t.Errorf("Hash length = %d, want 64", len(hash1))
	}
}

func TestAES256Service_Mask(t *testing.T) {
	service, _ := NewAES256Service("test-key-that-is-32-bytes-long!")

	tests := []struct {
		input    string
		expected string
	}{
		{"12345678901", "***78901"},
		{"AB123456", "***3456"},
		{"1234", "****"},
		{"12", "****"},
		{"", "****"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := service.Mask(tt.input)
			if result != tt.expected {
				t.Errorf("Mask(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAES256Service_GenerateDocumentHash(t *testing.T) {
	service, _ := NewAES256Service("test-key-that-is-32-bytes-long!")

	hash := service.GenerateDocumentHash("12345678901")

	// Should be 16 characters
	if len(hash) != 16 {
		t.Errorf("DocumentHash length = %d, want 16", len(hash))
	}

	// Same document should produce same hash
	hash2 := service.GenerateDocumentHash("12345678901")
	if hash != hash2 {
		t.Error("Same document should produce same partial hash")
	}
}

func TestNewAES256Service_InvalidKey(t *testing.T) {
	_, err := NewAES256Service("short")
	if err == nil {
		t.Error("Should fail with short key")
	}
}
