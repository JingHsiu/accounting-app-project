package query

import (
	"fmt"
	"strings"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
	"time"
)

type GetExpensesService struct {
	walletRepo repository.WalletRepository
}

func NewGetExpensesService(walletRepo repository.WalletRepository) *GetExpensesService {
	return &GetExpensesService{
		walletRepo: walletRepo,
	}
}

func (s *GetExpensesService) Execute(input usecase.GetExpensesInput) common.Output {
	// Get user's wallets to extract expense records
	wallets, err := s.walletRepo.FindByUserID(input.UserID)
	if err != nil {
		return usecase.GetExpensesOutput{
			ID:       input.UserID,
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Failed to retrieve wallets: %v", err),
		}
	}

	if len(wallets) == 0 {
		return usecase.GetExpensesOutput{
			ID:       input.UserID,
			ExitCode: common.Success,
			Message:  "No wallets found. Please create a wallet first.",
			Data:     []usecase.ExpenseRecordData{},
			Count:    0,
		}
	}

	// Collect expense records from all wallets
	allExpenseRecords := make([]usecase.ExpenseRecordData, 0)

	for _, wallet := range wallets {
		// Load wallet with transactions to get complete aggregate
		fullyLoadedWallet, err := s.walletRepo.FindByIDWithTransactions(wallet.ID)
		if err != nil {
			// Fallback to basic wallet if transaction loading fails
			fullyLoadedWallet = wallet
		}
		
		// Get expense records from fully loaded wallet aggregate
		expenseRecords := fullyLoadedWallet.GetExpenseRecords()
		
		for i := range expenseRecords {
			record := &expenseRecords[i]
			// Apply filters
			if input.WalletID != nil && *input.WalletID != record.WalletID {
				continue
			}
			if input.CategoryID != nil && *input.CategoryID != record.SubcategoryID {
				continue
			}
			if input.StartDate != nil && record.Date.Before(*input.StartDate) {
				continue
			}
			if input.EndDate != nil && record.Date.After(*input.EndDate) {
				continue
			}
			if input.MinAmount != nil && record.Amount.Amount < *input.MinAmount {
				continue
			}
			if input.MaxAmount != nil && record.Amount.Amount > *input.MaxAmount {
				continue
			}
			if input.Description != nil && *input.Description != "" {
				// Simple contains check for description filter
				// In production, you might want more sophisticated text search
				descriptionFilter := *input.Description
				if len(record.Description) == 0 || 
				   (len(record.Description) > 0 && len(descriptionFilter) > 0 && 
				    !strings.Contains(record.Description, descriptionFilter)) {
					continue
				}
			}

			// Convert to API format
			expenseData := usecase.ExpenseRecordData{
				ID:            record.ID,
				WalletID:      record.WalletID,
				SubcategoryID: record.SubcategoryID,
				Amount: struct {
					Amount   int64  `json:"amount"`
					Currency string `json:"currency"`
				}{
					Amount:   record.Amount.Amount,
					Currency: record.Amount.Currency,
				},
				Description: record.Description,
				Date:        record.Date.Format(time.RFC3339),
				CreatedAt:   record.CreatedAt.Format(time.RFC3339),
			}

			allExpenseRecords = append(allExpenseRecords, expenseData)
		}
	}

	return usecase.GetExpensesOutput{
		ID:       input.UserID,
		ExitCode: common.Success,
		Message:  fmt.Sprintf("Successfully retrieved %d expense records", len(allExpenseRecords)),
		Data:     allExpenseRecords,
		Count:    len(allExpenseRecords),
	}
}

