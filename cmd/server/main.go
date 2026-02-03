package main

import (
	"log"
	"net/http"
	"time"

	"Gofinal/internal/auth"
	"Gofinal/internal/booking"
	"Gofinal/internal/catalog"
	"Gofinal/internal/db"
	httpx "Gofinal/internal/http"
)

func main() {
	dbConn := db.NewPostgres()

	
	bookingRepo := booking.NewRepo(dbConn)
	bookingService := booking.NewService(bookingRepo)
	bookingHandler := booking.NewHandler(bookingService)

	
	authRepo := auth.NewRepo()
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)


	catalogRepo := catalog.NewRepo()
	catalogService := catalog.NewService(catalogRepo)
	catalogHandler := catalog.NewHandler(catalogService)

	router := httpx.NewRouter(
		bookingHandler,
		authHandler,
		catalogHandler,
	)

	
	go func() {
		for {
			log.Println("background worker alive")
			time.Sleep(10 * time.Second)
		}
	}()

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
