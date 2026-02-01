package catalog

import "Gofinal/internal/domain"

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

// Read-only use case (stub OK for milestone)
func (s *Service) GetAll() ([]domain.Room, error) {
	return s.repo.GetAll()
}
