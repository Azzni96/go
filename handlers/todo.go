package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/nihad/todo/models"
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
		h.rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	// simple validation
	if t.Title == "" {
		h.rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is required",
		})
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
		h.rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to save todo",
			"error":   err,
		})
		return
	}

	h.rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		"todo_id": tm.ID.Hex(),
	})
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	var t models.Todo

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	// simple validation
	if t.Title == "" {
		h.rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// if input is okay, update a todo
	_, err = h.db.Collection(collectionName).UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"title": t.Title, "completed": t.Completed}},
	)
	if err != nil {
		h.rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to update todo",
			"error":   err,
		})
		return
	}

	h.rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo updated successfully",
	})
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
	page := 1
	if p := query.Get("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	pageSize := 10
	if ps := query.Get("page_size"); ps != "" {
		if size, err := strconv.Atoi(ps); err == nil && size > 0 && size <= 100 {
			pageSize = size
		}
	}

	skip := (page - 1) * pageSize

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
		h.rnd.JSON(w, http.StatusInternalServerError, renderer.M{
			"message": "Failed to count todos",
			"error":   err,
		})
		return
	}

	// Find options with pagination and sorting
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(sortOptions)

	cursor, err := h.db.Collection(collectionName).Find(ctx, filter, findOptions)
	if err != nil {
		h.rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todos",
			"error":   err,
		})
		return
	}
	defer cursor.Close(ctx)

	var todos []models.TodoModel
	if err = cursor.All(ctx, &todos); err != nil {
		h.rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to decode todos",
			"error":   err,
		})
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

	// Calculate pagination info
	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	h.rnd.JSON(w, http.StatusOK, renderer.M{
		"data": todoList,
		"pagination": renderer.M{
			"current_page": page,
			"page_size":    pageSize,
			"total_count":  totalCount,
			"total_pages":  totalPages,
		},
		"filters": renderer.M{
			"completed": query.Get("completed"),
			"search":    query.Get("search"),
			"sort_by":   sortBy,
			"order":     order,
		},
	})
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = h.db.Collection(collectionName).DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		h.rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to delete todo",
			"error":   err,
		})
		return
	}

	h.rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo deleted successfully",
	})
}
