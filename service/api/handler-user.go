package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) createUserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Get the `username` query parameter
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, `{"error":"username parameter is required"}`, http.StatusBadRequest)
		return
	}

	// Call the SetUser function to create a new user
	user, err := rt.db.SetUser(username)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, `{"error":"username already exists"}`, http.StatusConflict)
		} else {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		}
		return
	}

	// Write the created user as a JSON response
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to encode response: %v"}`, err), http.StatusInternalServerError)
	}
}

func (rt *_router) updateUsernameHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Get the `userId` and `newUsername` query parameters
	userId := r.URL.Query().Get("userId")
	newUsername := r.URL.Query().Get("newUsername")

	// Validate that both parameters are provided
	if userId == "" {
		http.Error(w, `{"error":"userId parameter is required"}`, http.StatusBadRequest)
		return
	}
	if newUsername == "" {
		http.Error(w, `{"error":"newUsername parameter is required"}`, http.StatusBadRequest)
		return
	}

	// Call the UpdateUsername function to update the username
	user, err := rt.db.UpdateUsername(userId, newUsername)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, `{"error":"username already exists"}`, http.StatusConflict)
		} else if strings.Contains(err.Error(), "not found") {
			http.Error(w, `{"error":"user not found or conflict occurred"}`, http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		}
		return
	}

	// Write the updated user as a JSON response
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to encode response: %v"}`, err), http.StatusInternalServerError)
	}
}

func (rt *_router) getAllUsersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Call GetAllUsers to fetch all users
	users, err := rt.db.GetAllUsers()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	// Write the list of users as a JSON response
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to encode response: %v"}`, err), http.StatusInternalServerError)
	}
}
