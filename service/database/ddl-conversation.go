package database

import (
	api "AlChats/service/api/models"
	"fmt"
)

func (db *appdbimpl) SetConversation(userIDs []string, isGroup bool, groupName, groupPhoto string) (api.Conversation, error) {
	var conversation api.Conversation

	// Check for invalid userIDs length
	if len(userIDs) == 1 {
		return conversation, fmt.Errorf("cannot create a conversation with only one user")
	}

	// Check if userIDs is greater than 2 and isGroup is false
	if len(userIDs) > 2 && !isGroup {
		return conversation, fmt.Errorf("cannot create a group conversation with more than two users without setting isGroup to true")
	}

	// SQL to insert a new conversation and retrieve the generated ConversationID and other fields
	query := `
		INSERT INTO conversation_table (IsGroup, GroupName, GroupPhoto) 
		VALUES (?, ?, ?)
		RETURNING 
			ConversationID, 
			IsGroup, 
			COALESCE(GroupName, ''), 
			COALESCE(GroupPhoto, '')
	`

	// Insert the conversation and fetch the generated fields
	err := db.c.QueryRow(query, isGroup, groupName, groupPhoto).
		Scan(&conversation.ConversationID, &conversation.IsGroup, &conversation.GroupName, &conversation.GroupPhoto)
	if err != nil {
		// Handle any database error
		return conversation, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Insert user-conversation relationships into the user_conversation_table
	for _, userID := range userIDs {
		// Check if the UserID exists in the user_table
		var exists bool
		checkUserQuery := `SELECT EXISTS(SELECT 1 FROM user_table WHERE UserID = ?)`
		err := db.c.QueryRow(checkUserQuery, userID).Scan(&exists)
		if err != nil {
			return conversation, fmt.Errorf("failed to check if user exists: %w", err)
		}

		// If the user does not exist, return an error
		if !exists {
			return conversation, fmt.Errorf("user with UserID %s does not exist", userID)
		}

		// SQL to insert the relationship between user and conversation
		relationshipQuery := `
			INSERT INTO user_conversation_table (UserID, ConversationID)
			VALUES (?, ?)
		`

		// Insert the user-conversation relationship
		_, err = db.c.Exec(relationshipQuery, userID, conversation.ConversationID)
		if err != nil {
			// If an error occurs, return the error
			return conversation, fmt.Errorf("failed to create user-conversation relationship: %w", err)
		}
	}

	return conversation, nil
}

func (db *appdbimpl) GetAllConversations() ([]api.Conversation, error) {
	var conversations []api.Conversation

	// SQL to select all conversations from the conversation_table
	query := `
		SELECT 
			ConversationID, 
			IsGroup, 
			COALESCE(GroupName, ''), 
			COALESCE(GroupPhoto, '') 
		FROM conversation_table
	`

	// Execute the query and iterate through the results
	rows, err := db.c.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve conversations: %w", err)
	}
	defer rows.Close()

	// Loop through the rows and map them to the Conversation struct
	for rows.Next() {
		var conversation api.Conversation
		err := rows.Scan(&conversation.ConversationID, &conversation.IsGroup, &conversation.GroupName, &conversation.GroupPhoto)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation row: %w", err)
		}
		conversations = append(conversations, conversation)
	}

	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over conversation rows: %w", err)
	}

	return conversations, nil
}

func (db *appdbimpl) GetAllConversationsByMember(userID string) ([]api.Conversation, error) {
	var conversations []api.Conversation

	// SQL to select all conversations for a user by joining conversation_table and user_conversation_table
	query := `
		SELECT 
			c.ConversationID, 
			c.IsGroup, 
			COALESCE(c.GroupName, ''), 
			COALESCE(c.GroupPhoto, '') 
		FROM conversation_table c
		JOIN user_conversation_table uc ON c.ConversationID = uc.ConversationID
		WHERE uc.UserID = ?
	`

	// Execute the query
	rows, err := db.c.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve conversations for user %s: %w", userID, err)
	}
	defer rows.Close()

	// Loop through the rows and map them to the Conversation struct
	for rows.Next() {
		var conversation api.Conversation
		err := rows.Scan(&conversation.ConversationID, &conversation.IsGroup, &conversation.GroupName, &conversation.GroupPhoto)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation row: %w", err)
		}
		conversations = append(conversations, conversation)
	}

	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over conversation rows: %w", err)
	}

	return conversations, nil
}

func (db *appdbimpl) GetConversationMembers(conversationID string) ([]api.User, error) {
	var members []api.User

	// SQL query to retrieve all users associated with a given conversation
	query := `
		SELECT 
			u.UserID, 
			u.Username, 
			COALESCE(u.Photo, '') 
		FROM user_table u
		JOIN user_conversation_table uc ON u.UserID = uc.UserID
		WHERE uc.ConversationID = ?
	`

	// Execute the query
	rows, err := db.c.Query(query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve members for conversation %s: %w", conversationID, err)
	}
	defer rows.Close()

	// Iterate through the rows and map them to the User struct
	for rows.Next() {
		var user api.User
		err := rows.Scan(&user.UserID, &user.Username, &user.Photo)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		members = append(members, user)
	}

	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over user rows: %w", err)
	}

	return members, nil
}
