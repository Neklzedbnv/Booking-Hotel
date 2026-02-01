package domain

import "time"

type Room struct {
	ID         int64     `json:"id"`
	Code       string    `json:"code"`
	Capacity   int       `json:"capacity"`
	PriceBase  float64   `json:"price_base"`
	RoomTypeID int64     `json:"room_type_id"`
	CreatedAt  time.Time `json:"created_at"`
}
