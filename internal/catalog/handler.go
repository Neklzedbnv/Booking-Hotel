package catalog

import "net/http"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// read-only endpoint (достаточно для Assignment 4)
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"catalog endpoint works"}`))
}