package catalog

import (
	"Gofinal/internal/domain"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type RoomHandler struct {
	service *RoomService
}

func NewRoomHandler(s *RoomService) *RoomHandler {
	return &RoomHandler{service: s}
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Code     string  `json:"code"`
		TypeID   int64   `json:"type_id"`
		Capacity int     `json:"capacity"`
		Price    float64 `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	room := domain.Room{
		Code:      req.Code,
		TypeID:    req.TypeID,
		Capacity:  req.Capacity,
		Price:     req.Price,
		Status:    "available",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := h.service.CreateRoom(room)
	if err != nil {
		http.Error(w, "error creating room", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *RoomHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomID := r.URL.Query().Get("id")
	if roomID == "" {
		http.Error(w, "missing room id", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(roomID, 10, 64)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	room, err := h.service.GetRoomByID(id)
	if err != nil {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

func (h *RoomHandler) ListRooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := r.URL.Query().Get("status")
	typeID := r.URL.Query().Get("type_id")

	rooms, err := h.service.ListRooms(status, typeID)
	if err != nil {
		http.Error(w, "error fetching rooms", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

func (h *RoomHandler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomID := r.URL.Query().Get("id")
	if roomID == "" {
		http.Error(w, "missing room id", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(roomID, 10, 64)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	var req struct {
		Price  *float64 `json:"price"`
		Status *string  `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	updated, err := h.service.UpdateRoom(id, req.Price, req.Status)
	if err != nil {
		http.Error(w, "error updating room", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (h *RoomHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomID := r.URL.Query().Get("id")
	if roomID == "" {
		http.Error(w, "missing room id", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(roomID, 10, 64)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteRoom(id)
	if err != nil {
		http.Error(w, "error deleting room", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "room deleted successfully"})
}

func (h *RoomHandler) CreateRoomType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name      string  `json:"name"`
		Capacity  int     `json:"capacity"`
		BasePrice float64 `json:"base_price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	roomType := domain.RoomType{
		Name:      req.Name,
		Capacity:  req.Capacity,
		BasePrice: req.BasePrice,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := h.service.CreateRoomType(roomType)
	if err != nil {
		http.Error(w, "error creating room type", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *RoomHandler) ListRoomTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	types, err := h.service.ListRoomTypes()
	if err != nil {
		http.Error(w, "error fetching room types", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types)
}

func (h *RoomHandler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CheckInDate  string `json:"check_in_date"`
		CheckOutDate string `json:"check_out_date"`
		Capacity     *int   `json:"capacity"`
		TypeID       *int64 `json:"type_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	checkIn, err := time.Parse("2006-01-02", req.CheckInDate)
	if err != nil {
		http.Error(w, "invalid check_in_date format", http.StatusBadRequest)
		return
	}

	checkOut, err := time.Parse("2006-01-02", req.CheckOutDate)
	if err != nil {
		http.Error(w, "invalid check_out_date format", http.StatusBadRequest)
		return
	}

	available, err := h.service.CheckAvailability(checkIn, checkOut, req.Capacity, req.TypeID)
	if err != nil {
		http.Error(w, "error checking availability", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(available)
}
