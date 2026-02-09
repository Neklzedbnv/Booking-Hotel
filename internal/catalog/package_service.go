package catalog

import (
	"fmt"

	"Gofinal/internal/domain"
)

type PackageService struct {
	repo *PackageRepo
}

func NewPackageService(repo *PackageRepo) *PackageService {
	return &PackageService{repo: repo}
}


func (s *PackageService) CreatePackage(pkg domain.Package) (*domain.Package, error) {
	return s.repo.CreatePackage(pkg)
}


func (s *PackageService) GetPackageByID(id int) (*domain.Package, error) {
	return s.repo.GetPackageByID(id)
}


func (s *PackageService) ListPackages(onlyActive bool) ([]domain.Package, error) {
	return s.repo.ListPackages(onlyActive)
}


func (s *PackageService) UpdatePackage(id int, name *string, description *string, priceModifier *float64, isActive *bool) (*domain.Package, error) {
	if name == nil && description == nil && priceModifier == nil && isActive == nil {
		return nil, fmt.Errorf("no fields to update")
	}

	return s.repo.UpdatePackage(id, name, description, priceModifier, isActive)
}


func (s *PackageService) DeletePackage(id int) error {
	return s.repo.DeletePackage(id)
}


func (s *PackageService) AttachPackageToRoom(roomID int64, packageID int) error {
	return s.repo.AttachPackageToRoom(roomID, packageID)
}


func (s *PackageService) GetRoomPackages(roomID int64) ([]domain.Package, error) {
	return s.repo.GetRoomPackages(roomID)
}


func (s *PackageService) DetachPackageFromRoom(roomID int64, packageID int) error {
	return s.repo.DetachPackageFromRoom(roomID, packageID)
}


func (s *PackageService) CalculatePackagePrice(packageIDs []int) (float64, error) {
	total := 0.0

	for _, id := range packageIDs {
		pkg, err := s.GetPackageByID(id)
		if err != nil {
			return 0, fmt.Errorf("package %d not found", id)
		}

		if !pkg.IsActive {
			return 0, fmt.Errorf("package %d is not active", id)
		}

		total += pkg.PriceModifier
	}

	return total, nil
}
