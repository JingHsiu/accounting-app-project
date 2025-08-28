package query

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

type GetWalletService struct {
	walletRepo repository.WalletRepository
}

func NewGetWalletService(walletRepo repository.WalletRepository) *GetWalletService {
	return &GetWalletService{walletRepo: walletRepo}
}

func (s *GetWalletService) Execute(input usecase.GetWalletInput) common.Output {
	var wallet *model.Wallet
	var err error

	if input.IncludeTransactions {
		wallet, err = s.walletRepo.FindByIDWithTransactions(input.WalletID)
	} else {
		wallet, err = s.walletRepo.FindByID(input.WalletID)
	}

	if err != nil {
		return usecase.GetWalletOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Failed to retrieve wallet: %v", err),
		}
	}

	if wallet == nil {
		return usecase.GetWalletOutput{
			ExitCode: common.Failure,
			Message:  "Wallet not found",
		}
	}

	return usecase.GetWalletOutput{
		ID:       wallet.ID,
		ExitCode: common.Success,
		Message:  "Wallet retrieved successfully",
		Wallet:   wallet,
	}
}