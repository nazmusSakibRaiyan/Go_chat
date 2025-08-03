package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"password_hash"`
	DisplayName  string             `json:"display_name" bson:"display_name"`
	Avatar       int                `json:"avatar" bson:"avatar"` // Avatar ID (1-12)
	Status       string             `json:"status" bson:"status"` // User status: online, away, busy
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

// GetDisplayName returns display name or username as fallback
func (u *User) GetDisplayName() string {
	if u.DisplayName != "" {
		return u.DisplayName
	}
	return u.Username
}

type Message struct {
	ID        primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	RoomID    primitive.ObjectID  `json:"room_id" bson:"room_id"`
	UserID    *primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Username  string              `json:"username" bson:"username"`
	Content   string              `json:"content" bson:"content"`
	Type      string              `json:"type" bson:"type"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
}

type WebSocketMessage struct {
	Type     string      `json:"type"`
	RoomID   string      `json:"room_id,omitempty"`
	Username string      `json:"username,omitempty"`
	Content  string      `json:"content,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

type Avatar struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// GetAvailableAvatars returns the list of 12 available avatars
func GetAvailableAvatars() []Avatar {
	return []Avatar{
		{ID: 1, Name: "Cat", URL: "/static/avatars/cat.png"},
		{ID: 2, Name: "Dog", URL: "/static/avatars/dog.png"},
		{ID: 3, Name: "Bear", URL: "/static/avatars/bear.png"},
		{ID: 4, Name: "Fox", URL: "/static/avatars/fox.png"},
		{ID: 5, Name: "Lion", URL: "/static/avatars/lion.png"},
		{ID: 6, Name: "Panda", URL: "/static/avatars/panda.png"},
		{ID: 7, Name: "Robot", URL: "/static/avatars/robot.png"},
		{ID: 8, Name: "Alien", URL: "/static/avatars/alien.png"},
		{ID: 9, Name: "Ninja", URL: "/static/avatars/ninja.png"},
		{ID: 10, Name: "Pirate", URL: "/static/avatars/pirate.png"},
		{ID: 11, Name: "Knight", URL: "/static/avatars/knight.png"},
		{ID: 12, Name: "Wizard", URL: "/static/avatars/wizard.png"},
	}
}

// IsValidAvatar checks if the avatar ID is valid (1-12)
func IsValidAvatar(avatarID int) bool {
	return avatarID >= 1 && avatarID <= 12
}

// GetAvatarURL returns the URL for a given avatar ID
func GetAvatarURL(avatarID int) string {
	if !IsValidAvatar(avatarID) {
		return "/static/avatars/default.png"
	}
	avatars := GetAvailableAvatars()
	return avatars[avatarID-1].URL
}

// User status constants
const (
	StatusOnline = "online"
	StatusAway   = "away"
	StatusBusy   = "busy"
)

// GetValidStatuses returns the list of valid user statuses
func GetValidStatuses() []string {
	return []string{StatusOnline, StatusAway, StatusBusy}
}

// IsValidStatus checks if the status is valid
func IsValidStatus(status string) bool {
	validStatuses := GetValidStatuses()
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// GetDefaultStatus returns the default status for new users
func GetDefaultStatus() string {
	return StatusOnline
}
