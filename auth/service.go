package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/nihad/todo/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersCollection = "users"

// UserService handles user-related operations
type UserService struct {
	db *mongo.Database
}

// NewUserService creates a new user service
func NewUserService(db *mongo.Database) *UserService {
	return &UserService{db: db}
}

// Register creates a new user
func (s *UserService) Register(req RegisterRequest) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Validate input
	if !utils.IsValidString(req.Username) {
		return nil, errors.New("username is required")
	}
	if !utils.IsValidString(req.Email) {
		return nil, errors.New("email is required")
	}
	if !IsValidPassword(req.Password) {
		return nil, errors.New("password must be at least 6 characters")
	}

	// Normalize username and email
	req.Username = utils.TrimAndLower(req.Username)
	req.Email = utils.TrimAndLower(req.Email)

	// Check if user already exists
	existingUser, err := s.findByUsernameOrEmail(req.Username, req.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this username or email already exists")
	}

	// Hash password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		ID:        primitive.NewObjectID(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: utils.CurrentTime(),
		UpdatedAt: utils.CurrentTime(),
		IsActive:  true,
	}

	// Insert to database
	_, err = s.db.Collection(usersCollection).InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a token
func (s *UserService) Login(req LoginRequest) (*AuthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Validate input
	if !utils.IsValidString(req.Username) {
		return nil, errors.New("username is required")
	}
	if !utils.IsValidString(req.Password) {
		return nil, errors.New("password is required")
	}

	// Normalize username
	username := utils.TrimAndLower(req.Username)

	// Find user
	var user User
	err := s.db.Collection(usersCollection).FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": username},
			{"email": username}, // Allow login with email too
		},
		"is_active": true,
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	// Check password
	if !CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := GenerateJWT(user.ID.Hex(), user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user.ToUserResponse(),
	}, nil
}

// GetByID finds user by ID
func (s *UserService) GetByID(userID string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if !utils.IsValidObjectID(userID) {
		return nil, errors.New("invalid user ID")
	}

	objID, _ := utils.ObjectIDFromString(userID)

	var user User
	err := s.db.Collection(usersCollection).FindOne(ctx, bson.M{
		"_id":       objID,
		"is_active": true,
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// findByUsernameOrEmail is a helper function to check if user exists
func (s *UserService) findByUsernameOrEmail(username, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := s.db.Collection(usersCollection).FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": strings.ToLower(username)},
			{"email": strings.ToLower(email)},
		},
	}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
