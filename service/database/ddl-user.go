package database

import (
	api "AlChats/service/api/models"
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

func (db *appdbimpl) GetUserByID(userID string) (api.User, error) {
	var user api.User
	err := db.c.QueryRow("SELECT UserID, Username, Photo FROM user_table WHERE UserID = ?", userID).Scan(&user.UserID, &user.Username, &user.Photo)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (db *appdbimpl) UpdateUsername(userId string, newUsername string) (api.User, error) {
	var user api.User

	// SQL to update the username, ensuring no duplicate usernames exist
	query := `
		UPDATE user_table 
		SET Username = ? 
		WHERE UserID = ? 
		AND NOT EXISTS (
			SELECT 1 
			FROM user_table 
			WHERE Username = ?
			AND UserID != ?
		)
		RETURNING UserID, Username, COALESCE(Photo, '')
	`

	// Update the username and fetch the updated fields
	err := db.c.QueryRow(query, newUsername, userId, newUsername, userId).Scan(&user.UserID, &user.Username, &user.Photo)
	if err != nil {
		// Check if the error is a constraint violation
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			return user, fmt.Errorf("username %q already exists", newUsername)
		}

		// Handle case when no rows are affected (e.g., userId not found)
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user with ID %q not found or username %q already exists", userId, newUsername)
		}

		return user, err
	}

	return user, nil
}

func (db *appdbimpl) SetUser(username string) (api.User, error) {
	var user api.User

	// SQL to insert a new user and retrieve the generated UserID and other fields
	query := `
		INSERT INTO user_table (Username) 
		VALUES (?)
		RETURNING UserID, Username, COALESCE(Photo, '')
	`

	// Insert the user and fetch the generated fields
	err := db.c.QueryRow(query, username).Scan(&user.UserID, &user.Username, &user.Photo)
	if err != nil {
		// Check if the error is a unique constraint violation
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			return user, fmt.Errorf("username %q already exists", username)
		}
		return user, err
	}

	return user, nil
}

func (db *appdbimpl) GetAllUsers() ([]api.User, error) {
	var users []api.User

	// Query to select all users
	rows, err := db.c.Query("SELECT UserID, Username, COALESCE(Photo, '') FROM user_table")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows and scan data into the users slice
	for rows.Next() {
		var user api.User
		err := rows.Scan(&user.UserID, &user.Username, &user.Photo)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (db *appdbimpl) DeleteUserByID(userID string) error {
	// SQL query to delete a user by UserID
	query := "DELETE FROM user_table WHERE UserID = ?"

	// Execute the query
	result, err := db.c.Exec(query, userID)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with UserID %q", userID)
	}

	return nil
}
