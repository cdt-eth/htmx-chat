package models

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/cdt-eth/htmx-chat/internal/db"
	"golang.org/x/crypto/bcrypt"
)

// Validation constants
const (
	MinUsernameLength = 3
	MaxUsernameLength = 30
	MinPasswordLength = 8
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Won't be sent in JSON responses
}

// ValidateUsername checks username requirements
func ValidateUsername(username string) error {
	length := len(username)
	if length < MinUsernameLength {
		return fmt.Errorf("username must be at least %d characters", MinUsernameLength)
	}
	if length > MaxUsernameLength {
		return fmt.Errorf("username must be less than %d characters", MaxUsernameLength)
	}
	
	// Check for valid characters (letters, numbers, underscores only)
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) && char != '_' {
			return fmt.Errorf("username can only contain letters, numbers, and underscores")
		}
	}
	
	return nil
}

// ValidatePassword checks password strength
func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters", MinPasswordLength)
	}
	
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasNumber := strings.ContainsAny(password, "0123456789")
	
	if !hasUpper || !hasLower || !hasNumber {
		return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}
	
	return nil
}

// CreateUser handles registration with duplicate username check
func CreateUser(username, password string) (*User, error) {
	// Validate input
	if err := ValidateUsername(username); err != nil {
		return nil, err
	}
	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	// Insert user
	var user User
	err = db.DB.QueryRow(`
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, username
	`, username, string(hashedPassword)).Scan(&user.ID, &user.Username)

	// Handle unique constraint violation
	if err != nil && err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
		return nil, fmt.Errorf("username already taken")
	}

	return &user, err
}

// GetUserByUsername for login
func GetUserByUsername(username string) (*User, error) {
	var user User
	var hashedPassword string
	
	err := db.DB.QueryRow(`
		SELECT id, username, password_hash 
		FROM users 
		WHERE username = $1
	`, username).Scan(&user.ID, &user.Username, &hashedPassword)
	
	user.Password = hashedPassword
	return &user, err
} 