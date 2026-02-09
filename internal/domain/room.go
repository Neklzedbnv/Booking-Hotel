package domain

import "time"


type RoomType struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`        
	Capacity  int       `json:"capacity"`    
	BasePrice float64   `json:"base_price"`  
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type Room struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`             
	TypeID    int64     `json:"type_id"`           
	Type      *RoomType `json:"type,omitempty"`    
	Capacity  int       `json:"capacity"`          
	Price     float64   `json:"price"`             
	Status    string    `json:"status"`            
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
