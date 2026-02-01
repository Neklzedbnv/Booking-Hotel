package domain

import "time"

type MealPlan struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	PricePerDay float64 `json:"price_per_day"`
	CreatedAt time.Time `json:"created_at"`
}
