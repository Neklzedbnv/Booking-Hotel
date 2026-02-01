package main

import (
	"log"
	"net/http"

	"Gofinal/internal/booking"
	"Gofinal/internal/db"
	httpRouter "Gofinal/internal/http"
)

func main() {
	database := db.NewPostgres()

	bookingRepo := booking.NewRepo(database)
	bookingService := booking.NewService(bookingRepo)
	bookingHandler := booking.NewHandler(bookingService)

	router := httpRouter.NewRouter(bookingHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}