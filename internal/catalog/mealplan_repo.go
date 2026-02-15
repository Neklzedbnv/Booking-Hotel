package catalog

import (
	"Gofinal/internal/domain"
	"database/sql"
)

// MealPlanRepo — repository for mealplans table
type MealPlanRepo struct {
	db *sql.DB
}

func NewMealPlanRepo(db *sql.DB) *MealPlanRepo {
	return &MealPlanRepo{db: db}
}

func (r *MealPlanRepo) Create(m domain.MealPlan) (*domain.MealPlan, error) {
	query := `INSERT INTO mealplans (name, price_per_day) VALUES ($1, $2) RETURNING id, created_at`
	var ca sql.NullTime
	err := r.db.QueryRow(query, m.Name, m.PricePerDay).Scan(&m.ID, &ca)
	if err != nil {
		return nil, err
	}
	if ca.Valid {
		m.CreatedAt = ca.Time
	}
	return &m, nil
}

func (r *MealPlanRepo) GetByID(id int) (*domain.MealPlan, error) {
	m := &domain.MealPlan{}
	var name sql.NullString
	var ca sql.NullTime
	err := r.db.QueryRow(`SELECT id, name, price_per_day, created_at FROM mealplans WHERE id=$1`, id).
		Scan(&m.ID, &name, &m.PricePerDay, &ca)
	if err != nil {
		return nil, err
	}
	m.Name = name.String
	if ca.Valid {
		m.CreatedAt = ca.Time
	}
	return m, nil
}

func (r *MealPlanRepo) List() ([]domain.MealPlan, error) {
	rows, err := r.db.Query(`SELECT id, name, price_per_day, created_at FROM mealplans ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.MealPlan
	for rows.Next() {
		var m domain.MealPlan
		var name sql.NullString
		var ppd sql.NullFloat64
		var ca sql.NullTime
		if err := rows.Scan(&m.ID, &name, &ppd, &ca); err != nil {
			return nil, err
		}
		m.Name = name.String
		m.PricePerDay = ppd.Float64
		if ca.Valid {
			m.CreatedAt = ca.Time
		}
		list = append(list, m)
	}
	return list, nil
}

func (r *MealPlanRepo) Update(m domain.MealPlan) (*domain.MealPlan, error) {
	query := `UPDATE mealplans SET name=$1, price_per_day=$2 WHERE id=$3 RETURNING id, name, price_per_day, created_at`
	var ca sql.NullTime
	err := r.db.QueryRow(query, m.Name, m.PricePerDay, m.ID).Scan(&m.ID, &m.Name, &m.PricePerDay, &ca)
	if err != nil {
		return nil, err
	}
	if ca.Valid {
		m.CreatedAt = ca.Time
	}
	return &m, nil
}

func (r *MealPlanRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM mealplans WHERE id=$1`, id)
	return err
}
