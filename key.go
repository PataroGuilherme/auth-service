package main

import (
	"crypto/rand"
	"encoding/hex"
)

type KeyInfo struct {
	Key string `json:"key"`
}

func generateKey() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
