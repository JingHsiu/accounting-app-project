package query

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
	"time"
)

type GetIncomesService struct {
	walletRepo repository.WalletRepository
}

func NewGetIncomesService(walletRepo repository.WalletRepository) *GetIncomesService {
	return &GetIncomesService{walletRepo: walletRepo}
}

func (s *GetIncomesService) Execute(input usecase.GetIncomesInput) common.Output {
	// Get user's wallets to extract income records
	wallets, err := s.walletRepo.FindByUserID(input.UserID)
	if err != nil {
		return usecase.GetIncomesOutput{
			ID:       input.UserID,
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Failed to retrieve wallets: %v", err),
		}
	}

	if len(wallets) == 0 {
		return usecase.GetIncomesOutput{
			ID:       input.UserID,
			ExitCode: common.Success,
			Message:  "No wallets found. Please create a wallet first.",
			Data:     []usecase.IncomeRecordData{},
			Count:    0,
		}
	}

	// Collect income records from all wallets
	var allIncomeRecords []usecase.IncomeRecordData

	for _, wallet := range wallets {
		// Get income records for this wallet
		incomeRecords := wallet.GetIncomeRecords()
		
		for _, record := range incomeRecords {
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
				    !contains(record.Description, descriptionFilter)) {
					continue
				}
			}

			// Convert to API format
			incomeData := usecase.IncomeRecordData{
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

			allIncomeRecords = append(allIncomeRecords, incomeData)
		}
	}

	return usecase.GetIncomesOutput{
		ID:       input.UserID,
		ExitCode: common.Success,
		Message:  fmt.Sprintf("Successfully retrieved %d income records", len(allIncomeRecords)),
		Data:     allIncomeRecords,
		Count:    len(allIncomeRecords),
	}
}

// Helper function for simple string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (substr == "" || findInString(s, substr))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}