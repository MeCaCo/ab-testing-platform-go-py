package service

import (
	"crypto/sha256"
	"fmt"
)

// AssignVariant determines which group (A or B) a user belongs to
// Uses deterministic hashing: same userID + testID always gets the same variant
func AssignVariant(testID, userID string) string {
	// Combine test_id and user_id to make assignment unique per test
	input := fmt.Sprintf("%s:%s", testID, userID)

	hash := sha256.Sum256([]byte(input))

	// Take first 4 bytes and convert to int
	// This ensures even distribution (50/50 split)
	num := int(hash[0])<<24 | int(hash[1])<<16 | int(hash[2])<<8 | int(hash[3])

	if num < 0 {
		num = -num // Ensure positive number
	}

	// 50% chance for each variant
	if num%100 < 50 {
		return "A"
	}
	return "B"
}
