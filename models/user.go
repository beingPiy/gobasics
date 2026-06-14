package models

import "time"

// User represents a user in the system
type User struct {
	ID           int       `db:"id" 			json:"id"`             // unique identifier for the user
	Email        string    `db:"email" 			json:"email"`          // user's email address (unique)
	PasswordHash string    `db:"password_hash" 	json:"-"`  // hashed password (not exposed in JSON responses)
	Role         string    `db:"role" 			json:"role"`           // user's role
	CreatedAt    time.Time `db:"created_at" 	json:"created_at"`     // timestamp of user creation
}

// RegisterRequest represents the expected payload for user registration
// This is what the client will send when registering a new user to POST /register
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"` // user's email address (must be valid)
	Password string `json:"password" validate:"required,min=6"` // user's password (must be at least 6 characters)
}

// LoginRequest represents the expected payload for user login
// This is what the client will send when logging in to POST /login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"` // user's email address (must be valid)
	Password string `json:"password" validate:"required"`     // user's password (required)
}

// AuthResponse represents the response sent back to the client after successful authentication
type AuthResponse struct {
	Token string `json:"token"` // JWT token that the client can use for authenticated requests
	User User   `json:"user"`  // user information (excluding password hash)
}