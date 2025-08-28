package command

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

type CreateWalletService struct {
	repo repository.WalletRepository
}

func NewCreateWalletService(repo repository.WalletRepository) *CreateWalletService {
	return &CreateWalletService{repo: repo}
}

func (s *CreateWalletService) Execute(input usecase.CreateWalletInput) common.Output {
	parsedType, err := model.ParseWalletType(input.Type)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("ParseWalletType failed: %v", err),
		}
	}

	// Determine initial balance amount
	var initialBalanceAmount int64 = 0
	if input.InitialBalance != nil {
		initialBalanceAmount = *input.InitialBalance
	}

	wallet, err := model.NewWalletWithInitialBalance(input.UserID, input.Name, parsedType, input.Currency, initialBalanceAmount)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Creating wallet failed: %v", err),
		}
	}

	err = s.repo.Save(wallet)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Saving wallet failed: %v", err),
		}
	}

	return common.UseCaseOutput{
		ID:       wallet.ID,
		ExitCode: common.Success,
	}
}
