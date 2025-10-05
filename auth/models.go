package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"` // Never include in JSON
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
}

// UserResponse represents user data for API responses (without sensitive info)
type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

// LoginRequest represents login request data
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest represents registration request data
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ToUserResponse converts User to UserResponse (removes sensitive data)
func (u *User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:        u.ID.Hex(),
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		IsActive:  u.IsActive,
	}
}
