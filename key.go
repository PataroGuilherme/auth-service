package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// generateAPIKey gera uma chave de API segura e aleatória (32 bytes -> 64 chars hex)
func generateAPIKey() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// hashAPIKey cria o hash SHA-256 da chave fornecida
func hashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
