package auth

import (
	"go-chat-backend/utils"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidateUsername checks if username meets requirements
func ValidateUsername(username string) bool {
	sanitized := utils.SanitizeUsername(username)
	return len(sanitized) >= 3 && len(sanitized) <= 20
}

// ValidatePassword checks if password meets security requirements
func ValidatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	// Check for at least one letter and one number
	hasLetter := false
	hasNumber := false

	for _, char := range password {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			hasLetter = true
		}
		if char >= '0' && char <= '9' {
			hasNumber = true
		}
		if hasLetter && hasNumber {
			break
		}
	}

	return hasLetter && hasNumber
}

// GenerateSessionID creates a random session ID
func GenerateSessionID() string {
	return utils.GenerateID(32)
}
