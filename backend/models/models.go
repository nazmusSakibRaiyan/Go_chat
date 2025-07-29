package models

import "time"

type Room struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Message struct {
	ID        int       `json:"id"`
	RoomID    int       `json:"room_id"`
	UserID    *int      `json:"user_id,omitempty"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type WebSocketMessage struct {
	Type     string      `json:"type"`
	RoomID   int         `json:"room_id,omitempty"`
	Username string      `json:"username,omitempty"`
	Content  string      `json:"content,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}
