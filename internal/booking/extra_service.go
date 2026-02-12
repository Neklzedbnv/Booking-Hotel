package booking

import "Gofinal/internal/domain"

// ReviewService — business logic for reviews
type ReviewService struct {
	repo *ReviewRepo
}

func NewReviewService(repo *ReviewRepo) *ReviewService {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) Create(rev domain.Review) (*domain.Review, error) { return s.repo.Create(rev) }
func (s *ReviewService) GetByID(id int) (*domain.Review, error)           { return s.repo.GetByID(id) }
func (s *ReviewService) ListByBooking(bid int) ([]domain.Review, error) {
	return s.repo.ListByBooking(bid)
}
func (s *ReviewService) ListAll() ([]domain.Review, error)                { return s.repo.ListAll() }
func (s *ReviewService) Update(rev domain.Review) (*domain.Review, error) { return s.repo.Update(rev) }
func (s *ReviewService) Delete(id int) error                              { return s.repo.Delete(id) }

// PaymentService — business logic for payments
type PaymentService struct {
	repo *PaymentRepo
}

func NewPaymentService(repo *PaymentRepo) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) Create(p domain.Payment) (*domain.Payment, error) { return s.repo.Create(p) }
func (s *PaymentService) GetByID(id int) (*domain.Payment, error)          { return s.repo.GetByID(id) }
func (s *PaymentService) ListByBooking(bid int) ([]domain.Payment, error) {
	return s.repo.ListByBooking(bid)
}
func (s *PaymentService) ListAll() ([]domain.Payment, error) { return s.repo.ListAll() }
func (s *PaymentService) UpdateStatus(id int, st string) (*domain.Payment, error) {
	return s.repo.UpdateStatus(id, st)
}
func (s *PaymentService) Delete(id int) error { return s.repo.Delete(id) }
