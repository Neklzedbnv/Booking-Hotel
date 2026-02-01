package domain

type BookingRepository interface {
	Create(b Booking) (Booking, error)
	GetAll() ([]Booking, error)
	GetByID(id int) (Booking, error)
}
