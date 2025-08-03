package status

import (
	"go-chat-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StatusMiddleware provides middleware functions to check status-based permissions
type StatusMiddleware struct {
	statusService *StatusService
}

// NewStatusMiddleware creates a new status middleware
func NewStatusMiddleware(statusService *StatusService) *StatusMiddleware {
	return &StatusMiddleware{
		statusService: statusService,
	}
}

// RequireAction creates middleware that checks if user can perform a specific action
func (sm *StatusMiddleware) RequireAction(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		user, ok := userInterface.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Invalid user data",
			})
			c.Abort()
			return
		}

		err := sm.statusService.ValidateStatusForAction(user.ID.Hex(), action)
		if err != nil {
			if IsStatusActionError(err) {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"message": err.Error(),
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Failed to validate action",
				})
			}
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireCanJoinRooms middleware checks if user can join rooms
func (sm *StatusMiddleware) RequireCanJoinRooms() gin.HandlerFunc {
	return sm.RequireAction("join_rooms")
}

// RequireCanSendMessages middleware checks if user can send messages
func (sm *StatusMiddleware) RequireCanSendMessages() gin.HandlerFunc {
	return sm.RequireAction("send_messages")
}

// InjectStatusInfo middleware adds status information to the context
func (sm *StatusMiddleware) InjectStatusInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if exists {
			user, ok := userInterface.(*models.User)
			if ok {
				statusInfo, err := sm.statusService.GetUserStatusInfo(user.ID.Hex())
				if err == nil {
					c.Set("status_info", statusInfo)
					c.Set("user_status", statusInfo.Status)
				}
			}
		}
		c.Next()
	}
}

// RecordActivity middleware updates user activity (bringing auto-away users back online)
func (sm *StatusMiddleware) RecordActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only record activity on successful requests
		c.Next()

		// Update activity after request completion
		if c.Writer.Status() < 400 { // Only on successful requests
			userInterface, exists := c.Get("user")
			if exists {
				user, ok := userInterface.(*models.User)
				if ok {
					// Update user activity (this will bring auto-away users back online)
					sm.statusService.UpdateUserActivity(user.ID.Hex())
				}
			}
		}
	}
}

// StatusAwareResponse middleware modifies responses based on user status
func (sm *StatusMiddleware) StatusAwareResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Add status information to certain responses
		if c.Writer.Header().Get("Content-Type") == "application/json; charset=utf-8" {
			userInterface, exists := c.Get("user")
			if exists {
				user, ok := userInterface.(*models.User)
				if ok {
					statusInfo, err := sm.statusService.GetUserStatusInfo(user.ID.Hex())
					if err == nil {
						// Add status capabilities to response headers for client-side logic
						c.Header("X-User-Status", string(statusInfo.Status))
						c.Header("X-Can-Send-Messages", boolToString(statusInfo.CanSendMessages))
						c.Header("X-Can-Join-Rooms", boolToString(statusInfo.CanJoinRooms))
						c.Header("X-Can-Receive-Notifications", boolToString(statusInfo.CanReceiveNotifications))
					}
				}
			}
		}
	}
}

// Helper function to convert bool to string
func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
