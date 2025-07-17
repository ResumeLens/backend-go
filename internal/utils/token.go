package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomToken(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
