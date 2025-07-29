package db

import (
	"database/sql"
	"go-chat-backend/models"

	_ "github.com/lib/pq"
)

func Initialize(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS rooms (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			room_id INTEGER REFERENCES rooms(id),
			user_id INTEGER REFERENCES users(id),
			username VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			message_type VARCHAR(50) DEFAULT 'text',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	// Insert default room if not exists
	_, err := db.Exec(`
		INSERT INTO rooms (name, description) 
		VALUES ('general', 'General chat room for everyone') 
		ON CONFLICT (name) DO NOTHING
	`)

	return err
}

func GetRooms(db *sql.DB) ([]models.Room, error) {
	rows, err := db.Query("SELECT id, name, description, created_at FROM rooms ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.CreatedAt)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func CreateRoom(db *sql.DB, room models.Room) (*models.Room, error) {
	var newRoom models.Room
	err := db.QueryRow(
		"INSERT INTO rooms (name, description) VALUES ($1, $2) RETURNING id, name, description, created_at",
		room.Name, room.Description,
	).Scan(&newRoom.ID, &newRoom.Name, &newRoom.Description, &newRoom.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &newRoom, nil
}

func GetRoomMessages(db *sql.DB, roomID int, limit int) ([]models.Message, error) {
	query := `
		SELECT id, room_id, user_id, username, content, message_type, created_at 
		FROM messages 
		WHERE room_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2
	`

	rows, err := db.Query(query, roomID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.ID, &msg.RoomID, &msg.UserID, &msg.Username, &msg.Content, &msg.Type, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	// Reverse the slice to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func SaveMessage(db *sql.DB, msg models.Message) error {
	_, err := db.Exec(
		"INSERT INTO messages (room_id, user_id, username, content, message_type) VALUES ($1, $2, $3, $4, $5)",
		msg.RoomID, msg.UserID, msg.Username, msg.Content, msg.Type,
	)
	return err
}
