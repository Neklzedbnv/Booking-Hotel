package booking

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"Gofinal/internal/domain"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID     int     `json:"user_id"`
		RoomID     int     `json:"room_id"`
		MealplanID *int    `json:"mealplan_id"`
		PackageID  *int    `json:"package_id"`
		StartDate  string  `json:"start_date"`
		EndDate    string  `json:"end_date"`
		TotalPrice float64 `json:"total_price"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	start, _ := time.Parse("2006-01-02", req.StartDate)
	end, _ := time.Parse("2006-01-02", req.EndDate)

	b := domain.Booking{
		UserID:     req.UserID,
		RoomID:     req.RoomID,
		MealplanID: req.MealplanID,
		PackageID:  req.PackageID,
		StartDate:  start,
		EndDate:    end,
		TotalPrice: req.TotalPrice,
	}

	res, err := h.service.Create(b)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	data, _ := h.service.GetAll()
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	data, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	json.NewEncoder(w).Encode(data)
}
