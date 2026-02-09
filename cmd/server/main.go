package main

import (
	"Gofinal/internal/auth"
	"Gofinal/internal/booking"
	"Gofinal/internal/catalog"
	"Gofinal/internal/db"
	introuter "Gofinal/internal/http"
	"log"
	"net/http"
	"time"
)

func main() {
	
	dbConn := db.NewPostgres()

	
	authRepo := auth.NewRepo(dbConn)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	
	bookingRepo := booking.NewRepo(dbConn)
	bookingService := booking.NewService(bookingRepo)
	bookingHandler := booking.NewHandler(bookingService)

	
	roomRepo := catalog.NewRoomRepo(dbConn)
	roomService := catalog.NewRoomService(roomRepo)
	roomHandler := catalog.NewRoomHandler(roomService)

	packageRepo := catalog.NewPackageRepo(dbConn)
	packageService := catalog.NewPackageService(packageRepo)
	packageHandler := catalog.NewPackageHandler(packageService)

	
	router := introuter.NewRouter(authHandler, bookingHandler)
	router.SetRoomHandler(roomHandler)
	router.SetPackageHandler(packageHandler)
	handler := router.SetupRoutes()

	
	go func() {
		for {
			log.Println("background worker alive")
			time.Sleep(10 * time.Second)
		}
	}()

	
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
