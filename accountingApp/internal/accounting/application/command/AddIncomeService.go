package command

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
	"time"
)

type AddIncomeInput struct {
	WalletID      string
	SubcategoryID string
	Amount        int64
	Currency      string
	Description   string
	Date          time.Time
}

type AddIncomeService struct {
	walletRepo   repository.WalletRepository
	categoryRepo repository.IncomeCategoryRepository
}

func NewAddIncomeService(walletRepo repository.WalletRepository, categoryRepo repository.IncomeCategoryRepository) *AddIncomeService {
	return &AddIncomeService{
		walletRepo:   walletRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *AddIncomeService) Execute(input AddIncomeInput) common.Output {
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

	// 2. 驗證子分類存在 (透過聚合根驗證)
	// TODO: 完整實現需要 IncomeCategoryRepository
	if s.categoryRepo != nil {
		// 根據子分類ID找到包含它的分類聚合
		category, err := s.categoryRepo.FindBySubcategoryID(input.SubcategoryID)
		if err != nil {
			return common.UseCaseOutput{
				ExitCode: common.Failure,
				Message:  fmt.Sprintf("Failed to find category for subcategory: %v", err),
			}
		}
		if category == nil {
			return common.UseCaseOutput{
				ExitCode: common.Failure,
				Message:  "Subcategory not found in any category",
			}
		}

		// 透過聚合根驗證子分類存在性
		err = category.ValidateSubcategoryExists(input.SubcategoryID)
		if err != nil {
			return common.UseCaseOutput{
				ExitCode: common.Failure,
				Message:  fmt.Sprintf("Invalid subcategory: %v", err),
			}
		}
	}
	// 如果 categoryRepo 為 nil，跳過驗證（用於測試或簡化場景）

	// 3. 建立金額 Value Object
	amount, err := model.NewMoney(input.Amount, input.Currency)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Invalid amount: %v", err),
		}
	}

	// 4. 透過錢包聚合根新增收入
	income, err := wallet.AddIncome(*amount, input.SubcategoryID, input.Description, input.Date)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Adding income failed: %v", err),
		}
	}

	// 5. 持久化錢包聚合
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
