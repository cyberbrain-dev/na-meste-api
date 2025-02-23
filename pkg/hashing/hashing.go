// Contains tools for hashing the data
package hashing

import (
	"crypto/sha256"
	"fmt"
)

// Computes and returns hashed string using SHA-256 algorithm
func HashSHA256(input string) string {
	hash := sha256.Sum256([]byte(input)) // Compute SHA-256 hash
	return fmt.Sprintf("%x", hash)       // Convert hash to hexadecimal string
}
