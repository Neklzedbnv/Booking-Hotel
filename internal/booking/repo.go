package booking

import (
	"database/sql"

	"Gofinal/internal/domain"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(b domain.Booking) (domain.Booking, error) {
	if b.Status == "" {
		b.Status = "pending"
	}
	query := `
		INSERT INTO bookings
		(user_id, room_id, mealplan_id, package_id, start_date, end_date, stay_days, total_price, status, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		b.UserID,
		b.RoomID,
		b.MealplanID,
		b.PackageID,
		b.StartDate,
		b.EndDate,
		b.StayDays,
		b.TotalPrice,
		b.Status,
		b.CreatedAt,
	).Scan(&b.ID)

	return b, err
}

func (r *Repo) GetAll() ([]domain.Booking, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, room_id, mealplan_id, package_id,
		       start_date, end_date, stay_days, total_price, COALESCE(status, 'pending'), created_at
		FROM bookings
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Booking
	for rows.Next() {
		var b domain.Booking
		if err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.RoomID,
			&b.MealplanID,
			&b.PackageID,
			&b.StartDate,
			&b.EndDate,
			&b.StayDays,
			&b.TotalPrice,
			&b.Status,
			&b.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, b)
	}

	return list, nil
}

func (r *Repo) GetByID(id int) (domain.Booking, error) {
	var b domain.Booking
	err := r.db.QueryRow(`
		SELECT id, user_id, room_id, mealplan_id, package_id,
		       start_date, end_date, stay_days, total_price, COALESCE(status, 'pending'), created_at
		FROM bookings
		WHERE id=$1
	`, id).Scan(
		&b.ID,
		&b.UserID,
		&b.RoomID,
		&b.MealplanID,
		&b.PackageID,
		&b.StartDate,
		&b.EndDate,
		&b.StayDays,
		&b.TotalPrice,
		&b.Status,
		&b.CreatedAt,
	)

	return b, err
}

func (r *Repo) Update(b domain.Booking) (domain.Booking, error) {
	query := `
		UPDATE bookings
		SET start_date=$1, end_date=$2, stay_days=$3, total_price=$4
		WHERE id=$5
	`

	_, err := r.db.Exec(
		query,
		b.StartDate,
		b.EndDate,
		b.StayDays,
		b.TotalPrice,
		b.ID,
	)

	return b, err
}

func (r *Repo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM bookings WHERE id=$1`, id)
	return err
}

func (r *Repo) GetAllByUser(userID int) ([]domain.Booking, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, room_id, mealplan_id, package_id,
		       start_date, end_date, stay_days, total_price, COALESCE(status, 'pending'), created_at
		FROM bookings
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Booking
	for rows.Next() {
		var b domain.Booking
		if err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.RoomID,
			&b.MealplanID,
			&b.PackageID,
			&b.StartDate,
			&b.EndDate,
			&b.StayDays,
			&b.TotalPrice,
			&b.Status,
			&b.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, b)
	}

	return list, nil
}
