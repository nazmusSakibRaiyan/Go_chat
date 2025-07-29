package main

import (
	"database/sql"
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

	// Initialize database (optional for now)
	var database *sql.DB
	var err error
	if cfg.DatabaseURL != "" {
		database, err = db.Initialize(cfg.DatabaseURL)
		if err != nil {
			log.Printf("Warning: Failed to initialize database: %v", err)
			log.Println("Running without database support...")
		} else {
			defer database.Close()
		}
	}

	// Create chat hub
	hub := chat.NewHub()
	go hub.Run()

	// Setup router
	router := gin.Default()

	// CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
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
		api.GET("/rooms", chat.GetRooms)
		api.POST("/rooms", chat.CreateRoom)
		api.GET("/rooms/:id/messages", chat.GetRoomMessages)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
