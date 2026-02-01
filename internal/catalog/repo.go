package catalog

import "Gofinal/internal/domain"

type Repo struct{}

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) GetAll() ([]domain.Room, error) {
	return []domain.Room{}, nil
}
