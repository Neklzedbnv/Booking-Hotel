package admin

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// GetDashboardStats returns dashboard statistics
func (h *Handler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.service.GetDashboardStats()
	if err != nil {
		http.Error(w, "error getting stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// ListUsers returns list of all users
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.service.ListUsers()
	if err != nil {
		http.Error(w, "error listing users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// SetupAdmin sets admin role by email (one-time setup)
func (h *Handler) SetupAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "email required", http.StatusBadRequest)
		return
	}

	if err := h.service.SetAdminByEmail(req.Email); err != nil {
		http.Error(w, "error setting admin", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "admin role set for " + req.Email})
}

// ResetPassword resets user password
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	if err := h.service.ResetPasswordByEmail(req.Email, string(hashedPassword)); err != nil {
		http.Error(w, "error resetting password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "password reset for " + req.Email})
}

// UpdateUserRole updates user role
func (h *Handler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID int    `json:"user_id"`
		Role   string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Role != "user" && req.Role != "admin" {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateUserRole(req.UserID, req.Role); err != nil {
		http.Error(w, "error updating role", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "role updated"})
}

// BlockUser blocks/unblocks user
func (h *Handler) BlockUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID  int  `json:"user_id"`
		Blocked bool `json:"blocked"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.service.BlockUser(req.UserID, req.Blocked); err != nil {
		http.Error(w, "error blocking user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "user block status updated"})
}

// GetBookings returns all bookings with details
func (h *Handler) GetBookings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bookings, err := h.service.GetBookingsWithDetails()
	if err != nil {
		http.Error(w, "error getting bookings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

// UpdateBookingStatus updates booking status
func (h *Handler) UpdateBookingStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		BookingID int    `json:"booking_id"`
		Status    string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	validStatuses := map[string]bool{
		"pending":   true,
		"confirmed": true,
		"cancelled": true,
		"completed": true,
	}

	if !validStatuses[req.Status] {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateBookingStatus(req.BookingID, req.Status); err != nil {
		http.Error(w, "error updating status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "booking status updated"})
}

// UpdateRoomType updates room type
func (h *Handler) UpdateRoomType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID        int64   `json:"id"`
		Name      string  `json:"name"`
		Capacity  int     `json:"capacity"`
		BasePrice float64 `json:"base_price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateRoomType(req.ID, req.Name, req.Capacity, req.BasePrice); err != nil {
		http.Error(w, "error updating room type", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "room type updated"})
}

// DeleteRoomType deletes room type
func (h *Handler) DeleteRoomType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteRoomType(id); err != nil {
		http.Error(w, "error deleting room type", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "room type deleted"})
}

// UploadImage uploads image for room
func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Maximum 10MB
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get image type from form
	imageType := r.FormValue("type")
	if imageType == "" {
		imageType = "room"
	}

	// Generate filename
	ext := filepath.Ext(handler.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	var filename string
	switch imageType {
	case "hero":
		filename = "hero" + ext
	case "spa":
		filename = "spa" + ext
	case "gym":
		filename = "gym" + ext
	case "bike":
		filename = "bike" + ext
	case "tour":
		filename = "tour" + ext
	default:
		// For rooms: room-standard, room-deluxe, etc.
		roomType := r.FormValue("room_type")
		if roomType == "" {
			roomType = "custom"
		}
		filename = "room-" + strings.ToLower(roomType) + ext
	}

	// Path for saving
	destPath := filepath.Join("public", "images", filename)

	// Create file
	dst, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy contents
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "error saving file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "image uploaded",
		"filename": filename,
		"path":     "/public/images/" + filename,
	})
}

// ListImages returns list of all images
func (h *Handler) ListImages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	imagesDir := filepath.Join("public", "images")
	files, err := os.ReadDir(imagesDir)
	if err != nil {
		http.Error(w, "error reading images directory", http.StatusInternalServerError)
		return
	}

	var images []map[string]string
	for _, f := range files {
		if !f.IsDir() {
			images = append(images, map[string]string{
				"name": f.Name(),
				"path": "/public/images/" + f.Name(),
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}

// DeleteImage deletes image
func (h *Handler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "filename required", http.StatusBadRequest)
		return
	}

	// Check that file is only in images directory (security)
	filename = filepath.Base(filename)
	path := filepath.Join("public", "images", filename)

	if err := os.Remove(path); err != nil {
		http.Error(w, "error deleting file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "image deleted"})
}
