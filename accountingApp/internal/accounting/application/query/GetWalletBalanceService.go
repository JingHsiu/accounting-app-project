package query

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
)

type GetWalletBalanceService struct {
	walletRepo repository.WalletRepository
}

func NewGetWalletBalanceService(walletRepo repository.WalletRepository) *GetWalletBalanceService {
	return &GetWalletBalanceService{walletRepo: walletRepo}
}

func (s *GetWalletBalanceService) Execute(input usecase.GetWalletBalanceInput) common.Output {
	// 只需要基本資訊，不需要載入所有交易記錄 (效能優化)
	wallet, err := s.walletRepo.FindByID(input.WalletID)
	if err != nil {
		return usecase.GetWalletBalanceOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("wallet not found: %v", err),
		}
	}

	return usecase.GetWalletBalanceOutput{
		ID:       wallet.ID,
		ExitCode: common.Success,
		Message:  "Balance retrieved successfully",
		Balance:  fmt.Sprintf("%.2f", float64(wallet.Balance.Amount)/100),
		Currency: wallet.Balance.Currency,
	}
}
