package command

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

type AddIncomeService struct {
	walletRepo repository.WalletRepository
}

func NewAddIncomeService(walletRepo repository.WalletRepository) *AddIncomeService {
	return &AddIncomeService{
		walletRepo: walletRepo,
	}
}

func (s *AddIncomeService) Execute(input usecase.AddIncomeInput) common.Output {
	// 1. 驗證錢包存在
	wallet, err := s.walletRepo.FindByID(input.WalletID)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Failed to find wallet: %v", err),
		}
	}
	if wallet == nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  "Wallet not found",
		}
	}

	// 2. 建立金額 Value Object
	amount, err := model.NewMoney(input.Amount, input.Currency)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Invalid amount: %v", err),
		}
	}

	// 3. 透過錢包聚合根新增收入
	income, err := wallet.AddIncome(*amount, input.SubcategoryID, input.Description, input.Date)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Adding income failed: %v", err),
		}
	}

	// 4. 持久化錢包聚合 (包括新增的收入記錄)
	err = s.walletRepo.Save(wallet)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Saving wallet failed: %v", err),
		}
	}

	return common.UseCaseOutput{
		ID:       income.ID,
		ExitCode: common.Success,
	}
}
