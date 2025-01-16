package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cdt-eth/htmx-chat/internal/models"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	messages := models.GetMessages()
	
	// Return HTML instead of JSON
	w.Header().Set("Content-Type", "text/html")
	for _, msg := range messages {
		fmt.Fprintf(w, "<div class='message'>%s: %s</div>", 
			msg.Sender, msg.Content)
	}
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data instead of JSON
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get values from form
	content := r.FormValue("content")
	sender := r.FormValue("sender")

	// Add message
	newMsg := models.AddMessage(content, sender)

	// Return HTML for the new message
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<div class='message'>%s: %s</div>", 
		newMsg.Sender, newMsg.Content)
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