package booking

import (
	"Gofinal/internal/domain"
	"database/sql"
)

// ReviewRepo — repository for reviews table
type ReviewRepo struct {
	db *sql.DB
}

func NewReviewRepo(db *sql.DB) *ReviewRepo {
	return &ReviewRepo{db: db}
}

func (r *ReviewRepo) Create(rev domain.Review) (*domain.Review, error) {
	query := `INSERT INTO reviews (booking_id, rating, comment) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.db.QueryRow(query, rev.BookingID, rev.Rating, rev.Comment).Scan(&rev.ID, &rev.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &rev, nil
}

func (r *ReviewRepo) GetByID(id int) (*domain.Review, error) {
	rev := &domain.Review{}
	err := r.db.QueryRow(`SELECT id, booking_id, rating, comment, created_at FROM reviews WHERE id=$1`, id).
		Scan(&rev.ID, &rev.BookingID, &rev.Rating, &rev.Comment, &rev.CreatedAt)
	if err != nil {
		return nil, err
	}
	return rev, nil
}

func (r *ReviewRepo) ListByBooking(bookingID int) ([]domain.Review, error) {
	rows, err := r.db.Query(`SELECT id, booking_id, rating, comment, created_at FROM reviews WHERE booking_id=$1 ORDER BY id`, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Review
	for rows.Next() {
		var rev domain.Review
		if err := rows.Scan(&rev.ID, &rev.BookingID, &rev.Rating, &rev.Comment, &rev.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, rev)
	}
	return list, nil
}

func (r *ReviewRepo) ListAll() ([]domain.Review, error) {
	rows, err := r.db.Query(`SELECT id, booking_id, rating, comment, created_at FROM reviews ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Review
	for rows.Next() {
		var rev domain.Review
		if err := rows.Scan(&rev.ID, &rev.BookingID, &rev.Rating, &rev.Comment, &rev.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, rev)
	}
	return list, nil
}

func (r *ReviewRepo) Update(rev domain.Review) (*domain.Review, error) {
	query := `UPDATE reviews SET rating=$1, comment=$2 WHERE id=$3 RETURNING id, booking_id, rating, comment, created_at`
	err := r.db.QueryRow(query, rev.Rating, rev.Comment, rev.ID).
		Scan(&rev.ID, &rev.BookingID, &rev.Rating, &rev.Comment, &rev.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &rev, nil
}

func (r *ReviewRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM reviews WHERE id=$1`, id)
	return err
}

// PaymentRepo — repository for payments table
type PaymentRepo struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) Create(p domain.Payment) (*domain.Payment, error) {
	query := `INSERT INTO payments (booking_id, method, status, amount) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err := r.db.QueryRow(query, p.BookingID, p.Method, p.Status, p.Amount).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PaymentRepo) GetByID(id int) (*domain.Payment, error) {
	p := &domain.Payment{}
	err := r.db.QueryRow(`SELECT id, booking_id, method, status, amount, created_at FROM payments WHERE id=$1`, id).
		Scan(&p.ID, &p.BookingID, &p.Method, &p.Status, &p.Amount, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PaymentRepo) ListByBooking(bookingID int) ([]domain.Payment, error) {
	rows, err := r.db.Query(`SELECT id, booking_id, method, status, amount, created_at FROM payments WHERE booking_id=$1 ORDER BY id`, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Payment
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(&p.ID, &p.BookingID, &p.Method, &p.Status, &p.Amount, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *PaymentRepo) ListAll() ([]domain.Payment, error) {
	rows, err := r.db.Query(`SELECT id, booking_id, method, status, amount, created_at FROM payments ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Payment
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(&p.ID, &p.BookingID, &p.Method, &p.Status, &p.Amount, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *PaymentRepo) UpdateStatus(id int, status string) (*domain.Payment, error) {
	p := &domain.Payment{}
	query := `UPDATE payments SET status=$1 WHERE id=$2 RETURNING id, booking_id, method, status, amount, created_at`
	err := r.db.QueryRow(query, status, id).Scan(&p.ID, &p.BookingID, &p.Method, &p.Status, &p.Amount, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PaymentRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM payments WHERE id=$1`, id)
	return err
}
