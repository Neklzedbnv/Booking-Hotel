package main

import (
	"Gofinal/internal/admin"
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
	// ── Database connection ─────────────────────────────────────────
	dbConn := db.NewPostgres()

	// ── Migrations ─────────────────────────────────────────────────
	db.RunMigrations(dbConn)

	// ── Auth ─────────────────────────────────────────────────────
	authRepo := auth.NewRepo(dbConn)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	// ── Rooms ────────────────────────────────────────────────────
	roomRepo := catalog.NewRoomRepo(dbConn)
	roomService := catalog.NewRoomService(roomRepo)
	roomHandler := catalog.NewRoomHandler(roomService)

	// ── Packages ─────────────────────────────────────────────────
	packageRepo := catalog.NewPackageRepo(dbConn)
	packageService := catalog.NewPackageService(packageRepo)
	packageHandler := catalog.NewPackageHandler(packageService)

	// ── Services ─────────────────────────────────────────────────
	svcRepo := catalog.NewSvcRepo(dbConn)
	svcService := catalog.NewSvcService(svcRepo)
	svcHandler := catalog.NewSvcHandler(svcService)

	// ── MealPlans ────────────────────────────────────────────────
	mpRepo := catalog.NewMealPlanRepo(dbConn)
	mpService := catalog.NewMealPlanService(mpRepo)
	mpHandler := catalog.NewMealPlanHandler(mpService)

	// ── Payments ─────────────────────────────────────────────────
	paymentRepo := booking.NewPaymentRepo(dbConn)
	paymentService := booking.NewPaymentService(paymentRepo)
	paymentHandler := booking.NewPaymentHandler(paymentService)

	// ── Booking ──────────────────────────────────────────────────
	bookingRepo := booking.NewRepo(dbConn)
	bookingService := booking.NewService(bookingRepo, roomRepo, mpRepo, packageRepo, paymentRepo)
	bookingHandler := booking.NewHandler(bookingService)

	// ── Reviews ──────────────────────────────────────────────────
	reviewRepo := booking.NewReviewRepo(dbConn)
	reviewService := booking.NewReviewService(reviewRepo)
	reviewHandler := booking.NewReviewHandler(reviewService)

	// ── Admin ─────────────────────────────────────────────────────
	adminRepo := admin.NewRepo(dbConn)
	adminService := admin.NewService(adminRepo)
	adminHandler := admin.NewHandler(adminService)

	// ── Page Handler (Go templates) ────────────────────────────────
	pageHandler := introuter.NewPageHandler(
		"UI",
		roomService,
		packageService,
		svcService,
		mpService,
		bookingService,
		reviewService,
	)

	// ── Router + Middleware ───────────────────────────────────────
	router := introuter.NewRouter(
		authHandler,
		bookingHandler,
		roomHandler,
		packageHandler,
		svcHandler,
		mpHandler,
		reviewHandler,
		paymentHandler,
		adminHandler,
		pageHandler,
	)
	handler := router.SetupRoutes()

	// ── Background worker ────────────────────────────────────────
	go func() {
		for {
			log.Println("background worker alive")
			time.Sleep(30 * time.Second)
		}
	}()

	// ── Start server ────────────────────────────────────────────
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
