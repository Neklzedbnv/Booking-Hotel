package catalog

import (
	"database/sql"
	"fmt"

	"Gofinal/internal/domain"
)

type PackageRepo struct {
	db *sql.DB
}

func NewPackageRepo(db *sql.DB) *PackageRepo {
	return &PackageRepo{db: db}
}


func (r *PackageRepo) CreatePackage(pkg domain.Package) (*domain.Package, error) {
	query := `
		INSERT INTO packages (name, description, price_modifier, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(query, pkg.Name, pkg.Description, pkg.PriceModifier, pkg.IsActive, pkg.CreatedAt, pkg.UpdatedAt).Scan(&id)
	if err != nil {
		return nil, err
	}

	pkg.ID = id
	return &pkg, nil
}


func (r *PackageRepo) GetPackageByID(id int) (*domain.Package, error) {
	query := `
		SELECT id, name, description, price_modifier, is_active, created_at, updated_at
		FROM packages
		WHERE id = $1
	`

	pkg := &domain.Package{}
	err := r.db.QueryRow(query, id).Scan(
		&pkg.ID, &pkg.Name, &pkg.Description, &pkg.PriceModifier, &pkg.IsActive, &pkg.CreatedAt, &pkg.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}


func (r *PackageRepo) ListPackages(onlyActive bool) ([]domain.Package, error) {
	query := "SELECT id, name, description, price_modifier, is_active, created_at, updated_at FROM packages"

	if onlyActive {
		query += " WHERE is_active = true"
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packages []domain.Package
	for rows.Next() {
		pkg := domain.Package{}
		err := rows.Scan(&pkg.ID, &pkg.Name, &pkg.Description, &pkg.PriceModifier, &pkg.IsActive, &pkg.CreatedAt, &pkg.UpdatedAt)
		if err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}


func (r *PackageRepo) UpdatePackage(id int, name *string, description *string, priceModifier *float64, isActive *bool) (*domain.Package, error) {
	query := "UPDATE packages SET "
	args := []interface{}{}
	argIndex := 1

	if name != nil {
		query += fmt.Sprintf("name = $%d, ", argIndex)
		args = append(args, *name)
		argIndex++
	}

	if description != nil {
		query += fmt.Sprintf("description = $%d, ", argIndex)
		args = append(args, *description)
		argIndex++
	}

	if priceModifier != nil {
		query += fmt.Sprintf("price_modifier = $%d, ", argIndex)
		args = append(args, *priceModifier)
		argIndex++
	}

	if isActive != nil {
		query += fmt.Sprintf("is_active = $%d, ", argIndex)
		args = append(args, *isActive)
		argIndex++
	}

	query += fmt.Sprintf("updated_at = NOW() WHERE id = $%d RETURNING id, name, description, price_modifier, is_active, created_at, updated_at", argIndex)
	args = append(args, id)

	pkg := &domain.Package{}
	err := r.db.QueryRow(query, args...).Scan(
		&pkg.ID, &pkg.Name, &pkg.Description, &pkg.PriceModifier, &pkg.IsActive, &pkg.CreatedAt, &pkg.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}


func (r *PackageRepo) DeletePackage(id int) error {
	query := "UPDATE packages SET is_active = false, updated_at = NOW() WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}


func (r *PackageRepo) AttachPackageToRoom(roomID int64, packageID int) error {
	query := `
		INSERT INTO room_packages (room_id, package_id, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (room_id, package_id) DO NOTHING
	`

	_, err := r.db.Exec(query, roomID, packageID)
	return err
}


func (r *PackageRepo) GetRoomPackages(roomID int64) ([]domain.Package, error) {
	query := `
		SELECT p.id, p.name, p.description, p.price_modifier, p.is_active, p.created_at, p.updated_at
		FROM packages p
		JOIN room_packages rp ON p.id = rp.package_id
		WHERE rp.room_id = $1 AND p.is_active = true
	`

	rows, err := r.db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packages []domain.Package
	for rows.Next() {
		pkg := domain.Package{}
		err := rows.Scan(&pkg.ID, &pkg.Name, &pkg.Description, &pkg.PriceModifier, &pkg.IsActive, &pkg.CreatedAt, &pkg.UpdatedAt)
		if err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}


func (r *PackageRepo) DetachPackageFromRoom(roomID int64, packageID int) error {
	query := "DELETE FROM room_packages WHERE room_id = $1 AND package_id = $2"
	_, err := r.db.Exec(query, roomID, packageID)
	return err
}
