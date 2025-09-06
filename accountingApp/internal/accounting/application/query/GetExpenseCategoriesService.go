package query

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
)

type GetExpenseCategoriesService struct {
	expenseCategoryRepo repository.ExpenseCategoryRepository
}

func NewGetExpenseCategoriesService(expenseCategoryRepo repository.ExpenseCategoryRepository) *GetExpenseCategoriesService {
	return &GetExpenseCategoriesService{expenseCategoryRepo: expenseCategoryRepo}
}

func (s *GetExpenseCategoriesService) Execute(input usecase.GetExpenseCategoriesInput) common.Output {
	// First, try to get categories for the user
	categories, err := s.expenseCategoryRepo.FindByUserID(input.UserID)
	if err != nil {
		return usecase.GetExpenseCategoriesOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Failed to retrieve expense categories: %v", err),
		}
	}

	// If no categories found, the user needs to have default categories initialized
	if len(categories) == 0 {
		return usecase.GetExpenseCategoriesOutput{
			ID:         input.UserID,
			ExitCode:   common.Success,
			Message:    "No expense categories found. Please initialize default categories.",
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
			Type:          "expense",
			CreatedAt:     category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Subcategories: subcategories,
		}
	}

	return usecase.GetExpenseCategoriesOutput{
		ID:         input.UserID,
		ExitCode:   common.Success,
		Message:    "Expense categories retrieved successfully",
		Categories: categoriesData,
	}
}