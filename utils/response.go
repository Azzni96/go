package utils

import (
	"net/http"

	"github.com/thedevsaddam/renderer"
)

// RespondWithError sends an error response in JSON format
func RespondWithError(rnd *renderer.Render, w http.ResponseWriter, code int, message string) {
	rnd.JSON(w, code, renderer.M{
		"error":   true,
		"message": message,
	})
}

// RespondWithSuccess sends a success response in JSON format
func RespondWithSuccess(rnd *renderer.Render, w http.ResponseWriter, code int, message string, data interface{}) {
	response := renderer.M{
		"error":   false,
		"message": message,
	}

	if data != nil {
		response["data"] = data
	}

	rnd.JSON(w, code, response)
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(rnd *renderer.Render, w http.ResponseWriter, code int, payload interface{}) {
	rnd.JSON(w, code, payload)
}
