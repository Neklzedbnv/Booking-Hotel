package auth

import (
	"database/sql"
	"Gofinal/internal/domain"
	"fmt"
	"time"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}


func (r *Repo) SaveUser(user domain.User) error {
	query := `
		INSERT INTO users (fullname, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, user.FullName, user.Email, user.PasswordHash, user.Role, time.Now())
	return err
}


func (r *Repo) FindUserByEmail(email string) (domain.User, error) {
	var user domain.User
	query := `
		SELECT id, fullname, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}
		return user, err
	}

	return user, nil
}
