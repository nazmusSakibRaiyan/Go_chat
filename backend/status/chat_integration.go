package status

// ChatIntegration provides status-aware chat functionality
type ChatIntegration struct {
	statusService *StatusService
}

// NewChatIntegration creates a new chat integration
func NewChatIntegration(statusService *StatusService) *ChatIntegration {
	return &ChatIntegration{
		statusService: statusService,
	}
}

// FilterMessageRecipients filters users who should receive a message based on their status
func (ci *ChatIntegration) FilterMessageRecipients(userIDs []string, messageType string) []string {
	var allowedUsers []string

	for _, userID := range userIDs {
		if ci.ShouldReceiveMessage(userID, messageType) {
			allowedUsers = append(allowedUsers, userID)
		}
	}

	return allowedUsers
}

// ShouldReceiveMessage checks if a user should receive a message based on their status
func (ci *ChatIntegration) ShouldReceiveMessage(userID string, messageType string) bool {
	return ci.statusService.ShouldReceiveNotification(userID, messageType)
}

// CanUserJoinRoom checks if a user can join a specific room
func (ci *ChatIntegration) CanUserJoinRoom(userID string, roomID string) (bool, string) {
	return ci.statusService.CanUserJoinRoom(userID)
}

// CanUserSendMessageToRoom checks if a user can send a message to a room
func (ci *ChatIntegration) CanUserSendMessageToRoom(userID string, roomID string) (bool, string) {
	return ci.statusService.CanUserSendMessage(userID)
}

// GetRoomUserStatuses returns status information for all users in a room
func (ci *ChatIntegration) GetRoomUserStatuses(userIDs []string) map[string]*StatusInfo {
	userStatuses := make(map[string]*StatusInfo)

	for _, userID := range userIDs {
		statusInfo, err := ci.statusService.GetUserStatusInfo(userID)
		if err == nil {
			userStatuses[userID] = statusInfo
		}
	}

	return userStatuses
}

// CreateStatusAwareMessage creates a message with status-aware delivery rules
func (ci *ChatIntegration) CreateStatusAwareMessage(fromUserID string, content string, roomID string) (*StatusAwareMessage, error) {
	// Check if sender can send messages
	canSend, reason := ci.statusService.CanUserSendMessage(fromUserID)
	if !canSend {
		return nil, &StatusActionError{
			UserStatus: Status("unknown"),
			Action:     "send_messages",
			Message:    reason,
		}
	}

	senderStatusInfo, err := ci.statusService.GetUserStatusInfo(fromUserID)
	if err != nil {
		return nil, err
	}

	return &StatusAwareMessage{
		FromUserID:    fromUserID,
		Content:       content,
		RoomID:        roomID,
		SenderStatus:  senderStatusInfo.Status,
		DeliveryRules: ci.createDeliveryRules(senderStatusInfo),
		RequiresAck:   senderStatusInfo.Status == Busy, // Busy users might want acknowledgment
	}, nil
}

// StatusAwareMessage represents a message with status-based delivery rules
type StatusAwareMessage struct {
	FromUserID    string        `json:"from_user_id"`
	Content       string        `json:"content"`
	RoomID        string        `json:"room_id"`
	SenderStatus  Status        `json:"sender_status"`
	DeliveryRules DeliveryRules `json:"delivery_rules"`
	RequiresAck   bool          `json:"requires_ack"`
}

// DeliveryRules defines how a message should be delivered based on status
type DeliveryRules struct {
	DeliverToOnline bool `json:"deliver_to_online"`
	DeliverToAway   bool `json:"deliver_to_away"`
	DeliverToBusy   bool `json:"deliver_to_busy"`
	ShowPopupToAll  bool `json:"show_popup_to_all"`
	PriorityLevel   int  `json:"priority_level"`
}

// createDeliveryRules creates delivery rules based on sender status
func (ci *ChatIntegration) createDeliveryRules(senderStatus *StatusInfo) DeliveryRules {
	rules := DeliveryRules{
		DeliverToOnline: true,  // Always deliver to online users
		DeliverToAway:   true,  // Always store for away users
		DeliverToBusy:   true,  // Always deliver to busy users
		ShowPopupToAll:  false, // Default no popups
		PriorityLevel:   1,     // Default priority
	}

	// Adjust rules based on sender status
	switch senderStatus.Status {
	case Online:
		rules.ShowPopupToAll = true
		rules.PriorityLevel = 2
	case Busy:
		rules.ShowPopupToAll = false // Busy senders don't trigger popups
		rules.PriorityLevel = 3      // High priority from busy users
	case Away:
		rules.ShowPopupToAll = false
		rules.PriorityLevel = 1
	}

	return rules
}

// ProcessMessageDelivery processes a message for delivery based on recipient statuses
func (ci *ChatIntegration) ProcessMessageDelivery(message *StatusAwareMessage, recipientIDs []string) *MessageDeliveryResult {
	result := &MessageDeliveryResult{
		Message:           message,
		OnlineRecipients:  []string{},
		AwayRecipients:    []string{},
		BusyRecipients:    []string{},
		PopupRecipients:   []string{},
		StoredRecipients:  []string{},
		BlockedRecipients: []string{},
	}

	for _, recipientID := range recipientIDs {
		statusInfo, err := ci.statusService.GetUserStatusInfo(recipientID)
		if err != nil {
			result.BlockedRecipients = append(result.BlockedRecipients, recipientID)
			continue
		}

		// Categorize recipients by status
		switch statusInfo.Status {
		case Online:
			result.OnlineRecipients = append(result.OnlineRecipients, recipientID)
			if message.DeliveryRules.ShowPopupToAll && statusInfo.CanReceivePopups {
				result.PopupRecipients = append(result.PopupRecipients, recipientID)
			}
		case Away:
			result.AwayRecipients = append(result.AwayRecipients, recipientID)
			result.StoredRecipients = append(result.StoredRecipients, recipientID)
		case Busy:
			result.BusyRecipients = append(result.BusyRecipients, recipientID)
			// Busy users receive messages but no popups
			if message.DeliveryRules.ShowPopupToAll && message.DeliveryRules.PriorityLevel >= 3 {
				// Only high priority messages get popups for busy users
				result.PopupRecipients = append(result.PopupRecipients, recipientID)
			}
		}
	}

	return result
}

// MessageDeliveryResult contains the result of message delivery processing
type MessageDeliveryResult struct {
	Message           *StatusAwareMessage `json:"message"`
	OnlineRecipients  []string            `json:"online_recipients"`
	AwayRecipients    []string            `json:"away_recipients"`
	BusyRecipients    []string            `json:"busy_recipients"`
	PopupRecipients   []string            `json:"popup_recipients"`
	StoredRecipients  []string            `json:"stored_recipients"`
	BlockedRecipients []string            `json:"blocked_recipients"`
}

// GetDeliveryStats returns statistics about message delivery
func (result *MessageDeliveryResult) GetDeliveryStats() map[string]int {
	return map[string]int{
		"total_recipients":   len(result.OnlineRecipients) + len(result.AwayRecipients) + len(result.BusyRecipients),
		"online_recipients":  len(result.OnlineRecipients),
		"away_recipients":    len(result.AwayRecipients),
		"busy_recipients":    len(result.BusyRecipients),
		"popup_recipients":   len(result.PopupRecipients),
		"stored_recipients":  len(result.StoredRecipients),
		"blocked_recipients": len(result.BlockedRecipients),
	}
}
