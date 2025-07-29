package main

import (
	"go-chat-backend/chat"
	"go-chat-backend/config"
	"go-chat-backend/db"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize MongoDB
	var mongoDB *db.MongoDB
	var err error
	if cfg.MongoURI != "" {
		mongoDB, err = db.Initialize(cfg.MongoURI, cfg.MongoDBName)
		if err != nil {
			log.Printf("Warning: Failed to initialize MongoDB: %v", err)
			log.Println("Running without database support...")
		} else {
			defer mongoDB.Close()
			log.Println("Connected to MongoDB successfully!")
		}
	}

	// Create chat hub
	hub := chat.NewHub(mongoDB)
	go hub.Run()

	// Setup router
	router := gin.Default()

	// CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173", "http://127.0.0.1:5500"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	// API routes
	api := router.Group("/api")
	{
		// WebSocket endpoint
		api.GET("/ws", func(c *gin.Context) {
			chat.HandleWebSocket(hub, c.Writer, c.Request)
		})

		// REST API endpoints
		api.GET("/rooms", func(c *gin.Context) {
			chat.GetRooms(c, mongoDB)
		})
		api.POST("/rooms", func(c *gin.Context) {
			chat.CreateRoom(c, mongoDB)
		})
		api.GET("/rooms/:id/messages", func(c *gin.Context) {
			chat.GetRoomMessages(c, mongoDB)
		})
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
