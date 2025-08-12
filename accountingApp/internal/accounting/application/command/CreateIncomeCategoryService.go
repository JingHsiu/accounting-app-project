package command

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

type CreateIncomeCategoryInput struct {
	UserID string
	Name   string
}

type CreateIncomeCategoryService struct {
	repo repository.IncomeCategoryRepository
}

func NewCreateIncomeCategoryService(repo repository.IncomeCategoryRepository) *CreateIncomeCategoryService {
	return &CreateIncomeCategoryService{repo: repo}
}

func (s *CreateIncomeCategoryService) Execute(input CreateIncomeCategoryInput) common.Output {
	categoryName, err := model.NewCategoryName(input.Name)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Invalid category name: %v", err),
		}
	}

	category, err := model.NewIncomeCategory(input.UserID, *categoryName)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Creating income category failed: %v", err),
		}
	}

	err = s.repo.Save(category)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Saving income category failed: %v", err),
		}
	}

	return common.UseCaseOutput{
		ID:       category.ID,
		ExitCode: common.Success,
	}
}
