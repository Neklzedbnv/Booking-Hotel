package http

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"Gofinal/internal/booking"
	"Gofinal/internal/catalog"
	"Gofinal/internal/domain"
)

// PageHandler renders HTML pages via Go templates (html/template).
type PageHandler struct {
	tmpl            *template.Template
	roomService     *catalog.RoomService
	packageService  *catalog.PackageService
	svcService      *catalog.SvcService
	mealPlanService *catalog.MealPlanService
	bookingService  *booking.Service
	reviewService   *booking.ReviewService
}

// templateFuncMap — custom functions for templates.
var templateFuncMap = template.FuncMap{
	"json": func(v interface{}) template.JS {
		b, _ := json.Marshal(v)
		return template.JS(b)
	},
	"roomImg": func(typeID int64) string {
		switch typeID {
		case 1:
			return "/public/images/room-standard.jpg"
		case 2:
			return "/public/images/room-deluxe.jpg"
		case 3:
			return "/public/images/room-suite.jpg"
		case 4:
			return "/public/images/room-vip.jpg"
		case 5:
			return "/public/images/room-premium.jpg"
		case 6:
			return "/public/images/room-business.jpg"
		default:
			return "/public/images/room-standard.jpg"
		}
	},
	"statusColor": func(status string) string {
		if status == "available" {
			return "green"
		}
		return "#e74c3c"
	},
	"statusText": func(status string) string {
		if status == "available" {
			return "✓ Available"
		}
		return "✗ Booked"
	},
	"isBooked": func(status string) bool {
		return status != "available"
	},
	"mul": func(a, b interface{}) float64 {
		af, _ := toFloat64(a)
		bf, _ := toFloat64(b)
		return af * bf
	},
	"serviceImg": func(name string) string {
		switch name {
		case "Spa Treatments":
			return "/public/images/spa.jpg"
		case "Gym":
			return "/public/images/gym.jpg"
		case "Bike Rental":
			return "/public/images/bike.jpg"
		case "Tours":
			return "/public/images/tour.jpg"
		default:
			return "/public/images/spa.jpg"
		}
	},
}

func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

// NewPageHandler parses all templates from UI/ directory and returns PageHandler.
func NewPageHandler(
	templateDir string,
	roomSvc *catalog.RoomService,
	pkgSvc *catalog.PackageService,
	svcSvc *catalog.SvcService,
	mpSvc *catalog.MealPlanService,
	bookSvc *booking.Service,
	revSvc *booking.ReviewService,
) *PageHandler {
	pattern := filepath.Join(templateDir, "*.html")
	tmpl := template.Must(template.New("").Funcs(templateFuncMap).ParseGlob(pattern))

	return &PageHandler{
		tmpl:            tmpl,
		roomService:     roomSvc,
		packageService:  pkgSvc,
		svcService:      svcSvc,
		mealPlanService: mpSvc,
		bookingService:  bookSvc,
		reviewService:   revSvc,
	}
}

func (p *PageHandler) render(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := p.tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("template render error (%s): %v", name, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// populateRoomTypes populates the Type field for each room.
func (p *PageHandler) populateRoomTypes(rooms []domain.Room) []domain.Room {
	allTypes, _ := p.roomService.ListRoomTypes()
	tmap := make(map[int64]*domain.RoomType, len(allTypes))
	for i := range allTypes {
		tmap[allTypes[i].ID] = &allTypes[i]
	}
	for i := range rooms {
		rooms[i].Type = tmap[rooms[i].TypeID]
	}
	return rooms
}

// ─── Page Data structs ───────────────────────────────────────────────────────

type HomePageData struct {
	RoomTypes []domain.RoomType
	Services  []domain.Service
	Rooms     []domain.Room
	Reviews   []domain.Review
}

type RoomsPageData struct {
	Rooms     []domain.Room
	RoomTypes []domain.RoomType
}

type RoomDetailData struct {
	Room      *domain.Room
	RoomType  *domain.RoomType
	Services  []domain.Service
	MealPlans []domain.MealPlan
	Packages  []domain.Package
}

type BookingPageData struct {
	Rooms     []domain.Room
	RoomTypes []domain.RoomType
	MealPlans []domain.MealPlan
	Packages  []domain.Package
	Services  []domain.Service
}

type ProfilePageData struct {
	Bookings []domain.Booking
}

// ─── Page Handlers ───────────────────────────────────────────────────────────

func (p *PageHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	roomTypes, _ := p.roomService.ListRoomTypes()
	services, _ := p.svcService.List()
	rooms, _ := p.roomService.ListRooms("", "")
	rooms = p.populateRoomTypes(rooms)
	reviews, _ := p.reviewService.ListAll()

	p.render(w, "index.html", HomePageData{
		RoomTypes: roomTypes,
		Services:  services,
		Rooms:     rooms,
		Reviews:   reviews,
	})
}

func (p *PageHandler) RoomsPage(w http.ResponseWriter, r *http.Request) {
	rooms, _ := p.roomService.ListRooms("", "")
	rooms = p.populateRoomTypes(rooms)
	roomTypes, _ := p.roomService.ListRoomTypes()

	p.render(w, "rooms.html", RoomsPageData{
		Rooms:     rooms,
		RoomTypes: roomTypes,
	})
}

func (p *PageHandler) RoomDetailPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		id = 1
	}
	room, err := p.roomService.GetRoomByID(id)
	if err != nil {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}

	var roomType *domain.RoomType
	allTypes, _ := p.roomService.ListRoomTypes()
	for _, rt := range allTypes {
		if rt.ID == room.TypeID {
			cp := rt
			roomType = &cp
			room.Type = &cp
			break
		}
	}

	services, _ := p.svcService.List()
	mealPlans, _ := p.mealPlanService.List()
	packages, _ := p.packageService.ListPackages(true)

	p.render(w, "room-details.html", RoomDetailData{
		Room:      room,
		RoomType:  roomType,
		Services:  services,
		MealPlans: mealPlans,
		Packages:  packages,
	})
}

func (p *PageHandler) BookingPage(w http.ResponseWriter, r *http.Request) {
	rooms, _ := p.roomService.ListRooms("available", "")
	rooms = p.populateRoomTypes(rooms)
	roomTypes, _ := p.roomService.ListRoomTypes()
	mealPlans, _ := p.mealPlanService.List()
	packages, _ := p.packageService.ListPackages(true)
	services, _ := p.svcService.List()

	p.render(w, "booking.html", BookingPageData{
		Rooms:     rooms,
		RoomTypes: roomTypes,
		MealPlans: mealPlans,
		Packages:  packages,
		Services:  services,
	})
}

func (p *PageHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	p.render(w, "login.html", nil)
}

func (p *PageHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	p.render(w, "register.html", nil)
}

func (p *PageHandler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	bookings, _ := p.bookingService.GetAll()
	p.render(w, "profile.html", ProfilePageData{Bookings: bookings})
}

func (p *PageHandler) ContactPage(w http.ResponseWriter, r *http.Request) {
	p.render(w, "contact.html", nil)
}

func (p *PageHandler) AdminPage(w http.ResponseWriter, r *http.Request) {
	rooms, _ := p.roomService.ListRooms("", "")
	rooms = p.populateRoomTypes(rooms)
	roomTypes, _ := p.roomService.ListRoomTypes()
	services, _ := p.svcService.List()
	packages, _ := p.packageService.ListPackages(false)

	p.render(w, "admin.html", AdminPageData{
		Rooms:     rooms,
		RoomTypes: roomTypes,
		Services:  services,
		Packages:  packages,
	})
}

type AdminPageData struct {
	Rooms     []domain.Room
	RoomTypes []domain.RoomType
	Services  []domain.Service
	Packages  []domain.Package
}
