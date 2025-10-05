package auth

import (
	"context"
	"net/http"
	"strings"
)

// ContextKey is used for storing user info in request context
type ContextKey string

const UserContextKey ContextKey = "user"

// AuthMiddleware validates JWT token and adds user info to request context
func AuthMiddleware(userService *UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": true, "message": "Authorization header is required"}`))
				return
			}

			// Check if header starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": true, "message": "Invalid authorization header format"}`))
				return
			}

			// Extract token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": true, "message": "Token is required"}`))
				return
			}

			// Validate token
			claims, err := ValidateJWT(tokenString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": true, "message": "Invalid or expired token"}`))
				return
			}

			// Get user from database to ensure they still exist and are active
			user, err := userService.GetByID(claims.UserID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": true, "message": "User not found or inactive"}`))
				return
			}

			// Add user to request context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			r = r.WithContext(ctx)

			// Continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// OptionalAuthMiddleware validates token if present, but doesn't require it
func OptionalAuthMiddleware(userService *UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				// Extract and validate token
				tokenString := strings.TrimPrefix(authHeader, "Bearer ")
				if tokenString != "" {
					claims, err := ValidateJWT(tokenString)
					if err == nil {
						// Get user from database
						user, err := userService.GetByID(claims.UserID)
						if err == nil {
							// Add user to request context
							ctx := context.WithValue(r.Context(), UserContextKey, user)
							r = r.WithContext(ctx)
						}
					}
				}
			}

			// Continue to next handler (with or without user in context)
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext extracts user from request context
func GetUserFromContext(r *http.Request) (*User, bool) {
	user, ok := r.Context().Value(UserContextKey).(*User)
	return user, ok
}

// RequireUserFromContext gets user from context or returns error
func RequireUserFromContext(r *http.Request) (*User, error) {
	user, ok := GetUserFromContext(r)
	if !ok {
		return nil, http.ErrNoCookie // This shouldn't happen if middleware is used correctly
	}
	return user, nil
}
