package http

import (
	"net/http"

	"Gofinal/internal/auth"
	"Gofinal/internal/booking"
	"Gofinal/internal/catalog"
)

type Router struct {
	authHandler    *auth.Handler
	bookingHandler *booking.Handler
	roomHandler    *catalog.RoomHandler
	packageHandler *catalog.PackageHandler
}

func NewRouter(authHandler *auth.Handler, bookingHandler *booking.Handler) *Router {
	return &Router{
		authHandler:    authHandler,
		bookingHandler: bookingHandler,
	}
}

func (r *Router) SetRoomHandler(rh *catalog.RoomHandler) {
	r.roomHandler = rh
}

func (r *Router) SetPackageHandler(ph *catalog.PackageHandler) {
	r.packageHandler = ph
}

func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	
	mux.HandleFunc("/auth/register", r.authHandler.Register)
	mux.HandleFunc("/auth/login", r.authHandler.Login)

	
	mux.HandleFunc("/booking", r.bookingHandler.Create)
	mux.HandleFunc("/booking/all", r.bookingHandler.GetAll)
	mux.HandleFunc("/booking/get", r.bookingHandler.GetByID)
	mux.HandleFunc("/booking/update", r.bookingHandler.Update)
	mux.HandleFunc("/booking/delete", r.bookingHandler.Delete)

	
	if r.roomHandler != nil {
		mux.HandleFunc("/rooms/create", r.roomHandler.CreateRoom)
		mux.HandleFunc("/rooms/get", r.roomHandler.GetRoom)
		mux.HandleFunc("/rooms/list", r.roomHandler.ListRooms)
		mux.HandleFunc("/rooms/update", r.roomHandler.UpdateRoom)
		mux.HandleFunc("/rooms/delete", r.roomHandler.DeleteRoom)
		mux.HandleFunc("/rooms/types/create", r.roomHandler.CreateRoomType)
		mux.HandleFunc("/rooms/types/list", r.roomHandler.ListRoomTypes)
		mux.HandleFunc("/rooms/availability", r.roomHandler.CheckAvailability)
	}

	
	if r.packageHandler != nil {
		mux.HandleFunc("/packages/create", r.packageHandler.CreatePackage)
		mux.HandleFunc("/packages/get", r.packageHandler.GetPackage)
		mux.HandleFunc("/packages/list", r.packageHandler.ListPackages)
		mux.HandleFunc("/packages/update", r.packageHandler.UpdatePackage)
		mux.HandleFunc("/packages/delete", r.packageHandler.DeletePackage)
		mux.HandleFunc("/packages/attach", r.packageHandler.AttachPackageToRoom)
		mux.HandleFunc("/packages/room", r.packageHandler.GetRoomPackages)
		mux.HandleFunc("/packages/detach", r.packageHandler.DetachPackageFromRoom)
	}

	return mux
}
