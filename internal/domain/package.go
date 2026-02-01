package domain

import "time"

type Package struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	BasePriceModifier float64   `json:"base_price_modifier"`
	CreatedAt         time.Time `json:"created_at"`
}
