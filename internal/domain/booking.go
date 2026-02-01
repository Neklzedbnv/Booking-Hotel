package domain

import "time"

type Booking struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	RoomID     int       `json:"room_id"`
	MealplanID *int      `json:"mealplan_id"`
	PackageID  *int      `json:"package_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	StayDays   int       `json:"stay_days"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
}
