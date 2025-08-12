package service

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

type CategoryValidationService struct{}

func NewCategoryValidationService() *CategoryValidationService {
	return &CategoryValidationService{}
}

// ValidateExpenseSubcategory 驗證支出子分類是否存在於指定的分類聚合中
func (s *CategoryValidationService) ValidateExpenseSubcategory(category *model.ExpenseCategory, subcategoryID string) error {
	if category == nil {
		return fmt.Errorf("category cannot be nil")
	}

	return category.ValidateSubcategoryExists(subcategoryID)
}

// ValidateIncomeSubcategory 驗證收入子分類是否存在於指定的分類聚合中
func (s *CategoryValidationService) ValidateIncomeSubcategory(category *model.IncomeCategory, subcategoryID string) error {
	if category == nil {
		return fmt.Errorf("category cannot be nil")
	}

	return category.ValidateSubcategoryExists(subcategoryID)
}

// FindExpenseSubcategory 在分類聚合中查找指定的子分類
func (s *CategoryValidationService) FindExpenseSubcategory(category *model.ExpenseCategory, subcategoryID string) (*model.ExpenseSubcategory, error) {
	if category == nil {
		return nil, fmt.Errorf("category cannot be nil")
	}

	return category.GetSubcategory(subcategoryID)
}

// FindIncomeSubcategory 在分類聚合中查找指定的子分類
func (s *CategoryValidationService) FindIncomeSubcategory(category *model.IncomeCategory, subcategoryID string) (*model.IncomeSubcategory, error) {
	if category == nil {
		return nil, fmt.Errorf("category cannot be nil")
	}

	return category.GetSubcategory(subcategoryID)
}
