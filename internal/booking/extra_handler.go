package booking

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Gofinal/internal/domain"
)

// ─── Review Handler ──────────────────────────────────────────────────────────

type ReviewHandler struct {
	service *ReviewService
}

func NewReviewHandler(s *ReviewService) *ReviewHandler {
	return &ReviewHandler{service: s}
}

// CreateReview POST /api/reviews
func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		BookingID int    `json:"booking_id"`
		Rating    int    `json:"rating"`
		Comment   string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	rev := domain.Review{BookingID: req.BookingID, Rating: req.Rating, Comment: req.Comment}
	created, err := h.service.Create(rev)
	if err != nil {
		http.Error(w, "error creating review", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetReview GET /api/reviews/get?id=
func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	rev, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rev)
}

// ListReviews GET /api/reviews/list?booking_id= (optional)
func (h *ReviewHandler) ListReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bidStr := r.URL.Query().Get("booking_id")
	if bidStr != "" {
		bid, err := strconv.Atoi(bidStr)
		if err != nil {
			http.Error(w, "invalid booking_id", http.StatusBadRequest)
			return
		}
		list, err := h.service.ListByBooking(bid)
		if err != nil {
			http.Error(w, "error listing reviews", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
		return
	}

	list, err := h.service.ListAll()
	if err != nil {
		http.Error(w, "error listing reviews", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// UpdateReview PUT /api/reviews/update?id=
func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
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
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	rev := domain.Review{ID: id, Rating: req.Rating, Comment: req.Comment}
	updated, err := h.service.Update(rev)
	if err != nil {
		http.Error(w, "error updating review", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteReview DELETE /api/reviews/delete?id=
func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "error deleting review", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "review deleted"})
}

// ─── Payment Handler ─────────────────────────────────────────────────────────

type PaymentHandler struct {
	service *PaymentService
}

func NewPaymentHandler(s *PaymentService) *PaymentHandler {
	return &PaymentHandler{service: s}
}

// CreatePayment POST /api/payments
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		BookingID int     `json:"booking_id"`
		Method    string  `json:"method"`
		Amount    float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	p := domain.Payment{BookingID: req.BookingID, Method: req.Method, Status: "pending", Amount: req.Amount}
	created, err := h.service.Create(p)
	if err != nil {
		http.Error(w, "error creating payment", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetPayment GET /api/payments/get?id=
func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	p, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// ListPayments GET /api/payments/list?booking_id= (optional)
func (h *PaymentHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bidStr := r.URL.Query().Get("booking_id")
	if bidStr != "" {
		bid, err := strconv.Atoi(bidStr)
		if err != nil {
			http.Error(w, "invalid booking_id", http.StatusBadRequest)
			return
		}
		list, err := h.service.ListByBooking(bid)
		if err != nil {
			http.Error(w, "error listing payments", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
		return
	}

	list, err := h.service.ListAll()
	if err != nil {
		http.Error(w, "error listing payments", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// UpdatePaymentStatus PUT /api/payments/update?id=
func (h *PaymentHandler) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
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
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	p, err := h.service.UpdateStatus(id, req.Status)
	if err != nil {
		http.Error(w, "error updating payment", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// DeletePayment DELETE /api/payments/delete?id=
func (h *PaymentHandler) DeletePayment(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "error deleting payment", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "payment deleted"})
}
