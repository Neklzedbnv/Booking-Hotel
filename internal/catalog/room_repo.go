package catalog

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"Gofinal/internal/domain"
)

type RoomRepo struct {
	db *sql.DB
}

func NewRoomRepo(db *sql.DB) *RoomRepo {
	return &RoomRepo{db: db}
}

func (r *RoomRepo) CreateRoom(room domain.Room) (*domain.Room, error) {
	query := `
		INSERT INTO rooms (code, room_type_id, capacity, price, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRow(query, room.Code, room.TypeID, room.Capacity, room.Price, room.Status, room.CreatedAt, room.UpdatedAt).Scan(&id)
	if err != nil {
		return nil, err
	}

	room.ID = id
	return &room, nil
}

// GetRoomByID gets room by ID
func (r *RoomRepo) GetRoomByID(id int64) (*domain.Room, error) {
	query := `
		SELECT id, code, room_type_id, capacity, price, status, created_at, updated_at
		FROM rooms
		WHERE id = $1
	`

	room := &domain.Room{}
	err := r.db.QueryRow(query, id).Scan(
		&room.ID, &room.Code, &room.TypeID, &room.Capacity, &room.Price, &room.Status, &room.CreatedAt, &room.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return room, nil
}

// ListRooms gets list of rooms with filtering
func (r *RoomRepo) ListRooms(status, typeID string) ([]domain.Room, error) {
	query := "SELECT id, code, room_type_id, capacity, price, status, created_at, updated_at FROM rooms WHERE 1=1"
	var args []interface{}

	if status != "" {
		query += " AND status = $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, status)
	}

	if typeID != "" {
		query += " AND room_type_id = $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, typeID)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []domain.Room
	for rows.Next() {
		room := domain.Room{}
		err := rows.Scan(
			&room.ID, &room.Code, &room.TypeID, &room.Capacity, &room.Price, &room.Status, &room.CreatedAt, &room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// UpdateRoom updates room
func (r *RoomRepo) UpdateRoom(id int64, price *float64, status *string) (*domain.Room, error) {
	updates := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	if price != nil {
		updates = append(updates, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, *price)
		argIndex++
	}

	if status != nil {
		updates = append(updates, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *status)
		argIndex++
	}

	args = append(args, id)

	query := fmt.Sprintf("UPDATE rooms SET %s WHERE id = $%d RETURNING id, code, room_type_id, capacity, price, status, created_at, updated_at", strings.Join(updates, ", "), argIndex)

	room := &domain.Room{}
	err := r.db.QueryRow(query, args...).Scan(
		&room.ID, &room.Code, &room.TypeID, &room.Capacity, &room.Price, &room.Status, &room.CreatedAt, &room.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return room, nil
}

// DeleteRoom deletes room (updates status)
func (r *RoomRepo) DeleteRoom(id int64) error {
	query := "UPDATE rooms SET status = 'deleted', updated_at = NOW() WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

// CreateRoomType creates new room type
func (r *RoomRepo) CreateRoomType(roomType domain.RoomType) (*domain.RoomType, error) {
	query := `
		INSERT INTO room_types (name, capacity, base_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRow(query, roomType.Name, roomType.Capacity, roomType.BasePrice, roomType.CreatedAt, roomType.UpdatedAt).Scan(&id)
	if err != nil {
		return nil, err
	}

	roomType.ID = id
	return &roomType, nil
}

// ListRoomTypes gets all room types
func (r *RoomRepo) ListRoomTypes() ([]domain.RoomType, error) {
	query := "SELECT id, name, capacity, base_price, created_at, updated_at FROM room_types"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []domain.RoomType
	for rows.Next() {
		rt := domain.RoomType{}
		err := rows.Scan(&rt.ID, &rt.Name, &rt.Capacity, &rt.BasePrice, &rt.CreatedAt, &rt.UpdatedAt)
		if err != nil {
			return nil, err
		}
		types = append(types, rt)
	}

	return types, nil
}

// GetAvailableRooms gets all available rooms for selected dates
func (r *RoomRepo) GetAvailableRooms(checkIn, checkOut time.Time, filters []string, args []interface{}) ([]domain.Room, error) {
	baseQuery := `
		SELECT DISTINCT r.id, r.code, r.room_type_id, r.capacity, r.price, r.status, r.created_at, r.updated_at
		FROM rooms r
		WHERE r.status = 'available'
		AND r.id NOT IN (
			SELECT room_id FROM bookings 
			WHERE status IN ('pending', 'confirmed')
			AND (check_in_date < $1 AND check_out_date > $2)
		)
	`

	// Add date parameters to the beginning of args
	args = append([]interface{}{checkOut, checkIn}, args...)

	if len(filters) > 0 {
		for i, filter := range filters {
			// Recalculate parameter indices
			updatedFilter := filter
			baseQuery += " AND " + updatedFilter
			i++ // for next parameters
		}
	}

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []domain.Room
	for rows.Next() {
		room := domain.Room{}
		err := rows.Scan(
			&room.ID, &room.Code, &room.TypeID, &room.Capacity, &room.Price, &room.Status, &room.CreatedAt, &room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// UpdateRoomStatus updates room status
func (r *RoomRepo) UpdateRoomStatus(id int64, status string) error {
	_, err := r.db.Exec(`UPDATE rooms SET status = $1, updated_at = $2 WHERE id = $3`, status, time.Now(), id)
	return err
}
