package models

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Sender    string    `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	messages = make([]Message, 0)
	mu      sync.Mutex
	nextID  int
)

// AddMessage adds a new message to the store
func AddMessage(content, sender string) Message {
	mu.Lock()
	defer mu.Unlock()

	msg := Message{
		ID:        nextID,
		Content:   content,
		Sender:    sender,
		Timestamp: time.Now(),
	}
	messages = append(messages, msg)
	nextID++
	return msg
}

// GetMessages returns all messages
func GetMessages() []Message {
	mu.Lock()
	defer mu.Unlock()
	return messages
}


func DeleteMessage(id int) error {
	mu.Lock()
	defer mu.Unlock()
	
	fmt.Printf("Trying to delete message with ID: %d\n", id)
	for i, msg := range messages {
		fmt.Printf("Checking message: %+v\n", msg)
		if msg.ID == id {
			messages = append(messages[:i], messages[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("message with ID %d not found", id)
}
