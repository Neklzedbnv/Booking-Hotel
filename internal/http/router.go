package http

import (
	"net/http"

	"Gofinal/internal/auth"
	"Gofinal/internal/booking"
	"Gofinal/internal/catalog"
)

func NewRouter(
	bookingHandler *booking.Handler,
	authHandler *auth.Handler,
	catalogHandler *catalog.Handler,
) http.Handler {

	mux := http.NewServeMux()

	
	mux.HandleFunc("/bookings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bookingHandler.GetAll(w, r)
		case http.MethodPost:
			bookingHandler.Create(w, r)
		case http.MethodPut:
			bookingHandler.Update(w, r)
		case http.MethodDelete:
			bookingHandler.Delete(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	
	mux.HandleFunc("/auth/health", authHandler.Health)

	
	mux.HandleFunc("/catalog", catalogHandler.GetAll)

	
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	return mux
}
