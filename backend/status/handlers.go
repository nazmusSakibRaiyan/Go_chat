package status

import (
	"go-chat-backend/db"
	"go-chat-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StatusHandlers provides HTTP handlers for status-related operations
type StatusHandlers struct {
	statusService *StatusService
	mongoDB       *db.MongoDB
}

// NewStatusHandlers creates new status handlers
func NewStatusHandlers(mongoDB *db.MongoDB) *StatusHandlers {
	return &StatusHandlers{
		statusService: NewStatusService(mongoDB),
		mongoDB:       mongoDB,
	}
}

// GetStatusService returns the status service for use by other handlers
func (sh *StatusHandlers) GetStatusService() *StatusService {
	return sh.statusService
}

// GetStatuses returns all available statuses with their capabilities
func (sh *StatusHandlers) GetStatuses(c *gin.Context) {
	statuses := sh.statusService.GetStatusManager().GetAvailableStatusesForAPI()

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Available statuses retrieved",
		"statuses": statuses,
	})
}

// UpdateUserStatus allows authenticated users to update their status
func (sh *StatusHandlers) UpdateUserStatus(c *gin.Context) {
	// Get user from context (set by AuthMiddleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not found in context",
		})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user data in context",
		})
		return
	}

	var req struct {
		Status        string `json:"status" binding:"required"`
		CustomMessage string `json:"custom_message,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	newStatus := Status(req.Status)

	// Validate status
	if !sh.statusService.GetStatusManager().IsValidStatus(newStatus) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid status. Must be one of: online, away, busy",
		})
		return
	}

	// Update user status
	err := sh.statusService.SetUserStatus(user.ID.Hex(), newStatus, req.CustomMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update status: " + err.Error(),
		})
		return
	}

	// Get updated user data
	updatedUser, err := sh.mongoDB.GetUserByID(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch updated user data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status updated successfully",
		"user":    updatedUser,
	})
}

// GetUserStatus returns the current status of the authenticated user
func (sh *StatusHandlers) GetUserStatus(c *gin.Context) {
	// Get user from context (set by AuthMiddleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not found in context",
		})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user data in context",
		})
		return
	}

	statusInfo, err := sh.statusService.GetUserStatusInfo(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get status info: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "User status retrieved",
		"status_info": statusInfo,
	})
}

// CheckUserAction checks if a user can perform a specific action
func (sh *StatusHandlers) CheckUserAction(c *gin.Context) {
	// Get user from context (set by AuthMiddleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not found in context",
		})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user data in context",
		})
		return
	}

	action := c.Query("action")
	if action == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Action parameter is required",
		})
		return
	}

	err := sh.statusService.ValidateStatusForAction(user.ID.Hex(), action)
	if err != nil {
		if IsStatusActionError(err) {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": err.Error(),
				"allowed": false,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to check action: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Action allowed",
		"allowed": true,
	})
}

// GetStatusCapabilities returns what the current user can do based on their status
func (sh *StatusHandlers) GetStatusCapabilities(c *gin.Context) {
	// Get user from context (set by AuthMiddleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not found in context",
		})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user data in context",
		})
		return
	}

	statusInfo, err := sh.statusService.GetUserStatusInfo(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get status info: " + err.Error(),
		})
		return
	}

	capabilities := map[string]bool{
		"can_receive_messages":      statusInfo.CanReceiveMessages,
		"can_send_messages":         statusInfo.CanSendMessages,
		"can_join_rooms":            statusInfo.CanJoinRooms,
		"can_receive_notifications": statusInfo.CanReceiveNotifications,
		"can_receive_popups":        statusInfo.CanReceivePopups,
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "Status capabilities retrieved",
		"status":       statusInfo.Status,
		"capabilities": capabilities,
	})
}
