package http

import (
	"net/http"

	"Gofinal/internal/admin"
	"Gofinal/internal/auth"
	"Gofinal/internal/booking"
	"Gofinal/internal/catalog"
)

// Router combines all HTTP handlers of the application.
type Router struct {
	authHandler     *auth.Handler
	bookingHandler  *booking.Handler
	roomHandler     *catalog.RoomHandler
	packageHandler  *catalog.PackageHandler
	svcHandler      *catalog.SvcHandler
	mealPlanHandler *catalog.MealPlanHandler
	reviewHandler   *booking.ReviewHandler
	paymentHandler  *booking.PaymentHandler
	adminHandler    *admin.Handler
	pageHandler     *PageHandler
}

// NewRouter creates a router with required handlers.
func NewRouter(
	authH *auth.Handler,
	bookingH *booking.Handler,
	roomH *catalog.RoomHandler,
	packageH *catalog.PackageHandler,
	svcH *catalog.SvcHandler,
	mpH *catalog.MealPlanHandler,
	reviewH *booking.ReviewHandler,
	paymentH *booking.PaymentHandler,
	adminH *admin.Handler,
	pageH *PageHandler,
) *Router {
	return &Router{
		authHandler:     authH,
		bookingHandler:  bookingH,
		roomHandler:     roomH,
		packageHandler:  packageH,
		svcHandler:      svcH,
		mealPlanHandler: mpH,
		reviewHandler:   reviewH,
		paymentHandler:  paymentH,
		adminHandler:    adminH,
		pageHandler:     pageH,
	}
}

// SetupRoutes registers all routes and applies middleware.
func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// ── Static files ───────────────────────────────────────────
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/public/", http.StripPrefix("/public/", fs))

	// ── Pages (Go templates) ───────────────────────────────────
	mux.HandleFunc("/", r.pageHandler.HomePage)
	mux.HandleFunc("/rooms", r.pageHandler.RoomsPage)
	mux.HandleFunc("/room-details", r.pageHandler.RoomDetailPage)
	mux.HandleFunc("/booking", r.pageHandler.BookingPage)
	mux.HandleFunc("/login", r.pageHandler.LoginPage)
	mux.HandleFunc("/register", r.pageHandler.RegisterPage)
	mux.HandleFunc("/profile", r.pageHandler.ProfilePage)
	mux.HandleFunc("/contact", r.pageHandler.ContactPage)
	// Admin panel - authorization check happens on the client (token in localStorage)
	mux.HandleFunc("/admin", r.pageHandler.AdminPage)

	// ── API: Auth ─────────────────────────────────────────────────
	mux.HandleFunc("/api/auth/register", r.authHandler.Register)
	mux.HandleFunc("/api/auth/login", r.authHandler.Login)

	// ── API: Bookings (CRUD) ──────────────────────────────────────
	mux.HandleFunc("/api/booking", r.bookingHandler.Create)
	mux.HandleFunc("/api/booking/all", r.bookingHandler.GetAll)
	mux.HandleFunc("/api/booking/get", r.bookingHandler.GetByID)
	mux.HandleFunc("/api/booking/update", r.bookingHandler.Update)
	mux.HandleFunc("/api/booking/delete", r.bookingHandler.Delete)

	// ── API: Rooms (CRUD + extras) ────────────────────────────────
	mux.HandleFunc("/api/rooms/create", r.roomHandler.CreateRoom)
	mux.HandleFunc("/api/rooms/get", r.roomHandler.GetRoom)
	mux.HandleFunc("/api/rooms/list", r.roomHandler.ListRooms)
	mux.HandleFunc("/api/rooms/update", r.roomHandler.UpdateRoom)
	mux.HandleFunc("/api/rooms/delete", r.roomHandler.DeleteRoom)
	mux.HandleFunc("/api/rooms/types/create", r.roomHandler.CreateRoomType)
	mux.HandleFunc("/api/rooms/types/list", r.roomHandler.ListRoomTypes)
	mux.HandleFunc("/api/rooms/availability", r.roomHandler.CheckAvailability)

	// ── API: Packages (CRUD + extras) ─────────────────────────────
	mux.HandleFunc("/api/packages/create", r.packageHandler.CreatePackage)
	mux.HandleFunc("/api/packages/get", r.packageHandler.GetPackage)
	mux.HandleFunc("/api/packages/list", r.packageHandler.ListPackages)
	mux.HandleFunc("/api/packages/update", r.packageHandler.UpdatePackage)
	mux.HandleFunc("/api/packages/delete", r.packageHandler.DeletePackage)
	mux.HandleFunc("/api/packages/attach", r.packageHandler.AttachPackageToRoom)
	mux.HandleFunc("/api/packages/room", r.packageHandler.GetRoomPackages)
	mux.HandleFunc("/api/packages/detach", r.packageHandler.DetachPackageFromRoom)

	// ── API: Services (CRUD) ──────────────────────────────────────
	mux.HandleFunc("/api/services/create", r.svcHandler.CreateService)
	mux.HandleFunc("/api/services/get", r.svcHandler.GetService)
	mux.HandleFunc("/api/services/list", r.svcHandler.ListServices)
	mux.HandleFunc("/api/services/update", r.svcHandler.UpdateService)
	mux.HandleFunc("/api/services/delete", r.svcHandler.DeleteService)

	// ── API: MealPlans (CRUD) ─────────────────────────────────────
	mux.HandleFunc("/api/mealplans/create", r.mealPlanHandler.CreateMealPlan)
	mux.HandleFunc("/api/mealplans/get", r.mealPlanHandler.GetMealPlan)
	mux.HandleFunc("/api/mealplans/list", r.mealPlanHandler.ListMealPlans)
	mux.HandleFunc("/api/mealplans/update", r.mealPlanHandler.UpdateMealPlan)
	mux.HandleFunc("/api/mealplans/delete", r.mealPlanHandler.DeleteMealPlan)

	// ── API: Reviews (CRUD) ───────────────────────────────────────
	mux.HandleFunc("/api/reviews/create", r.reviewHandler.CreateReview)
	mux.HandleFunc("/api/reviews/get", r.reviewHandler.GetReview)
	mux.HandleFunc("/api/reviews/list", r.reviewHandler.ListReviews)
	mux.HandleFunc("/api/reviews/update", r.reviewHandler.UpdateReview)
	mux.HandleFunc("/api/reviews/delete", r.reviewHandler.DeleteReview)

	// ── API: Payments (CRUD) ──────────────────────────────────────
	mux.HandleFunc("/api/payments/create", r.paymentHandler.CreatePayment)
	mux.HandleFunc("/api/payments/get", r.paymentHandler.GetPayment)
	mux.HandleFunc("/api/payments/list", r.paymentHandler.ListPayments)
	mux.HandleFunc("/api/payments/update", r.paymentHandler.UpdatePaymentStatus)
	mux.HandleFunc("/api/payments/delete", r.paymentHandler.DeletePayment)

	// ── API: Admin (only for abzalbahktiarow2006@gmail.com) ─────
	mux.Handle("/api/admin/stats", WrapFunc(r.adminHandler.GetDashboardStats, RequireAdminAPI))
	mux.Handle("/api/admin/users", WrapFunc(r.adminHandler.ListUsers, RequireAdminAPI))
	mux.Handle("/api/admin/users/role", WrapFunc(r.adminHandler.UpdateUserRole, RequireAdminAPI))
	mux.Handle("/api/admin/users/block", WrapFunc(r.adminHandler.BlockUser, RequireAdminAPI))
	mux.Handle("/api/admin/bookings", WrapFunc(r.adminHandler.GetBookings, RequireAdminAPI))
	mux.Handle("/api/admin/bookings/status", WrapFunc(r.adminHandler.UpdateBookingStatus, RequireAdminAPI))
	mux.Handle("/api/admin/room-types/update", WrapFunc(r.adminHandler.UpdateRoomType, RequireAdminAPI))
	mux.Handle("/api/admin/room-types/delete", WrapFunc(r.adminHandler.DeleteRoomType, RequireAdminAPI))
	mux.Handle("/api/admin/images", WrapFunc(r.adminHandler.ListImages, RequireAdminAPI))
	mux.Handle("/api/admin/images/upload", WrapFunc(r.adminHandler.UploadImage, RequireAdminAPI))
	mux.Handle("/api/admin/images/delete", WrapFunc(r.adminHandler.DeleteImage, RequireAdminAPI))
	// These endpoints are available without authorization (for initial setup)
	mux.HandleFunc("/api/admin/setup", r.adminHandler.SetupAdmin)
	mux.HandleFunc("/api/admin/reset-password", r.adminHandler.ResetPassword)

	// ── Middleware chain ──────────────────────────────────────────
	// Recovery → RequestID → Logging → CORS → Authenticate
	handler := Chain(mux,
		Recovery,
		RequestID,
		Logging,
		CORS,
		Authenticate,
		RateLimiter(RateLimiterConfig{RequestsPerSecond: 10, Burst: 30}),
	)

	return handler
}
