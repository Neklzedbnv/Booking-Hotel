package booking

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"Gofinal/internal/domain"
	"Gofinal/pkg/common"
)

// jsonError sends JSON error
func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user_id from JWT token (context)
	userID, ok := r.Context().Value(common.ContextUserID).(int)
	if !ok || userID == 0 {
		jsonError(w, "authentication required", http.StatusUnauthorized)
		return
	}

	var req struct {
		RoomID        int     `json:"room_id"`
		MealplanID    *int    `json:"mealplan_id"`
		PackageID     *int    `json:"package_id"`
		StartDate     string  `json:"start_date"`
		EndDate       string  `json:"end_date"`
		TotalPrice    float64 `json:"total_price"`
		PaymentMethod string  `json:"payment_method"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid json", http.StatusBadRequest)
		return
	}

	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		jsonError(w, "invalid start_date", http.StatusBadRequest)
		return
	}

	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		jsonError(w, "invalid end_date", http.StatusBadRequest)
		return
	}

	b := domain.Booking{
		UserID:     userID,
		RoomID:     req.RoomID,
		MealplanID: req.MealplanID,
		PackageID:  req.PackageID,
		StartDate:  start,
		EndDate:    end,
		TotalPrice: req.TotalPrice,
	}

	log.Printf("Creating booking: user_id=%d, room_id=%d, payment_method=%s", b.UserID, b.RoomID, req.PaymentMethod)

	res, err := h.service.Create(b, req.PaymentMethod)
	if err != nil {
		log.Printf("booking create error: %v", err)
		jsonError(w, "error creating booking", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user_id from JWT token
	userID, ok := r.Context().Value(common.ContextUserID).(int)
	if !ok || userID == 0 {
		jsonError(w, "authentication required", http.StatusUnauthorized)
		return
	}

	// Get only current user's bookings
	data, err := h.service.GetAllByUser(userID)
	if err != nil {
		jsonError(w, "error fetching bookings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	data, err := h.service.GetByID(id)
	if err != nil {
		jsonError(w, "booking not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID         int     `json:"id"`
		StartDate  string  `json:"start_date"`
		EndDate    string  `json:"end_date"`
		TotalPrice float64 `json:"total_price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid json", http.StatusBadRequest)
		return
	}

	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		jsonError(w, "invalid start_date", http.StatusBadRequest)
		return
	}

	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		jsonError(w, "invalid end_date", http.StatusBadRequest)
		return
	}

	b := domain.Booking{
		ID:         req.ID,
		StartDate:  start,
		EndDate:    end,
		TotalPrice: req.TotalPrice,
	}

	res, err := h.service.Update(b)
	if err != nil {
		jsonError(w, "error updating booking", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user_id from JWT token
	userID, ok := r.Context().Value(common.ContextUserID).(int)
	if !ok || userID == 0 {
		jsonError(w, "authentication required", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Use CancelBooking to verify owner and release room
	if err := h.service.CancelBooking(id, userID); err != nil {
		jsonError(w, "error deleting booking", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "booking deleted"})
}
