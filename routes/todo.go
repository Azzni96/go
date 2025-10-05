package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nihad/todo/auth"
	"github.com/nihad/todo/handlers"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/mongo"
)

func TodoRoutes(db *mongo.Database, rnd *renderer.Render) http.Handler {
	todoHandler := handlers.NewTodoHandler(db, rnd)
	userService := auth.GetUserService(db)

	rg := chi.NewRouter()

	// Apply optional auth middleware - works with or without authentication
	rg.Use(auth.OptionalAuthMiddleware(userService))

	rg.Group(func(r chi.Router) {
		r.Get("/", todoHandler.FetchTodos)
		r.Post("/", todoHandler.CreateTodo)
		r.Put("/{id}", todoHandler.UpdateTodo)
		r.Delete("/{id}", todoHandler.DeleteTodo)
	})
	return rg
}
