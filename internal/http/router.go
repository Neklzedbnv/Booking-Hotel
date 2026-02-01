package http

import (
	"net/http"

	"Gofinal/internal/booking"
	"Gofinal/internal/auth"
	"Gofinal/internal/catalog"
)

func NewRouter(
	bookingHandler *booking.Handler,
	authHandler *auth.Handler,
	catalogHandler *catalog.Handler,
) http.Handler {

	mux := http.NewServeMux()

	// booking
	mux.HandleFunc("/bookings", bookingHandler.GetAll)

	// auth
	mux.HandleFunc("/auth/health", authHandler.Health)

	// catalog
	mux.HandleFunc("/catalog", catalogHandler.GetAll)

	return mux
}
