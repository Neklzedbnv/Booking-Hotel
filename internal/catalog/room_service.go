package catalog

import (
	"fmt"
	"time"

	"Gofinal/internal/domain"
)

type RoomService struct {
	repo *RoomRepo
}

func NewRoomService(repo *RoomRepo) *RoomService {
	return &RoomService{repo: repo}
}


func (s *RoomService) CreateRoom(room domain.Room) (*domain.Room, error) {
	return s.repo.CreateRoom(room)
}


func (s *RoomService) GetRoomByID(id int64) (*domain.Room, error) {
	return s.repo.GetRoomByID(id)
}


func (s *RoomService) ListRooms(status, typeID string) ([]domain.Room, error) {
	return s.repo.ListRooms(status, typeID)
}


func (s *RoomService) UpdateRoom(id int64, price *float64, status *string) (*domain.Room, error) {
	if price == nil && status == nil {
		return nil, fmt.Errorf("no fields to update")
	}

	return s.repo.UpdateRoom(id, price, status)
}


func (s *RoomService) DeleteRoom(id int64) error {
	return s.repo.DeleteRoom(id)
}


func (s *RoomService) CreateRoomType(roomType domain.RoomType) (*domain.RoomType, error) {
	return s.repo.CreateRoomType(roomType)
}


func (s *RoomService) ListRoomTypes() ([]domain.RoomType, error) {
	return s.repo.ListRoomTypes()
}


func (s *RoomService) CheckAvailability(checkIn, checkOut time.Time, capacity *int, typeID *int64) ([]domain.Room, error) {
	
	var filters []string
	var args []interface{}

	if capacity != nil {
		filters = append(filters, "capacity >= ?")
		args = append(args, *capacity)
	}

	if typeID != nil {
		filters = append(filters, "type_id = ?")
		args = append(args, *typeID)
	}

	
	filters = append(filters, "status = 'available'")

	return s.repo.GetAvailableRooms(checkIn, checkOut, filters, args)
}
