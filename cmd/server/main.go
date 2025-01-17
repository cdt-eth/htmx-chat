package main

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/cdt-eth/htmx-chat/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {    
    // Don't fail if .env missing
    godotenv.Load() // Remove error check

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" 
    }

    // Auth routes first (most specific)
    http.HandleFunc("/auth/login", handlers.LoginHandler)

    // Chat routes
    http.HandleFunc("/chat/messages", handlers.GetMessages)
    http.HandleFunc("/chat/send", handlers.SendMessage)
    http.HandleFunc("/chat/delete", handlers.DeleteMessage)
    
    // WebSocket endpoint
    http.HandleFunc("/ws", handlers.HandleWebSocket)

    // Root route 
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("token")
        if err != nil || cookie.Value == "" {
            http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
            return
        }
        tmpl := template.Must(template.ParseFiles("templates/index.html"))
        tmpl.Execute(w, nil)
    })

    log.Printf("Server starting on port %s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
}
