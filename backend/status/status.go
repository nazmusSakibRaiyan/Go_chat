package status

import (
	"errors"
	"time"
)

// Status represents a user's availability status
type Status string

// Status constants
const (
	Online Status = "online"
	Away   Status = "away"
	Busy   Status = "busy"
)

// StatusInfo contains detailed information about a status
type StatusInfo struct {
	Status                  Status        `json:"status"`
	DisplayName             string        `json:"display_name"`
	Description             string        `json:"description"`
	Icon                    string        `json:"icon"`
	CanReceiveMessages      bool          `json:"can_receive_messages"`
	CanSendMessages         bool          `json:"can_send_messages"`
	CanJoinRooms            bool          `json:"can_join_rooms"`
	CanReceiveNotifications bool          `json:"can_receive_notifications"`
	CanReceivePopups        bool          `json:"can_receive_popups"`
	AutoAwayTimeout         time.Duration `json:"auto_away_timeout,omitempty"`
	Priority                int           `json:"priority"` // Higher number = higher priority
}

// UserStatus represents a user's current status with additional metadata
type UserStatus struct {
	UserID        string     `json:"user_id" bson:"user_id"`
	Status        Status     `json:"status" bson:"status"`
	CustomMessage string     `json:"custom_message,omitempty" bson:"custom_message,omitempty"`
	SetAt         time.Time  `json:"set_at" bson:"set_at"`
	LastActivity  time.Time  `json:"last_activity" bson:"last_activity"`
	AutoSet       bool       `json:"auto_set" bson:"auto_set"` // True if status was set automatically
	ExpiresAt     *time.Time `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
}

// StatusManager handles all status-related operations and business logic
type StatusManager struct {
	statusDefinitions map[Status]*StatusInfo
}

// NewStatusManager creates a new status manager with default configurations
func NewStatusManager() *StatusManager {
	sm := &StatusManager{
		statusDefinitions: make(map[Status]*StatusInfo),
	}

	// Initialize default status definitions
	sm.initializeDefaultStatuses()

	return sm
}

// initializeDefaultStatuses sets up the default status configurations
func (sm *StatusManager) initializeDefaultStatuses() {
	sm.statusDefinitions[Online] = &StatusInfo{
		Status:                  Online,
		DisplayName:             "Online",
		Description:             "Available and ready to chat",
		Icon:                    "🟢",
		CanReceiveMessages:      true,
		CanSendMessages:         true,
		CanJoinRooms:            true,
		CanReceiveNotifications: true,
		CanReceivePopups:        true,
		Priority:                3,
	}

	sm.statusDefinitions[Away] = &StatusInfo{
		Status:                  Away,
		DisplayName:             "Away",
		Description:             "Temporarily unavailable",
		Icon:                    "🟡",
		CanReceiveMessages:      true,
		CanSendMessages:         false,
		CanJoinRooms:            false,
		CanReceiveNotifications: false,
		CanReceivePopups:        false,
		AutoAwayTimeout:         30 * time.Minute,
		Priority:                1,
	}

	sm.statusDefinitions[Busy] = &StatusInfo{
		Status:                  Busy,
		DisplayName:             "Busy",
		Description:             "Occupied and may not respond immediately",
		Icon:                    "🔴",
		CanReceiveMessages:      true,
		CanSendMessages:         true,
		CanJoinRooms:            true,
		CanReceiveNotifications: true,
		CanReceivePopups:        false, // No popup notifications
		Priority:                2,
	}
}

// GetAllStatuses returns all available status definitions
func (sm *StatusManager) GetAllStatuses() map[Status]*StatusInfo {
	return sm.statusDefinitions
}

// GetStatusInfo returns detailed information about a specific status
func (sm *StatusManager) GetStatusInfo(status Status) (*StatusInfo, error) {
	info, exists := sm.statusDefinitions[status]
	if !exists {
		return nil, errors.New("invalid status")
	}
	return info, nil
}

// IsValidStatus checks if a status is valid
func (sm *StatusManager) IsValidStatus(status Status) bool {
	_, exists := sm.statusDefinitions[status]
	return exists
}

// GetDefaultStatus returns the default status for new users
func (sm *StatusManager) GetDefaultStatus() Status {
	return Online
}

// CanUserPerformAction checks if a user with given status can perform a specific action
func (sm *StatusManager) CanUserPerformAction(status Status, action string) bool {
	statusInfo, err := sm.GetStatusInfo(status)
	if err != nil {
		return false
	}

	switch action {
	case "receive_messages":
		return statusInfo.CanReceiveMessages
	case "send_messages":
		return statusInfo.CanSendMessages
	case "join_rooms":
		return statusInfo.CanJoinRooms
	case "receive_notifications":
		return statusInfo.CanReceiveNotifications
	case "receive_popups":
		return statusInfo.CanReceivePopups
	default:
		return false
	}
}

// CreateUserStatus creates a new user status record
func (sm *StatusManager) CreateUserStatus(userID string, status Status, customMessage string) (*UserStatus, error) {
	if !sm.IsValidStatus(status) {
		return nil, errors.New("invalid status")
	}

	now := time.Now()
	userStatus := &UserStatus{
		UserID:        userID,
		Status:        status,
		CustomMessage: customMessage,
		SetAt:         now,
		LastActivity:  now,
		AutoSet:       false,
	}

	// Set expiration if this is a temporary status
	statusInfo, _ := sm.GetStatusInfo(status)
	if statusInfo.AutoAwayTimeout > 0 && status != Away {
		expiresAt := now.Add(statusInfo.AutoAwayTimeout)
		userStatus.ExpiresAt = &expiresAt
	}

	return userStatus, nil
}

// UpdateUserActivity updates the last activity time for a user
func (sm *StatusManager) UpdateUserActivity(userStatus *UserStatus) {
	userStatus.LastActivity = time.Now()

	// If user was auto-away and becomes active, reset to online
	if userStatus.AutoSet && userStatus.Status == Away {
		userStatus.Status = Online
		userStatus.AutoSet = false
		userStatus.SetAt = time.Now()
		userStatus.ExpiresAt = nil
	}
}

// CheckAutoStatusUpdate checks if a user's status should be automatically updated
func (sm *StatusManager) CheckAutoStatusUpdate(userStatus *UserStatus) (bool, Status) {
	now := time.Now()

	// Check if status has expired
	if userStatus.ExpiresAt != nil && now.After(*userStatus.ExpiresAt) {
		return true, Away
	}

	// Check for auto-away based on inactivity
	statusInfo, err := sm.GetStatusInfo(userStatus.Status)
	if err != nil {
		return false, userStatus.Status
	}

	if statusInfo.AutoAwayTimeout > 0 && userStatus.Status == Online {
		inactiveTime := now.Sub(userStatus.LastActivity)
		if inactiveTime > statusInfo.AutoAwayTimeout {
			return true, Away
		}
	}

	return false, userStatus.Status
}

// GetStatusPriority returns the priority of a status (higher = more important)
func (sm *StatusManager) GetStatusPriority(status Status) int {
	statusInfo, err := sm.GetStatusInfo(status)
	if err != nil {
		return 0
	}
	return statusInfo.Priority
}

// CompareStatuses returns true if status1 has higher priority than status2
func (sm *StatusManager) CompareStatuses(status1, status2 Status) bool {
	return sm.GetStatusPriority(status1) > sm.GetStatusPriority(status2)
}

// ValidateStatusTransition checks if a status transition is allowed
func (sm *StatusManager) ValidateStatusTransition(fromStatus, toStatus Status) error {
	if !sm.IsValidStatus(fromStatus) {
		return errors.New("invalid source status")
	}
	if !sm.IsValidStatus(toStatus) {
		return errors.New("invalid target status")
	}

	// Add custom transition rules here if needed
	// For now, all transitions are allowed

	return nil
}

// GetAvailableStatusesForAPI returns statuses formatted for API response
func (sm *StatusManager) GetAvailableStatusesForAPI() []map[string]interface{} {
	var statuses []map[string]interface{}

	for _, statusInfo := range sm.statusDefinitions {
		statuses = append(statuses, map[string]interface{}{
			"status":       string(statusInfo.Status),
			"display_name": statusInfo.DisplayName,
			"description":  statusInfo.Description,
			"icon":         statusInfo.Icon,
			"priority":     statusInfo.Priority,
		})
	}

	return statuses
}
