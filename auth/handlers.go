package auth

import (
	"encoding/json"
	"net/http"

	"github.com/nihad/todo/utils"
	"github.com/thedevsaddam/renderer"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	userService *UserService
	rnd         *renderer.Render
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService *UserService, rnd *renderer.Render) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		rnd:         rnd,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	user, err := h.userService.Register(req)
	if err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusBadRequest, err.Error())
		return
	}

	// Generate token for the new user
	token, err := GenerateJWT(user.ID.Hex(), user.Username, user.Email)
	if err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := AuthResponse{
		Token: token,
		User:  user.ToUserResponse(),
	}

	utils.RespondWithSuccess(h.rnd, w, http.StatusCreated, "User registered successfully", response)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	authResponse, err := h.userService.Login(req)
	if err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.RespondWithSuccess(h.rnd, w, http.StatusOK, "Login successful", authResponse)
}

// Profile returns current user profile
func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	user, err := RequireUserFromContext(r)
	if err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusUnauthorized, "User not found in context")
		return
	}

	utils.RespondWithSuccess(h.rnd, w, http.StatusOK, "Profile retrieved successfully", user.ToUserResponse())
}

// RefreshToken generates a new token
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	user, err := RequireUserFromContext(r)
	if err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// Generate new token
	newToken, err := GenerateJWT(user.ID.Hex(), user.Username, user.Email)
	if err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := map[string]string{
		"token": newToken,
	}

	utils.RespondWithSuccess(h.rnd, w, http.StatusOK, "Token refreshed successfully", response)
}

// Logout handles user logout (client-side token removal)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Since we're using JWT (stateless), logout is handled client-side
	// The client should remove the token from storage
	// We can optionally add token blacklisting here if needed

	utils.RespondWithSuccess(h.rnd, w, http.StatusOK, "Logout successful", nil)
}
