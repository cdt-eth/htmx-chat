package handlers

import (
	"fmt"
	"net/http"

	"github.com/cdt-eth/htmx-chat/internal/models"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	messages := models.GetMessages()
	
	// Return HTML instead of JSON
	w.Header().Set("Content-Type", "text/html")
	for _, msg := range messages {
		fmt.Fprintf(w, `
			<div class="message" id="msg-%d">
				<span>%s: %s</span>
				<button 
					class="delete-btn"
					hx-delete="/chat/delete"
					hx-vals='{"id": %d}'
					hx-target="#msg-%d"
					hx-swap="outerHTML">
					×
				</button>
			</div>`, 
			msg.ID, msg.Sender, msg.Content, msg.ID, msg.ID)
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
	fmt.Fprintf(w, `
		<div class="message" id="msg-%d">
			<span>%s: %s</span>
			<button 
				class="delete-btn"
				hx-delete="/chat/delete"
				hx-vals='{"id": %d}'
				hx-target="#msg-%d"
				hx-swap="outerHTML">
				×
			</button>
		</div>`, 
		newMsg.ID, newMsg.Sender, newMsg.Content, newMsg.ID, newMsg.ID)
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get ID from form values
	var id int
	fmt.Sscanf(r.FormValue("id"), "%d", &id)

	// Delete message
	if err := models.DeleteMessage(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}