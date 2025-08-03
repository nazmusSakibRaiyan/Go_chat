package auth

import (
	"go-chat-backend/db"
	"go-chat-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandlers struct {
	mongoDB     *db.MongoDB
	jwtSecret   string
	rateLimiter *RateLimiter
}

func NewAuthHandlers(mongoDB *db.MongoDB, jwtSecret string) *AuthHandlers {
	// Create rate limiter: 10 attempts per 5 minutes (more lenient for development)
	rateLimiter := NewRateLimiter(10, 5*time.Minute)

	// Start cleanup routine
	go func() {
		ticker := time.NewTicker(time.Minute) // Clean more frequently
		defer ticker.Stop()
		for range ticker.C {
			rateLimiter.Clean()
		}
	}()

	return &AuthHandlers{
		mongoDB:     mongoDB,
		jwtSecret:   jwtSecret,
		rateLimiter: rateLimiter,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	DisplayName string `json:"display_name" binding:"required,min=1,max=50"`
	Avatar      *int   `json:"avatar,omitempty"` // Optional avatar ID (1-12)
}

type AuthResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Token   string       `json:"token,omitempty"`
	User    *models.User `json:"user,omitempty"`
}

// JWT Claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Register handles user registration
func (h *AuthHandlers) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	// Check if database is available
	if h.mongoDB == nil {
		c.JSON(http.StatusServiceUnavailable, AuthResponse{
			Success: false,
			Message: "Database not available",
		})
		return
	}

	// Validate username
	if !ValidateUsername(req.Username) {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Message: "Username must be 3-20 characters long and contain only valid characters",
		})
		return
	}

	// Validate password strength
	if !ValidatePassword(req.Password) {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Message: "Password must be at least 6 characters long and contain a mix of letters and numbers",
		})
		return
	}

	// Check if user already exists
	existingUser, _ := h.mongoDB.GetUserByUsername(req.Username)
	if existingUser != nil {
		c.JSON(http.StatusConflict, AuthResponse{
			Success: false,
			Message: "Username already exists",
		})
		return
	}

	// Check if email already exists
	existingEmail, _ := h.mongoDB.GetUserByEmail(req.Email)
	if existingEmail != nil {
		c.JSON(http.StatusConflict, AuthResponse{
			Success: false,
			Message: "Email already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Message: "Failed to process password",
		})
		return
	}

	// Create user
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, err := h.mongoDB.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Message: "Failed to create user",
		})
		return
	}

	// Generate JWT token
	token, err := h.generateJWT(createdUser.ID.Hex(), createdUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Message: "Failed to generate token",
		})
		return
	}

	// Remove password hash from response
	createdUser.PasswordHash = ""

	c.JSON(http.StatusCreated, AuthResponse{
		Success: true,
		Message: "User registered successfully",
		Token:   token,
		User:    createdUser,
	})
}

// Login handles user authentication
func (h *AuthHandlers) Login(c *gin.Context) {
	// Check rate limiting - Temporarily disabled for development
	// clientIP := c.ClientIP()
	// if !h.rateLimiter.IsAllowed(clientIP) {
	// 	c.JSON(http.StatusTooManyRequests, AuthResponse{
	// 		Success: false,
	// 		Message: "Too many login attempts. Please try again later.",
	// 	})
	// 	return
	// }

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	// Check if database is available
	if h.mongoDB == nil {
		c.JSON(http.StatusServiceUnavailable, AuthResponse{
			Success: false,
			Message: "Database not available",
		})
		return
	}

	// Get user by username
	user, err := h.mongoDB.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, AuthResponse{
			Success: false,
			Message: "Invalid username or password",
		})
		return
	}

	// Check password
	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, AuthResponse{
			Success: false,
			Message: "Invalid username or password",
		})
		return
	}

	// Generate JWT token
	token, err := h.generateJWT(user.ID.Hex(), user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Message: "Failed to generate token",
		})
		return
	}

	// Remove password hash from response
	user.PasswordHash = ""

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		User:    user,
	})
}

// Logout handles user logout (mainly for client-side token removal)
func (h *AuthHandlers) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Message: "Logout successful",
	})
}

// Me returns current user information
func (h *AuthHandlers) Me(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, AuthResponse{
			Success: false,
			Message: "User not found in context",
		})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Message: "Invalid user data",
		})
		return
	}

	// Remove password hash from response
	user.PasswordHash = ""

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Message: "User data retrieved",
		User:    user,
	})
}

// generateJWT creates a new JWT token
func (h *AuthHandlers) generateJWT(userID, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

// ValidateJWT validates a JWT token and returns claims
func (h *AuthHandlers) ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

// UpdateProfile allows users to update their display name and avatar
func (h *AuthHandlers) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	// Validate avatar if provided
	if req.Avatar != nil && !models.IsValidAvatar(*req.Avatar) {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Message: "Invalid avatar ID. Must be between 1 and 12",
		})
		return
	}

	// Get user from context (set by AuthMiddleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, AuthResponse{
			Success: false,
			Message: "User not found in context",
		})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Message: "Invalid user data in context",
		})
		return
	}

	// Check if database is available
	if h.mongoDB == nil {
		c.JSON(http.StatusServiceUnavailable, AuthResponse{
			Success: false,
			Message: "Database not available",
		})
		return
	}

	// Update user's display name and avatar
	err := h.mongoDB.UpdateUserProfile(user.ID.Hex(), req.DisplayName, req.Avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Message: "Failed to update profile: " + err.Error(),
		})
		return
	}

	// Get updated user data
	updatedUser, err := h.mongoDB.GetUserByID(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Message: "Failed to fetch updated user data",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Message: "Profile updated successfully",
		User:    updatedUser,
	})
}

// GetAvatars returns the list of available avatars
func (h *AuthHandlers) GetAvatars(c *gin.Context) {
	avatars := models.GetAvailableAvatars()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Available avatars retrieved",
		"avatars": avatars,
	})
}
