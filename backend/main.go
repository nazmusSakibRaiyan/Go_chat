package main

import (
	"go-chat-backend/auth"
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

	// Create auth handlers
	authHandlers := auth.NewAuthHandlers(mongoDB, cfg.JWTSecret)

	// Setup router
	router := gin.Default()

	// CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // Allow all origins for development (including file://)
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(corsConfig))

	// Serve static files
	router.Static("/static", "../frontend/public/static")

	// API routes
	api := router.Group("/api")
	{
		// Authentication routes (no auth required)
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandlers.Register)
			authRoutes.POST("/login", authHandlers.Login)
			authRoutes.POST("/logout", authHandlers.Logout)
		}

		// Protected routes (auth required)
		protected := api.Group("/")
		protected.Use(auth.AuthMiddleware(authHandlers))
		{
			protected.GET("/me", authHandlers.Me)
			protected.PUT("/profile", authHandlers.UpdateProfile)
		}

		// Public routes (no auth required)
		api.GET("/avatars", authHandlers.GetAvatars)

		// WebSocket endpoint (optional auth)
		api.GET("/ws", func(c *gin.Context) {
			chat.HandleWebSocket(hub, c.Writer, c.Request)
		})

		// REST API endpoints (no auth required for now, but can be protected later)
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
