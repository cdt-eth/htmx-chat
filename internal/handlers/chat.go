package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cdt-eth/htmx-chat/internal/models"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	messages := models.GetMessages()
	json.NewEncoder(w).Encode(messages)
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg struct {
		Content string `json:"content"`
		Sender  string `json:"sender"`
	}

	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newMsg := models.AddMessage(msg.Content, msg.Sender)
	json.NewEncoder(w).Encode(newMsg)
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse ID from request
	var msg struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delete message
	if err := models.DeleteMessage(msg.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}