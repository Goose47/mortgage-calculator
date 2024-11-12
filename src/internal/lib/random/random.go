// Package random provides functions to generate cryptographically secure random entities.
package random

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// String generates a cryptographically secure random string of length n.
func String(n int) (string, error) {
	byteStr := make([]byte, n)
	for i := 0; i < n; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random string: %w", err)
		}
		byteStr[i] = charset[index.Int64()]
	}
	return string(byteStr), nil
}
