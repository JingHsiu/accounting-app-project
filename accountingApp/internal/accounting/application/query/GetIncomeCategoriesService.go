package query

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
)

type GetIncomeCategoriesService struct {
	incomeCategoryRepo repository.IncomeCategoryRepository
}

func NewGetIncomeCategoriesService(incomeCategoryRepo repository.IncomeCategoryRepository) *GetIncomeCategoriesService {
	return &GetIncomeCategoriesService{incomeCategoryRepo: incomeCategoryRepo}
}

func (s *GetIncomeCategoriesService) Execute(input usecase.GetIncomeCategoriesInput) common.Output {
	// First, try to get categories for the user
	categories, err := s.incomeCategoryRepo.FindByUserID(input.UserID)
	if err != nil {
		return usecase.GetIncomeCategoriesOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Failed to retrieve income categories: %v", err),
		}
	}

	// If no categories found, the user needs to have default categories initialized
	if len(categories) == 0 {
		return usecase.GetIncomeCategoriesOutput{
			ID:         input.UserID,
			ExitCode:   common.Success,
			Message:    "No income categories found. Please initialize default categories.",
			Categories: []usecase.CategoryData{},
		}
	}

	// Convert domain models to API response format
	categoriesData := make([]usecase.CategoryData, len(categories))
	for i, category := range categories {
		// Convert subcategories
		subcategories := make([]usecase.SubcategoryData, len(category.Subcategories))
		for j, subcategory := range category.Subcategories {
			subcategories[j] = usecase.SubcategoryData{
				ID:   subcategory.ID,
				Name: subcategory.Name.Value,
			}
		}

		categoriesData[i] = usecase.CategoryData{
			ID:            category.ID,
			Name:          category.Name.Value,
			Type:          "income",
			CreatedAt:     category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Subcategories: subcategories,
		}
	}

	return usecase.GetIncomeCategoriesOutput{
		ID:         input.UserID,
		ExitCode:   common.Success,
		Message:    "Income categories retrieved successfully",
		Categories: categoriesData,
	}
}