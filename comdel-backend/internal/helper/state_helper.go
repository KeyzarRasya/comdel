package helper

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandState() string {
	bytes := make([]byte, 32);

	_, err := rand.Read(bytes);
	if err != nil {
		panic("Failed to create random state");
	}

	return hex.EncodeToString(bytes);
}