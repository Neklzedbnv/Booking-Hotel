package domain

import "time"

type Review struct {
	ID        int       `json:"id"`
	BookingID int       `json:"booking_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
