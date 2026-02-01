package catalog

import "net/http"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}


func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("catalog is not implemented yet"))
}
