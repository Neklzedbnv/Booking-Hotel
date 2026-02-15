package catalog

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Gofinal/internal/domain"
)

// MealPlanHandler — HTTP CRUD handlers for meal plans
type MealPlanHandler struct {
	service *MealPlanService
}

func NewMealPlanHandler(s *MealPlanService) *MealPlanHandler {
	return &MealPlanHandler{service: s}
}

// CreateMealPlan POST /api/mealplans
func (h *MealPlanHandler) CreateMealPlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Name        string  `json:"name"`
		PricePerDay float64 `json:"price_per_day"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	m := domain.MealPlan{Name: req.Name, PricePerDay: req.PricePerDay}
	created, err := h.service.Create(m)
	if err != nil {
		http.Error(w, "error creating meal plan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetMealPlan GET /api/mealplans/get?id=
func (h *MealPlanHandler) GetMealPlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	m, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

// ListMealPlans GET /api/mealplans/list
func (h *MealPlanHandler) ListMealPlans(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	list, err := h.service.List()
	if err != nil {
		http.Error(w, "error listing meal plans", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// UpdateMealPlan PUT /api/mealplans/update?id=
func (h *MealPlanHandler) UpdateMealPlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req struct {
		Name        string  `json:"name"`
		PricePerDay float64 `json:"price_per_day"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	m := domain.MealPlan{ID: int64(id), Name: req.Name, PricePerDay: req.PricePerDay}
	updated, err := h.service.Update(m)
	if err != nil {
		http.Error(w, "error updating meal plan", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteMealPlan DELETE /api/mealplans/delete?id=
func (h *MealPlanHandler) DeleteMealPlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.service.Delete(id); err != nil {
		http.Error(w, "error deleting meal plan", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "meal plan deleted"})
}
