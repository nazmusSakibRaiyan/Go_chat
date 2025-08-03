package chat

import (
	"encoding/json"
	"go-chat-backend/db"
	"go-chat-backend/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
	roomID   string
	userID   string // Add user ID for authenticated users
	isAuth   bool   // Track if user is authenticated
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	rooms      map[string]map[*Client]bool
	mongoDB    *db.MongoDB
}

func NewHub(mongoDB *db.MongoDB) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		mongoDB:    mongoDB,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			if h.rooms[client.roomID] == nil {
				h.rooms[client.roomID] = make(map[*Client]bool)
			}
			h.rooms[client.roomID][client] = true

			// Send welcome message
			welcomeMsg := models.WebSocketMessage{
				Type:     "user_joined",
				Username: client.username,
				Content:  client.username + " joined the room",
			}
			message, _ := json.Marshal(welcomeMsg)
			h.broadcastToRoom(client.roomID, message)

			log.Printf("Client %s connected to room %s", client.username, client.roomID)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.rooms[client.roomID], client)
				close(client.send)

				// Send leave message
				leaveMsg := models.WebSocketMessage{
					Type:     "user_left",
					Username: client.username,
					Content:  client.username + " left the room",
				}
				message, _ := json.Marshal(leaveMsg)
				h.broadcastToRoom(client.roomID, message)

				log.Printf("Client %s disconnected from room %s", client.username, client.roomID)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) broadcastToRoom(roomID string, message []byte) {
	if room, exists := h.rooms[roomID]; exists {
		for client := range room {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
				delete(room, client)
			}
		}
	}
}

func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// Get username and room from query parameters
	username := r.URL.Query().Get("username")
	roomID := r.URL.Query().Get("room_id")

	if username == "" {
		username = "Anonymous"
	}

	// Default to first room if not provided
	if roomID == "" {
		// Try to get the first available room
		if hub.mongoDB != nil {
			rooms, err := hub.mongoDB.GetRooms()
			if err == nil && len(rooms) > 0 {
				roomID = rooms[0].ID.Hex()
			}
		}
		// Fallback to a default ObjectID-like string
		if roomID == "" {
			roomID = "general"
		}
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
		roomID:   roomID,
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(120 * time.Second)) // Increased timeout
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(120 * time.Second)) // Increased timeout
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var wsMsg models.WebSocketMessage
		if err := json.Unmarshal(messageBytes, &wsMsg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Handle different message types
		switch wsMsg.Type {
		case "chat_message":
			// Save message to database if MongoDB is available
			if c.hub.mongoDB != nil {
				roomObjectID, err := primitive.ObjectIDFromHex(c.roomID)
				if err != nil {
					log.Printf("Invalid room ID: %v", err)
					continue
				}

				msg := models.Message{
					RoomID:    roomObjectID,
					Username:  c.username,
					Content:   wsMsg.Content,
					Type:      "text",
					CreatedAt: time.Now(),
				}

				_, err = c.hub.mongoDB.SaveMessage(msg)
				if err != nil {
					log.Printf("Error saving message: %v", err)
				}
			}

			// Broadcast to room
			response := models.WebSocketMessage{
				Type:     "chat_message",
				Username: c.username,
				Content:  wsMsg.Content,
				RoomID:   c.roomID,
			}

			message, _ := json.Marshal(response)
			c.hub.broadcastToRoom(c.roomID, message)

			log.Printf("Message from %s in room %s: %s", c.username, c.roomID, wsMsg.Content)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(110 * time.Second) // Increased to 110 seconds
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(20 * time.Second)) // Increased timeout
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(20 * time.Second)) // Increased timeout
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// REST API handlers
func GetRooms(c *gin.Context, mongoDB *db.MongoDB) {
	if mongoDB == nil {
		// Return mock data if no database
		rooms := []map[string]interface{}{
			{"id": "general", "name": "general", "description": "General chat room", "created_at": time.Now()},
			{"id": "random", "name": "random", "description": "Random discussions", "created_at": time.Now()},
			{"id": "tech", "name": "tech", "description": "Technology discussions", "created_at": time.Now()},
		}
		c.JSON(http.StatusOK, rooms)
		return
	}

	rooms, err := mongoDB.GetRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms"})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func CreateRoom(c *gin.Context, mongoDB *db.MongoDB) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if mongoDB == nil {
		// Return mock data if no database
		room.ID = primitive.NewObjectID()
		room.CreatedAt = time.Now()
		room.UpdatedAt = time.Now()
		c.JSON(http.StatusCreated, room)
		return
	}

	createdRoom, err := mongoDB.CreateRoom(room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}

	c.JSON(http.StatusCreated, createdRoom)
}

func GetRoomMessages(c *gin.Context, mongoDB *db.MongoDB) {
	roomID := c.Param("id")

	if mongoDB == nil {
		// Return mock messages if no database
		messages := []map[string]interface{}{
			{
				"id":         primitive.NewObjectID().Hex(),
				"room_id":    roomID,
				"username":   "System",
				"content":    "Welcome to the chat room!",
				"type":       "system",
				"created_at": time.Now().Add(-1 * time.Hour),
			},
		}
		c.JSON(http.StatusOK, messages)
		return
	}

	messages, err := mongoDB.GetRoomMessages(roomID, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
