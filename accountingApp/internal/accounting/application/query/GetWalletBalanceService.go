package query

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
)

type GetWalletBalanceInput struct {
	WalletID string
}

type GetWalletBalanceOutput struct {
	ID       string          `json:"id"`
	ExitCode common.ExitCode `json:"exit_code"`
	Message  string          `json:"message"`
	Balance  string          `json:"balance,omitempty"`
	Currency string          `json:"currency,omitempty"`
}

func (o GetWalletBalanceOutput) GetID() string                { return o.ID }
func (o GetWalletBalanceOutput) GetExitCode() common.ExitCode { return o.ExitCode }
func (o GetWalletBalanceOutput) GetMessage() string           { return o.Message }

type GetWalletBalanceService struct {
	walletRepo repository.WalletRepository
}

func NewGetWalletBalanceService(walletRepo repository.WalletRepository) *GetWalletBalanceService {
	return &GetWalletBalanceService{walletRepo: walletRepo}
}

func (s *GetWalletBalanceService) Execute(input GetWalletBalanceInput) common.Output {
	// 只需要基本資訊，不需要載入所有交易記錄 (效能優化)
	wallet, err := s.walletRepo.FindByID(input.WalletID)
	if err != nil {
		return GetWalletBalanceOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("wallet not found: %v", err),
		}
	}

	return GetWalletBalanceOutput{
		ID:       wallet.ID,
		ExitCode: common.Success,
		Message:  "Balance retrieved successfully",
		Balance:  fmt.Sprintf("%.2f", float64(wallet.Balance.Amount)/100),
		Currency: wallet.Balance.Currency,
	}
}
