package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result) + "@"
}

func GenerateHashFromString(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	emailHash := hex.EncodeToString(h.Sum(nil))
	truncatedHash := emailHash[:8]

	return truncatedHash
}
