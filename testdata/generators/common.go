package generators

import (
	"crypto/rand"
	"github.com/google/uuid"
	"math/big"
)

func RandomUuidString() string {
	u := uuid.New()
	return u.String()
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString generates a random string of the specified length.
func RandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return ""
		}
		result[i] = charset[randomIndex.Int64()]
	}
	return string(result)
}
