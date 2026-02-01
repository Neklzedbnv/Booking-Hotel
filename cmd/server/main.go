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
	// одно подключение к БД
	dbConn := db.NewPostgres()

	// booking layers
	repo := booking.NewRepo(dbConn)
	service := booking.NewService(repo)
	handler := booking.NewHandler(service)

	// router
	router := httpx.NewRouter(handler)

	// goroutine (concurrency requirement)
	go func() {
		for {
			log.Println("background worker alive")
			time.Sleep(10 * time.Second)
		}
	}()

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
