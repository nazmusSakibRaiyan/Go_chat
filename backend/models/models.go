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
