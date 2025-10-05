package auth

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthRoutes creates authentication routes
func AuthRoutes(db *mongo.Database, rnd *renderer.Render) http.Handler {
	userService := NewUserService(db)
	authHandler := NewAuthHandler(userService, rnd)

	r := chi.NewRouter()

	// Public routes (no authentication required)
	r.Group(func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// Protected routes (authentication required)
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(userService))
		r.Get("/profile", authHandler.Profile)
		r.Post("/refresh", authHandler.RefreshToken)
		r.Post("/logout", authHandler.Logout)
	})

	return r
}

// GetUserService returns user service instance (for use in other packages)
func GetUserService(db *mongo.Database) *UserService {
	return NewUserService(db)
}
