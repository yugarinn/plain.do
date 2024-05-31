package repository

import (
	"math/rand"
	"time"
	"log"
)

func generateHash() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		hash := stringWithCharset(7, charset, seededRand)
		exists, err := hashExists(hash)

		if err != nil {
			log.Fatalf("Error checking hash existence: %v", err)
		}

		if !exists {
			return hash
		}
	}
}

func stringWithCharset(length int, charset string, seededRand *rand.Rand) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func hashExists(hash string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM todos_lists WHERE hash = ? LIMIT 1)"
	err := Database.QueryRow(query, hash).Scan(&exists)

	return exists, err
}
