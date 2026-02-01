package domain

import "time"

type Payment struct {
	ID        int       `json:"id"`
	BookingID int       `json:"booking_id"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
