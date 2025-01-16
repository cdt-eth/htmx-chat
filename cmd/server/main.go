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
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // Get port from environment variable or fallback to 8080
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" 
    }

    
    // Chat routes
    http.HandleFunc("/chat/messages", handlers.GetMessages)
    http.HandleFunc("/chat/send", handlers.SendMessage)
    http.HandleFunc("/chat/delete", handlers.DeleteMessage)
    
	// Home route
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(template.ParseFiles("templates/index.html"))
        tmpl.Execute(w, nil)
    })

    // Start server
    log.Printf("Server starting on port %s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
}
