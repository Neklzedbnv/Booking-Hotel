package common

import (
	"encoding/json"
	"net/http"
)

// ─── Context keys ────────────────────────────────────────────────────────────

type ContextKey string

const (
	ContextUserID    ContextKey = "userID"
	ContextUserRole  ContextKey = "userRole"
	ContextUserEmail ContextKey = "userEmail"
	ContextRequestID ContextKey = "requestID"
)

// JSONError returns an error in JSON format
func JSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// JSONSuccess returns a successful response in JSON format
func JSONSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
