package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and sets user in context
func AuthMiddleware(authHandlers *AuthHandlers) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check Bearer token format
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate JWT token
		claims, err := authHandlers.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Get user from database
		if authHandlers.mongoDB != nil {
			user, err := authHandlers.mongoDB.GetUserByID(claims.UserID)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "User not found",
				})
				c.Abort()
				return
			}

			// Set user in context
			c.Set("user", user)
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
		}

		c.Next()
	}
}

// OptionalAuthMiddleware checks for authentication but doesn't require it
func OptionalAuthMiddleware(authHandlers *AuthHandlers) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		token := tokenParts[1]
		claims, err := authHandlers.ValidateJWT(token)
		if err != nil {
			c.Next()
			return
		}

		if authHandlers.mongoDB != nil {
			user, err := authHandlers.mongoDB.GetUserByID(claims.UserID)
			if err == nil {
				c.Set("user", user)
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
			}
		}

		c.Next()
	}
}
