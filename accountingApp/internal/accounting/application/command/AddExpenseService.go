package command

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
	"time"
)

type AddExpenseInput struct {
	WalletID      string
	SubcategoryID string
	Amount        int64
	Currency      string
	Description   string
	Date          time.Time
}

type AddExpenseService struct {
	walletRepo   repository.WalletRepository
	categoryRepo repository.ExpenseCategoryRepository
}

func NewAddExpenseService(walletRepo repository.WalletRepository, categoryRepo repository.ExpenseCategoryRepository) *AddExpenseService {
	return &AddExpenseService{
		walletRepo:   walletRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *AddExpenseService) Execute(input AddExpenseInput) common.Output {
	// 1. 透過Repository取得錢包 (可能需要完整聚合取決於業務需求)
	wallet, err := s.walletRepo.FindByIDWithTransactions(input.WalletID)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("wallet not found: %v", err),
		}
	}

	// 2. 驗證分類存在 (透過Repository而非Inquiry)
	category, err := s.categoryRepo.FindBySubcategoryID(input.SubcategoryID)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("category not found: %v", err),
		}
	}

	// 透過聚合根驗證子分類存在性
	err = category.ValidateSubcategoryExists(input.SubcategoryID)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("invalid subcategory: %v", err),
		}
	}
	
	// 3. 建立金額物件
	amount, err := model.NewMoney(input.Amount, input.Currency)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("invalid amount: %v", err),
		}
	}

	// 4. 透過Domain Model執行業務邏輯
	expense, err := wallet.AddExpense(*amount, input.SubcategoryID, input.Description, input.Date)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("failed to add expense: %v", err),
		}
	}

	// 5. 儲存完整聚合 (包含新的交易記錄)
	if err := s.walletRepo.Save(wallet); err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("failed to save wallet: %v", err),
		}
	}

	return common.UseCaseOutput{
		ID:       expense.ID,
		ExitCode: common.Success,
		Message:  "Expense added successfully",
	}
}
