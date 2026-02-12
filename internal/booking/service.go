package booking

import (
	"database/sql"
	"time"

	"Gofinal/internal/domain"
)

// RoomGetter interface for getting room
type RoomGetter interface {
	GetRoomByID(id int64) (*domain.Room, error)
	UpdateRoomStatus(id int64, status string) error
}

// MealPlanGetter interface for getting meal plan
type MealPlanGetter interface {
	GetByID(id int) (*domain.MealPlan, error)
}

// PackageGetter interface for getting package
type PackageGetter interface {
	GetPackageByID(id int) (*domain.Package, error)
}

// PaymentCreator interface for creating payment
type PaymentCreator interface {
	Create(p domain.Payment) (*domain.Payment, error)
}

type Service struct {
	repo        domain.BookingRepository
	roomRepo    RoomGetter
	mealRepo    MealPlanGetter
	packageRepo PackageGetter
	paymentRepo PaymentCreator
}

func NewService(repo domain.BookingRepository, roomRepo RoomGetter, mealRepo MealPlanGetter, packageRepo PackageGetter, paymentRepo PaymentCreator) *Service {
	return &Service{
		repo:        repo,
		roomRepo:    roomRepo,
		mealRepo:    mealRepo,
		packageRepo: packageRepo,
		paymentRepo: paymentRepo,
	}
}

func (s *Service) Create(b domain.Booking, paymentMethod string) (domain.Booking, error) {
	b.StayDays = int(b.EndDate.Sub(b.StartDate).Hours() / 24)
	if b.StayDays < 1 {
		b.StayDays = 1
	}
	b.CreatedAt = time.Now()

	// Calculate price
	var totalPrice float64 = 0

	// Room price
	if s.roomRepo != nil {
		room, err := s.roomRepo.GetRoomByID(int64(b.RoomID))
		if err == nil && room != nil {
			totalPrice += room.Price * float64(b.StayDays)
		}
	}

	// Meal plan price
	if s.mealRepo != nil && b.MealplanID != nil && *b.MealplanID > 0 {
		meal, err := s.mealRepo.GetByID(*b.MealplanID)
		if err == nil && meal != nil {
			totalPrice += meal.PricePerDay * float64(b.StayDays)
		}
	}

	// Package discount
	if s.packageRepo != nil && b.PackageID != nil && *b.PackageID > 0 {
		pkg, err := s.packageRepo.GetPackageByID(*b.PackageID)
		if err == nil && pkg != nil {
			totalPrice = totalPrice * (1 - pkg.PriceModifier)
		}
	}

	b.TotalPrice = totalPrice

	// Create booking
	res, err := s.repo.Create(b)
	if err != nil {
		return res, err
	}

	// Update room status to "booked"
	if s.roomRepo != nil {
		_ = s.roomRepo.UpdateRoomStatus(int64(b.RoomID), "booked")
	}

	// Create payment record
	if s.paymentRepo != nil && paymentMethod != "" {
		payment := domain.Payment{
			BookingID: res.ID,
			Method:    paymentMethod,
			Status:    "pending",
			Amount:    res.TotalPrice,
		}
		_, _ = s.paymentRepo.Create(payment)
	}

	return res, nil
}

func (s *Service) GetAll() ([]domain.Booking, error) {
	return s.repo.GetAll()
}

func (s *Service) GetAllByUser(userID int) ([]domain.Booking, error) {
	return s.repo.GetAllByUser(userID)
}

func (s *Service) GetByID(id int) (domain.Booking, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(b domain.Booking) (domain.Booking, error) {
	b.StayDays = int(b.EndDate.Sub(b.StartDate).Hours() / 24)
	return s.repo.Update(b)
}

func (s *Service) Delete(id int) error {
	// Get booking to find room_id
	booking, err := s.repo.GetByID(id)
	if err == nil && s.roomRepo != nil {
		// Release room
		_ = s.roomRepo.UpdateRoomStatus(int64(booking.RoomID), "available")
	}
	return s.repo.Delete(id)
}

// CancelBooking cancels booking (marks status as cancelled)
func (s *Service) CancelBooking(id int, userID int) error {
	booking, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}

	// Check that booking belongs to user
	if booking.UserID != userID {
		return sql.ErrNoRows // return "not found" if not owner
	}

	// Release room
	if s.roomRepo != nil {
		_ = s.roomRepo.UpdateRoomStatus(int64(booking.RoomID), "available")
	}

	return s.repo.Delete(id)
}
