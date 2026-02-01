package domain

import "time"

type User struct {
	ID           int64     `json:"id"`
	FullName     string    `json:"fullname"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}