package catalog

import "Gofinal/internal/domain"

// SvcService — business logic for services
type SvcService struct {
	repo *SvcRepo
}

func NewSvcService(repo *SvcRepo) *SvcService {
	return &SvcService{repo: repo}
}

func (s *SvcService) Create(svc domain.Service) (*domain.Service, error) {
	return s.repo.Create(svc)
}

func (s *SvcService) GetByID(id int) (*domain.Service, error) {
	return s.repo.GetByID(id)
}

func (s *SvcService) List() ([]domain.Service, error) {
	return s.repo.List()
}

func (s *SvcService) Update(svc domain.Service) (*domain.Service, error) {
	return s.repo.Update(svc)
}

func (s *SvcService) Delete(id int) error {
	return s.repo.Delete(id)
}
