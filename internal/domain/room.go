package domain

import "time"

type Room struct {
	ID         int       `json:"id"`
	Code       string    `json:"code"`
	Capacity   int       `json:"capacity"`
	PriceBase  float64   `json:"price_base"`
	RoomTypeID int       `json:"room_type_id"`
	CreatedAt  time.Time `json:"created_at"`
}
