package catalog

import (
	"Gofinal/internal/domain"
	"database/sql"
)

// SvcRepo — repository for services table
type SvcRepo struct {
	db *sql.DB
}

func NewSvcRepo(db *sql.DB) *SvcRepo {
	return &SvcRepo{db: db}
}

func (r *SvcRepo) Create(s domain.Service) (*domain.Service, error) {
	query := `INSERT INTO services (name, price) VALUES ($1, $2) RETURNING id, created_at`
	err := r.db.QueryRow(query, s.Name, s.Price).Scan(&s.ID, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SvcRepo) GetByID(id int) (*domain.Service, error) {
	s := &domain.Service{}
	err := r.db.QueryRow(`SELECT id, name, price, created_at FROM services WHERE id=$1`, id).
		Scan(&s.ID, &s.Name, &s.Price, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *SvcRepo) List() ([]domain.Service, error) {
	rows, err := r.db.Query(`SELECT id, name, price, created_at FROM services ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Service
	for rows.Next() {
		var s domain.Service
		if err := rows.Scan(&s.ID, &s.Name, &s.Price, &s.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

func (r *SvcRepo) Update(s domain.Service) (*domain.Service, error) {
	query := `UPDATE services SET name=$1, price=$2 WHERE id=$3 RETURNING id, name, price, created_at`
	err := r.db.QueryRow(query, s.Name, s.Price, s.ID).Scan(&s.ID, &s.Name, &s.Price, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SvcRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM services WHERE id=$1`, id)
	return err
}
