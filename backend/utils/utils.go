package utils

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
	"unicode"
)

// GenerateID generates a random hexadecimal ID
func GenerateID(length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// SanitizeUsername removes invalid characters from username
func SanitizeUsername(username string) string {
	// Remove leading/trailing whitespace
	username = strings.TrimSpace(username)

	// Replace multiple spaces with single space
	re := regexp.MustCompile(`\s+`)
	username = re.ReplaceAllString(username, " ")

	// Remove any non-printable characters
	result := strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, username)

	return result
}

// ValidateRoomName checks if room name is valid
func ValidateRoomName(name string) bool {
	name = strings.TrimSpace(name)
	if len(name) < 1 || len(name) > 50 {
		return false
	}

	// Only allow alphanumeric, spaces, hyphens, and underscores
	re := regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)
	return re.MatchString(name)
}

// TruncateText truncates text to specified length with ellipsis
func TruncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}

	if maxLength <= 3 {
		return text[:maxLength]
	}

	return text[:maxLength-3] + "..."
}
