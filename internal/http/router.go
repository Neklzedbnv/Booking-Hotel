package http

import (
	"net/http"

	"Gofinal/internal/booking"
)

func NewRouter(h *booking.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/bookings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Create(w, r)
		}
		if r.Method == http.MethodGet {
			h.GetAll(w, r)
		}
	})

	mux.HandleFunc("/booking", h.GetByID)

	return mux
}
