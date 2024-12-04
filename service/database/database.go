/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"

	api "AlChats/service/api/models"

	"os"

	"github.com/sirupsen/logrus"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetUserByID(userID string) (api.User, error)
	SetUser(username string) (api.User, error)
	GetAllUsers() ([]api.User, error)
	DeleteUserByID(userID string) error
	UpdateUsername(userId string, newUsername string) (api.User, error)

	SetConversation(userIDs []string, isGroup bool, groupName, groupPhoto string) (api.Conversation, error)
	GetAllConversations() ([]api.Conversation, error)
	GetAllConversationsByMember(userID string) ([]api.Conversation, error)
	GetConversationMembers(conversationID string) ([]api.User, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	if err := initializeDatabase(db); err != nil {
		return nil, errors.New("failed to initialize database")
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

func initializeDatabase(db *sql.DB) error {
	// User table schema
	user_table_schema := `
		CREATE TABLE user_table (
			UserID TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))), 
			Username TEXT NOT NULL UNIQUE, 
			Photo TEXT 
		);
	`

	// Conversation table schema
	conversation_table_schema := `
		CREATE TABLE conversation_table (
			ConversationID TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))), -- Generate a UUID
			IsGroup BOOLEAN NOT NULL CHECK (IsGroup IN (0, 1)), -- Boolean value for group indicator
			GroupName TEXT, -- Optional group name
			GroupPhoto TEXT -- Optional group photo
		);
	`

	user_conversation_table_schema := `
		CREATE TABLE user_conversation_table (
			UserID TEXT NOT NULL,                       -- User ID
			ConversationID TEXT NOT NULL,               -- Conversation ID
			PRIMARY KEY (UserID, ConversationID),       -- Composite primary key
			FOREIGN KEY (UserID) REFERENCES user_table(UserID) ON DELETE CASCADE,  -- Link to user_table
			FOREIGN KEY (ConversationID) REFERENCES conversation_table(ConversationID) ON DELETE CASCADE -- Link to conversation_table
		);
	`

	// Initialize TABLES
	if err := ensureTableExists(db, "user_table", user_table_schema); err != nil {
		return err
	}

	if err := ensureTableExists(db, "conversation_table", conversation_table_schema); err != nil {
		return err
	}

	if err := ensureTableExists(db, "user_conversation_table", user_conversation_table_schema); err != nil {
		return err
	}

	return nil
}

func ensureTableExists(db *sql.DB, tableName, createStmt string) error {

	// Init logging
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	// logger.SetLevel(logrus.InfoLevel)

	logger.Infof("ensureTableExists %s", tableName)

	var existingTable string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?;`, tableName).Scan(&existingTable)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Infof("do not exists %s", tableName)
		_, err := db.Exec(createStmt)
		logger.Infof("ddl command result %w", err)
		if err != nil {
			return fmt.Errorf("error creating table %s: %w", tableName, err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking for table %s: %w", tableName, err)
	}
	return nil
}
