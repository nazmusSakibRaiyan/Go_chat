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

// GenerateSessionID creates a random session ID
func GenerateSessionID() string {
	return utils.GenerateID(32)
}
