package main

import (
	"log"
	"net/http"
	"time"

	"Gofinal/internal/booking"
	"Gofinal/internal/db"
	httpx "Gofinal/internal/http"
)

func main() {
	dbConn := db.NewPostgres()

	repo := booking.NewRepo(dbConn)
	service := booking.NewService(repo)
	handler := booking.NewHandler(service)

	router := httpx.NewRouter(handler)

	go func() {
		for {
			log.Println("background worker alive")
			time.Sleep(10 * time.Second)
		}
	}()

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
