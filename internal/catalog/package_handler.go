package catalog

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"Gofinal/internal/domain"
)

type PackageHandler struct {
	service *PackageService
}

func NewPackageHandler(s *PackageService) *PackageHandler {
	return &PackageHandler{service: s}
}

// CreatePackage creates new service package
func (h *PackageHandler) CreatePackage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name          string  `json:"name"`
		Description   string  `json:"description"`
		PriceModifier float64 `json:"price_modifier"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	pkg := domain.Package{
		Name:          req.Name,
		Description:   req.Description,
		PriceModifier: req.PriceModifier,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	created, err := h.service.CreatePackage(pkg)
	if err != nil {
		http.Error(w, "error creating package", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *PackageHandler) GetPackage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	packageID := r.URL.Query().Get("id")
	if packageID == "" {
		http.Error(w, "missing package id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(packageID)
	if err != nil {
		http.Error(w, "invalid package id", http.StatusBadRequest)
		return
	}

	pkg, err := h.service.GetPackageByID(id)
	if err != nil {
		http.Error(w, "package not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pkg)
}

func (h *PackageHandler) ListPackages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	onlyActive := r.URL.Query().Get("active")

	pkgs, err := h.service.ListPackages(onlyActive == "true")
	if err != nil {
		http.Error(w, "error fetching packages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pkgs)
}

func (h *PackageHandler) UpdatePackage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	packageID := r.URL.Query().Get("id")
	if packageID == "" {
		http.Error(w, "missing package id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(packageID)
	if err != nil {
		http.Error(w, "invalid package id", http.StatusBadRequest)
		return
	}

	var req struct {
		Name          *string  `json:"name"`
		Description   *string  `json:"description"`
		PriceModifier *float64 `json:"price_modifier"`
		IsActive      *bool    `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	updated, err := h.service.UpdatePackage(id, req.Name, req.Description, req.PriceModifier, req.IsActive)
	if err != nil {
		http.Error(w, "error updating package", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (h *PackageHandler) DeletePackage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	packageID := r.URL.Query().Get("id")
	if packageID == "" {
		http.Error(w, "missing package id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(packageID)
	if err != nil {
		http.Error(w, "invalid package id", http.StatusBadRequest)
		return
	}

	err = h.service.DeletePackage(id)
	if err != nil {
		http.Error(w, "error deleting package", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "package deleted successfully"})
}

func (h *PackageHandler) AttachPackageToRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RoomID    int64 `json:"room_id"`
		PackageID int   `json:"package_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err := h.service.AttachPackageToRoom(req.RoomID, req.PackageID)
	if err != nil {
		http.Error(w, "error attaching package to room", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "package attached to room successfully"})
}

func (h *PackageHandler) GetRoomPackages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		http.Error(w, "missing room id", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(roomID, 10, 64)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	packages, err := h.service.GetRoomPackages(id)
	if err != nil {
		http.Error(w, "error fetching room packages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(packages)
}

func (h *PackageHandler) DetachPackageFromRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomID := r.URL.Query().Get("room_id")
	packageID := r.URL.Query().Get("package_id")

	if roomID == "" || packageID == "" {
		http.Error(w, "missing room_id or package_id", http.StatusBadRequest)
		return
	}

	rID, err := strconv.ParseInt(roomID, 10, 64)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	pID, err := strconv.Atoi(packageID)
	if err != nil {
		http.Error(w, "invalid package id", http.StatusBadRequest)
		return
	}

	err = h.service.DetachPackageFromRoom(rID, pID)
	if err != nil {
		http.Error(w, "error detaching package from room", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "package detached from room successfully"})
}
