package domain

type BookingRepository interface {
	Create(b Booking) (Booking, error)
	GetAll() ([]Booking, error)
	GetAllByUser(userID int) ([]Booking, error)
	GetByID(id int) (Booking, error)
	Update(b Booking) (Booking, error)
	Delete(id int) error
}
