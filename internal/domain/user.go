package domain

import "time"

type User struct {
	ID           int       `json:"id"`
	FullName     string    `json:"fullname"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	IsBlocked    bool      `json:"is_blocked"`
	CreatedAt    time.Time `json:"created_at"`
}
// UserInput is used for registration and login requests.