package command

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

type CreateExpenseCategoryInput struct {
	UserID string
	Name   string
}

type CreateExpenseCategoryService struct {
	repo repository.ExpenseCategoryRepository
}

func NewCreateExpenseCategoryService(repo repository.ExpenseCategoryRepository) *CreateExpenseCategoryService {
	return &CreateExpenseCategoryService{repo: repo}
}

func (s *CreateExpenseCategoryService) Execute(input CreateExpenseCategoryInput) common.Output {
	categoryName, err := model.NewCategoryName(input.Name)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Invalid category name: %v", err),
		}
	}

	category, err := model.NewExpenseCategory(input.UserID, *categoryName)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Creating expense category failed: %v", err),
		}
	}

	err = s.repo.Save(category)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Saving expense category failed: %v", err),
		}
	}

	return common.UseCaseOutput{
		ID:       category.ID,
		ExitCode: common.Success,
	}
}
