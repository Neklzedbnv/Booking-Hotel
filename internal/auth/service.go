package auth

import (
	"Gofinal/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}


func (s *Service) SaveUser(user domain.User) error {
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword) 

	return s.repo.SaveUser(user)
}


func (s *Service) FindUserByEmail(email string) (domain.User, error) {
	return s.repo.FindUserByEmail(email)
}
