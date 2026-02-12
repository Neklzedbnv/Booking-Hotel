package admin

import (
	"Gofinal/internal/domain"
	"database/sql"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

// ListUsers returns all users
func (r *Repo) ListUsers() ([]domain.User, error) {
	query := `SELECT id, fullname, email, role, is_blocked, created_at FROM users ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.FullName, &u.Email, &u.Role, &u.IsBlocked, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// UpdateUserRole updates user role
func (r *Repo) UpdateUserRole(userID int, role string) error {
	_, err := r.db.Exec(`UPDATE users SET role = $1 WHERE id = $2`, role, userID)
	return err
}

// UpdateUserRoleByEmail updates user role by email
func (r *Repo) UpdateUserRoleByEmail(email string, role string) error {
	_, err := r.db.Exec(`UPDATE users SET role = $1 WHERE email = $2`, role, email)
	return err
}

// ResetPasswordByEmail resets user password by email
func (r *Repo) ResetPasswordByEmail(email string, hashedPassword string) error {
	_, err := r.db.Exec(`UPDATE users SET password_hash = $1 WHERE email = $2`, hashedPassword, email)
	return err
}

// BlockUser blocks/unblocks a user
func (r *Repo) BlockUser(userID int, blocked bool) error {
	_, err := r.db.Exec(`UPDATE users SET is_blocked = $1 WHERE id = $2`, blocked, userID)
	return err
}

// GetUserByID returns user by ID
func (r *Repo) GetUserByID(userID int) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(
		`SELECT id, fullname, email, role, is_blocked, created_at FROM users WHERE id = $1`,
		userID,
	).Scan(&u.ID, &u.FullName, &u.Email, &u.Role, &u.IsBlocked, &u.CreatedAt)
	return u, err
}

// UpdateBookingStatus updates booking status
func (r *Repo) UpdateBookingStatus(bookingID int, status string) error {
	_, err := r.db.Exec(`UPDATE bookings SET status = $1 WHERE id = $2`, status, bookingID)
	return err
}

// GetBookingsWithDetails returns bookings with user and room data
func (r *Repo) GetBookingsWithDetails() ([]map[string]interface{}, error) {
	query := `
		SELECT 
			b.id, b.user_id, b.room_id, b.start_date, b.end_date, 
			b.stay_days, b.total_price, COALESCE(b.status, 'pending') as status, b.created_at,
			u.fullname as user_name, u.email as user_email,
			r.code as room_code
		FROM bookings b
		JOIN users u ON b.user_id = u.id
		JOIN rooms r ON b.room_id = r.id
		ORDER BY b.created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []map[string]interface{}
	for rows.Next() {
		var id, userID, roomID, stayDays int
		var totalPrice float64
		var status, userName, userEmail, roomCode string
		var startDate, endDate, createdAt interface{}

		if err := rows.Scan(&id, &userID, &roomID, &startDate, &endDate, &stayDays, &totalPrice, &status, &createdAt, &userName, &userEmail, &roomCode); err != nil {
			return nil, err
		}

		bookings = append(bookings, map[string]interface{}{
			"id":          id,
			"user_id":     userID,
			"room_id":     roomID,
			"start_date":  startDate,
			"end_date":    endDate,
			"stay_days":   stayDays,
			"total_price": totalPrice,
			"status":      status,
			"created_at":  createdAt,
			"user_name":   userName,
			"user_email":  userEmail,
			"room_code":   roomCode,
		})
	}
	return bookings, nil
}

// UpdateRoomType updates room type
func (r *Repo) UpdateRoomType(id int64, name string, capacity int, basePrice float64) error {
	_, err := r.db.Exec(
		`UPDATE room_types SET name = $1, capacity = $2, base_price = $3, updated_at = NOW() WHERE id = $4`,
		name, capacity, basePrice, id,
	)
	return err
}

// DeleteRoomType deletes room type
func (r *Repo) DeleteRoomType(id int64) error {
	_, err := r.db.Exec(`DELETE FROM room_types WHERE id = $1`, id)
	return err
}

// GetDashboardStats returns dashboard statistics
func (r *Repo) GetDashboardStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// User count
	var userCount int
	r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&userCount)
	stats["users_count"] = userCount

	// Room count
	var roomCount int
	r.db.QueryRow(`SELECT COUNT(*) FROM rooms`).Scan(&roomCount)
	stats["rooms_count"] = roomCount

	// Booking count
	var bookingCount int
	r.db.QueryRow(`SELECT COUNT(*) FROM bookings`).Scan(&bookingCount)
	stats["bookings_count"] = bookingCount

	// Revenue
	var totalRevenue float64
	r.db.QueryRow(`SELECT COALESCE(SUM(total_price), 0) FROM bookings`).Scan(&totalRevenue)
	stats["total_revenue"] = totalRevenue

	// Bookings by status
	rows, _ := r.db.Query(`SELECT COALESCE(status, 'pending'), COUNT(*) FROM bookings GROUP BY status`)
	if rows != nil {
		defer rows.Close()
		statusCounts := make(map[string]int)
		for rows.Next() {
			var status string
			var count int
			rows.Scan(&status, &count)
			statusCounts[status] = count
		}
		stats["bookings_by_status"] = statusCounts
	}

	return stats, nil
}
