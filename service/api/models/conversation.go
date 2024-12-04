package models

// Conversation represents a conversation entity
type Conversation struct {
	ConversationID string `json:"conversationId"` // Unique identifier for the conversation
	IsGroup        bool   `json:"isGroup"`        // Indicates if the conversation is a group
	GroupName      string `json:"groupName"`      // Name of the group (optional, only for group conversations)
	GroupPhoto     string `json:"groupPhoto"`     // Photo of the group (optional, only for group conversations)
}
