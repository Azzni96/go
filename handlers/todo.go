package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/nihad/todo/models"
	"github.com/nihad/todo/utils"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "todo"

type TodoHandler struct {
	db  *mongo.Database
	rnd *renderer.Render
}

func NewTodoHandler(db *mongo.Database, rnd *renderer.Render) *TodoHandler {
	return &TodoHandler{db: db, rnd: rnd}
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var t models.Todo

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// simple validation
	if !utils.IsValidString(t.Title) {
		utils.RespondWithError(h.rnd, w, http.StatusBadRequest, "The title field is required")
		return
	}

	// if input is okay, create a todo
	tm := models.TodoModel{
		ID:        primitive.NewObjectID(),
		Title:     t.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.db.Collection(collectionName).InsertOne(ctx, tm)
	if err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusInternalServerError, "Failed to save todo")
		return
	}

	h.rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		"todo_id": tm.ID.Hex(),
		"todo": models.Todo{
			ID:        tm.ID.Hex(),
			Title:     tm.Title,
			Completed: tm.Completed,
			CreatedAt: tm.CreatedAt,
		},
	})
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	if !utils.IsValidObjectID(id) {
		utils.RespondWithError(h.rnd, w, http.StatusBadRequest, "The id is invalid")
		return
	}

	objID, _ := utils.ObjectIDFromString(id)

	var t models.Todo

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// simple validation
	if !utils.IsValidString(t.Title) {
		utils.RespondWithError(h.rnd, w, http.StatusBadRequest, "The title field is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// if input is okay, update a todo
	_, err := h.db.Collection(collectionName).UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"title": t.Title, "completed": t.Completed}},
	)
	if err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusInternalServerError, "Failed to update todo")
		return
	}

	utils.RespondWithSuccess(h.rnd, w, http.StatusOK, "Todo updated successfully", nil)
}

func (h *TodoHandler) FetchTodos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse query parameters
	query := r.URL.Query()

	// 1. FILTERING - Build MongoDB filter
	filter := bson.M{}

	// Filter by completion status
	if completed := query.Get("completed"); completed != "" {
		if completed == "true" {
			filter["completed"] = true
		} else if completed == "false" {
			filter["completed"] = false
		}
	}

	// Filter by title search (case-insensitive)
	if search := query.Get("search"); search != "" {
		filter["title"] = bson.M{"$regex": search, "$options": "i"}
	}

	// 3. PAGINATION - Parse pagination parameters
	page := utils.StringToInt(query.Get("page"), 1)
	pageSize := utils.StringToInt(query.Get("page_size"), 10)
	paginationParams := utils.NewPaginationParams(page, pageSize)

	// 2. SORTING - Build sort options
	sortOptions := bson.D{}
	sortBy := query.Get("sort_by")
	order := query.Get("order")

	sortOrder := -1 // Default descending
	if order == "asc" {
		sortOrder = 1
	}

	switch sortBy {
	case "title":
		sortOptions = append(sortOptions, bson.E{Key: "title", Value: sortOrder})
	case "completed":
		sortOptions = append(sortOptions, bson.E{Key: "completed", Value: sortOrder})
	case "created_at":
		sortOptions = append(sortOptions, bson.E{Key: "createAt", Value: sortOrder})
	default:
		// Default sort by creation date (newest first)
		sortOptions = append(sortOptions, bson.E{Key: "createAt", Value: -1})
	}

	// Count total documents for pagination
	totalCount, err := h.db.Collection(collectionName).CountDocuments(ctx, filter)
	if err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusInternalServerError, "Failed to count todos")
		return
	}

	// Find options with pagination and sorting
	findOptions := options.Find().
		SetSkip(int64(paginationParams.Skip)).
		SetLimit(int64(paginationParams.PageSize)).
		SetSort(sortOptions)

	cursor, err := h.db.Collection(collectionName).Find(ctx, filter, findOptions)
	if err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusInternalServerError, "Failed to fetch todos")
		return
	}
	defer cursor.Close(ctx)

	var todos []models.TodoModel
	if err = cursor.All(ctx, &todos); err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusInternalServerError, "Failed to decode todos")
		return
	}

	todoList := []models.Todo{}
	for _, t := range todos {
		todoList = append(todoList, models.Todo{
			ID:        t.ID.Hex(),
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}

	// Create pagination response with filters
	paginationResponse := utils.NewPaginationResponse(todoList, paginationParams.Page, paginationParams.PageSize, totalCount)
	responseData := renderer.M{
		"data": paginationResponse.Data,
		"pagination": renderer.M{
			"current_page": paginationResponse.Page,
			"page_size":    paginationResponse.PageSize,
			"total_items":  paginationResponse.TotalItems,
			"total_pages":  paginationResponse.TotalPages,
		},
		"filters": renderer.M{
			"completed": query.Get("completed"),
			"search":    query.Get("search"),
			"sort_by":   sortBy,
			"order":     order,
		},
	}

	utils.RespondWithJSON(h.rnd, w, http.StatusOK, responseData)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	if !utils.IsValidObjectID(id) {
		utils.RespondWithError(h.rnd, w, http.StatusBadRequest, "The id is invalid")
		return
	}

	objID, _ := utils.ObjectIDFromString(id)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.db.Collection(collectionName).DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		utils.RespondWithError(h.rnd, w, http.StatusInternalServerError, "Failed to delete todo")
		return
	}

	utils.RespondWithSuccess(h.rnd, w, http.StatusOK, "Todo deleted successfully", nil)
}
