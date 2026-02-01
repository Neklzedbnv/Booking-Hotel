package booking

import (
	"Gofinal/internal/domain"
)

type Service struct {
	repo domain.BookingRepository
}

func NewService(repo domain.BookingRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(b domain.Booking) (domain.Booking, error) {
	b.StayDays = int(b.EndDate.Sub(b.StartDate).Hours() / 24)
	return s.repo.Create(b)
}

func (s *Service) GetAll() ([]domain.Booking, error) {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id int) (domain.Booking, error) {
	return s.repo.GetByID(id)
}
