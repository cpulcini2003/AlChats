package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type ConversationRequest struct {
	UserIDs    []string `json:"user_ids"`
	IsGroup    bool     `json:"is_group"`
	GroupName  string   `json:"group_name,omitempty"`
	GroupPhoto string   `json:"group_photo,omitempty"`
}

// SetConversationHandler handles the creation of new conversations.
func (h *_router) setConversationHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Parse the JSON request body
	var req ConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate the required fields
	if len(req.UserIDs) == 0 {
		http.Error(w, `{"error":"user_ids is required"}`, http.StatusBadRequest)
		return
	}

	// Call the SetConversation function
	conversation, err := h.db.SetConversation(req.UserIDs, req.IsGroup, req.GroupName, req.GroupPhoto)
	if err != nil {
		if strings.Contains(err.Error(), "cannot create a conversation with only one user") ||
			strings.Contains(err.Error(), "cannot create a group conversation") {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusBadRequest)
		} else if strings.Contains(err.Error(), "user with UserID") {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		}
		return
	}

	// Respond with the created conversation
	if err := json.NewEncoder(w).Encode(conversation); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to encode response: %v"}`, err), http.StatusInternalServerError)
	}
}
