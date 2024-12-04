package models

// User represents a user entity
type User struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Photo    string `json:"photo,omitempty"`
}
