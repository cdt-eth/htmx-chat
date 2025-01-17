package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/cdt-eth/htmx-chat/internal/models"
	"github.com/gorilla/websocket"
)

// Think of this as a TV broadcast station
// - Each connected client is like a TV in someone's house
// - When we have a new message, we want to broadcast it to all TVs
type Hub struct {
    // clients holds all connected clients (like a list of all TVs tuned to our channel)
    clients map[*websocket.Conn]bool
    
    // broadcast is our message channel (like the TV signal we're broadcasting)
    broadcast chan models.Message
    
    // mutex prevents data races (like making sure we don't change the TV channel while someone's recording)
    mu sync.Mutex
}

// Create our hub (like setting up our TV station)
var hub = Hub{
    clients:   make(map[*websocket.Conn]bool),
    broadcast: make(chan models.Message),
}

// Configure websocket (like setting up our broadcast equipment)
var upgrader = websocket.Upgrader{
    // Allow any origin (in production, you'd restrict this)
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

// Add this type
type WSMessage struct {
    Type    string         `json:"type,omitempty"`
    ID      int           `json:"id,omitempty"`
    Content string        `json:"content,omitempty"`
    Sender  string        `json:"sender,omitempty"`
}

// HandleWebSocket upgrades HTTP to WebSocket
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    // Upgrade HTTP connection to WebSocket (like converting a regular TV antenna to a digital one)
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade connection: %v", err)
        return
    }
    defer conn.Close()

    // Register new client (like adding a new TV to our broadcast list)
    hub.mu.Lock()
    hub.clients[conn] = true
    hub.mu.Unlock()

    // Remove client when they disconnect (like removing a TV that's turned off)
    defer func() {
        hub.mu.Lock()
        delete(hub.clients, conn)
        hub.mu.Unlock()
    }()

    // Listen for messages from this client
    for {
        var wsMsg WSMessage
        if err := conn.ReadJSON(&wsMsg); err != nil {
            log.Printf("Error reading message: %v", err)
            break
        }

        if wsMsg.Type == "delete" {
            if err := models.DeleteMessage(wsMsg.ID); err != nil {
                log.Printf("Error deleting message: %v", err)
                continue
            }
            // Broadcast deletion to all clients
            hub.broadcast <- models.Message{ID: wsMsg.ID, Content: "DELETED"}
        } else {
            newMsg := models.AddMessage(wsMsg.Content, wsMsg.Sender)
            hub.broadcast <- newMsg
        }
    }
}

// Start listening for messages to broadcast (like starting our TV station)
func init() {
    go func() {
        for {
            // Wait for message to broadcast
            msg := <-hub.broadcast

            // Send to all connected clients (like broadcasting to all TVs)
            hub.mu.Lock()
            for client := range hub.clients {
                if err := client.WriteJSON(msg); err != nil {
                    log.Printf("Error broadcasting to client: %v", err)
                    client.Close()
                    delete(hub.clients, client)
                }
            }
            hub.mu.Unlock()
        }
    }()
} 