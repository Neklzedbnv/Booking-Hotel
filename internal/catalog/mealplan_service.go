package catalog

import "Gofinal/internal/domain"

// MealPlanService — business logic for meal plans
type MealPlanService struct {
	repo *MealPlanRepo
}

func NewMealPlanService(repo *MealPlanRepo) *MealPlanService {
	return &MealPlanService{repo: repo}
}

func (s *MealPlanService) Create(m domain.MealPlan) (*domain.MealPlan, error) {
	return s.repo.Create(m)
}

func (s *MealPlanService) GetByID(id int) (*domain.MealPlan, error) {
	return s.repo.GetByID(id)
}

func (s *MealPlanService) List() ([]domain.MealPlan, error) {
	return s.repo.List()
}

func (s *MealPlanService) Update(m domain.MealPlan) (*domain.MealPlan, error) {
	return s.repo.Update(m)
}

func (s *MealPlanService) Delete(id int) error {
	return s.repo.Delete(id)
}
