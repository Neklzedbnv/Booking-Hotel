package admin

import "Gofinal/internal/domain"

type Service struct {
	repo *Repo
}

func NewService(r *Repo) *Service {
	return &Service{repo: r}
}

func (s *Service) ListUsers() ([]domain.User, error) {
	return s.repo.ListUsers()
}

func (s *Service) UpdateUserRole(userID int, role string) error {
	return s.repo.UpdateUserRole(userID, role)
}

func (s *Service) SetAdminByEmail(email string) error {
	return s.repo.UpdateUserRoleByEmail(email, "admin")
}

func (s *Service) ResetPasswordByEmail(email string, hashedPassword string) error {
	return s.repo.ResetPasswordByEmail(email, hashedPassword)
}

func (s *Service) BlockUser(userID int, blocked bool) error {
	return s.repo.BlockUser(userID, blocked)
}

func (s *Service) GetUserByID(userID int) (domain.User, error) {
	return s.repo.GetUserByID(userID)
}

func (s *Service) UpdateBookingStatus(bookingID int, status string) error {
	return s.repo.UpdateBookingStatus(bookingID, status)
}

func (s *Service) GetBookingsWithDetails() ([]map[string]interface{}, error) {
	return s.repo.GetBookingsWithDetails()
}

func (s *Service) UpdateRoomType(id int64, name string, capacity int, basePrice float64) error {
	return s.repo.UpdateRoomType(id, name, capacity, basePrice)
}

func (s *Service) DeleteRoomType(id int64) error {
	return s.repo.DeleteRoomType(id)
}

func (s *Service) GetDashboardStats() (map[string]interface{}, error) {
	return s.repo.GetDashboardStats()
}
