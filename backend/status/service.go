package status

import (
	"go-chat-backend/db"
	"go-chat-backend/models"
)

// StatusService provides high-level status management operations
type StatusService struct {
	statusManager *StatusManager
	mongoDB       *db.MongoDB
}

// NewStatusService creates a new status service
func NewStatusService(mongoDB *db.MongoDB) *StatusService {
	return &StatusService{
		statusManager: NewStatusManager(),
		mongoDB:       mongoDB,
	}
}

// GetStatusManager returns the underlying status manager
func (ss *StatusService) GetStatusManager() *StatusManager {
	return ss.statusManager
}

// SetUserStatus sets a user's status with validation and business logic
func (ss *StatusService) SetUserStatus(userID string, newStatus Status, customMessage string) error {
	// Get current user to validate status transition
	user, err := ss.mongoDB.GetUserByID(userID)
	if err != nil {
		return err
	}

	currentStatus := Status(user.Status)

	// Validate status transition
	if err := ss.statusManager.ValidateStatusTransition(currentStatus, newStatus); err != nil {
		return err
	}

	// Update user status in database
	statusString := string(newStatus)
	err = ss.mongoDB.UpdateUserProfile(userID, user.DisplayName, &user.Avatar, statusString)
	if err != nil {
		return err
	}

	// TODO: Store detailed status information in separate collection if needed
	// This could include custom messages, expiration times, etc.

	return nil
}

// CanUserJoinRoom checks if a user can join a room based on their status
func (ss *StatusService) CanUserJoinRoom(userID string) (bool, string) {
	user, err := ss.mongoDB.GetUserByID(userID)
	if err != nil {
		return false, "Unable to verify user status"
	}

	userStatus := Status(user.Status)
	if !ss.statusManager.CanUserPerformAction(userStatus, "join_rooms") {
		statusInfo, _ := ss.statusManager.GetStatusInfo(userStatus)
		return false, "Cannot join rooms while " + statusInfo.DisplayName
	}

	return true, ""
}

// CanUserSendMessage checks if a user can send messages based on their status
func (ss *StatusService) CanUserSendMessage(userID string) (bool, string) {
	user, err := ss.mongoDB.GetUserByID(userID)
	if err != nil {
		return false, "Unable to verify user status"
	}

	userStatus := Status(user.Status)
	if !ss.statusManager.CanUserPerformAction(userStatus, "send_messages") {
		statusInfo, _ := ss.statusManager.GetStatusInfo(userStatus)
		return false, "Cannot send messages while " + statusInfo.DisplayName
	}

	return true, ""
}

// ShouldReceiveNotification checks if a user should receive a notification
func (ss *StatusService) ShouldReceiveNotification(userID string, notificationType string) bool {
	user, err := ss.mongoDB.GetUserByID(userID)
	if err != nil {
		return false
	}

	userStatus := Status(user.Status)

	switch notificationType {
	case "popup":
		return ss.statusManager.CanUserPerformAction(userStatus, "receive_popups")
	case "notification":
		return ss.statusManager.CanUserPerformAction(userStatus, "receive_notifications")
	default:
		return ss.statusManager.CanUserPerformAction(userStatus, "receive_notifications")
	}
}

// GetUserStatusInfo gets detailed status information for a user
func (ss *StatusService) GetUserStatusInfo(userID string) (*StatusInfo, error) {
	user, err := ss.mongoDB.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	userStatus := Status(user.Status)
	return ss.statusManager.GetStatusInfo(userStatus)
}

// GetAllUsersInRoom gets status information for all users in a room
func (ss *StatusService) GetAllUsersInRoom(roomID string) (map[string]*StatusInfo, error) {
	// TODO: Implement when room membership is tracked
	// For now, return empty map
	return make(map[string]*StatusInfo), nil
}

// FilterUsersByStatus filters a list of user IDs by their status
func (ss *StatusService) FilterUsersByStatus(userIDs []string, allowedStatuses []Status) ([]string, error) {
	var filteredUsers []string

	for _, userID := range userIDs {
		user, err := ss.mongoDB.GetUserByID(userID)
		if err != nil {
			continue // Skip users we can't fetch
		}

		userStatus := Status(user.Status)
		for _, allowedStatus := range allowedStatuses {
			if userStatus == allowedStatus {
				filteredUsers = append(filteredUsers, userID)
				break
			}
		}
	}

	return filteredUsers, nil
}

// GetOnlineUsers returns all users with "online" status
func (ss *StatusService) GetOnlineUsers() ([]*models.User, error) {
	// TODO: Implement efficient query when needed
	// This would require an index on the status field
	return nil, nil
}

// UpdateUserActivity updates a user's last activity time
func (ss *StatusService) UpdateUserActivity(userID string) error {
	// For now, we'll just update the user's updated_at timestamp
	// In a full implementation, we'd have a separate activity tracking
	user, err := ss.mongoDB.GetUserByID(userID)
	if err != nil {
		return err
	}

	// If user is auto-away, bring them back online
	if user.Status == string(Away) {
		return ss.SetUserStatus(userID, Online, "")
	}

	return nil
}

// AutoUpdateStatuses checks and updates statuses for users who should be auto-away
func (ss *StatusService) AutoUpdateStatuses() error {
	// TODO: Implement periodic status updates
	// This would run as a background job to set users to "away" after inactivity
	return nil
}

// GetStatusStats returns statistics about user statuses
func (ss *StatusService) GetStatusStats() (map[Status]int, error) {
	// TODO: Implement status statistics
	// This could be useful for admin dashboards
	stats := make(map[Status]int)
	return stats, nil
}

// ValidateStatusForAction checks if a user can perform an action and returns appropriate error
func (ss *StatusService) ValidateStatusForAction(userID string, action string) error {
	user, err := ss.mongoDB.GetUserByID(userID)
	if err != nil {
		return err
	}

	userStatus := Status(user.Status)
	if !ss.statusManager.CanUserPerformAction(userStatus, action) {
		statusInfo, _ := ss.statusManager.GetStatusInfo(userStatus)
		return &StatusActionError{
			UserStatus: userStatus,
			Action:     action,
			Message:    "Action not allowed while " + statusInfo.DisplayName,
		}
	}

	return nil
}

// StatusActionError represents an error when a user cannot perform an action due to their status
type StatusActionError struct {
	UserStatus Status
	Action     string
	Message    string
}

func (e *StatusActionError) Error() string {
	return e.Message
}

// IsStatusActionError checks if an error is a status action error
func IsStatusActionError(err error) bool {
	_, ok := err.(*StatusActionError)
	return ok
}
