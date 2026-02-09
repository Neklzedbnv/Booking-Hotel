package domain

import "time"


type Package struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`             
	Description     string    `json:"description"`      
	PriceModifier   float64   `json:"price_modifier"`   
	IsActive        bool      `json:"is_active"`        
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type RoomPackage struct {
	ID        int64     `json:"id"`
	RoomID    int64     `json:"room_id"`
	PackageID int       `json:"package_id"`
	CreatedAt time.Time `json:"created_at"`
}
