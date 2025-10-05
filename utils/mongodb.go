package utils

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectIDFromString converts string to MongoDB ObjectID
func ObjectIDFromString(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

// IsValidObjectID checks if string is a valid MongoDB ObjectID
func IsValidObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

// CurrentTime returns current time
func CurrentTime() time.Time {
	return time.Now()
}

// FormatTime formats time to string
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
