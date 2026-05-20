package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// Service interface para operações criptográficas
type Service interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
	Hash(data string) string
	Mask(data string) string
}

// AES256Service implementação AES-256-GCM
type AES256Service struct {
	key []byte
}

// NewAES256Service cria um novo serviço de criptografia
func NewAES256Service(key string) (*AES256Service, error) {
	if len(key) < 32 {
		return nil, errors.New("encryption key must be at least 32 characters")
	}
	
	// Derivar chave de 32 bytes usando SHA-256
	hash := sha256.Sum256([]byte(key))
	
	return &AES256Service{
		key: hash[:],
	}, nil
}

// Encrypt criptografa texto usando AES-256-GCM
func (s *AES256Service) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt descriptografa texto usando AES-256-GCM
func (s *AES256Service) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	
	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}
	
	return string(plaintext), nil
}

// Hash gera hash SHA-256 de um dado
func (s *AES256Service) Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Mask mascara um dado sensível para exibição
func (s *AES256Service) Mask(data string) string {
	if len(data) <= 4 {
		return "****"
	}
	return fmt.Sprintf("***%s", data[len(data)-4:])
}

// GenerateDocumentHash gera um hash parcial para correlacionar em logs
func (s *AES256Service) GenerateDocumentHash(document string) string {
	// Usar apenas parte do hash para não expor o documento completo
	fullHash := s.Hash(document)
	return fullHash[:16] // Primeiros 16 caracteres do hash
}
