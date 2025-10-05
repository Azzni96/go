package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Completed bool               `bson:"completed"`
	CreatedAt time.Time          `bson:"createAt"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty"` // New field for user association
}

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UserID    string    `json:"user_id,omitempty"` // Optional in JSON response
}
