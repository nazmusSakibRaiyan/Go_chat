package chat

import (
	"encoding/json"
	"go-chat-backend/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	roomID   int
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	rooms      map[int]map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[int]map[*Client]bool),
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

			log.Printf("Client %s connected to room %d", client.username, client.roomID)

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

				log.Printf("Client %s disconnected from room %d", client.username, client.roomID)
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

func (h *Hub) broadcastToRoom(roomID int, message []byte) {
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
	roomIDStr := r.URL.Query().Get("room_id")

	if username == "" {
		username = "Anonymous"
	}

	roomID := 1 // Default room
	if roomIDStr != "" {
		if id, err := strconv.Atoi(roomIDStr); err == nil {
			roomID = id
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
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
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
			// Broadcast to room
			response := models.WebSocketMessage{
				Type:     "chat_message",
				Username: c.username,
				Content:  wsMsg.Content,
				RoomID:   c.roomID,
			}

			message, _ := json.Marshal(response)
			c.hub.broadcastToRoom(c.roomID, message)

			// TODO: Save to database when db connection is available
			log.Printf("Message from %s in room %d: %s", c.username, c.roomID, wsMsg.Content)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
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
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// REST API handlers
func GetRooms(c *gin.Context) {
	// For now, return mock data. In production, fetch from database
	rooms := []models.Room{
		{ID: 1, Name: "general", Description: "General chat room", CreatedAt: time.Now()},
		{ID: 2, Name: "random", Description: "Random discussions", CreatedAt: time.Now()},
		{ID: 3, Name: "tech", Description: "Technology discussions", CreatedAt: time.Now()},
	}

	c.JSON(http.StatusOK, rooms)
}

func CreateRoom(c *gin.Context) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set creation time
	room.CreatedAt = time.Now()

	// In production, save to database and return the created room
	room.ID = int(time.Now().Unix()) // Mock ID

	c.JSON(http.StatusCreated, room)
}

func GetRoomMessages(c *gin.Context) {
	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	// For now, return mock messages. In production, fetch from database
	messages := []models.Message{
		{
			ID:        1,
			RoomID:    roomID,
			Username:  "System",
			Content:   "Welcome to the chat room!",
			Type:      "system",
			CreatedAt: time.Now().Add(-1 * time.Hour),
		},
	}

	c.JSON(http.StatusOK, messages)
}
