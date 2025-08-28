package query

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
)

type GetWalletsService struct {
	walletRepo repository.WalletRepository
}

func NewGetWalletsService(walletRepo repository.WalletRepository) *GetWalletsService {
	return &GetWalletsService{walletRepo: walletRepo}
}

func (s *GetWalletsService) Execute(input usecase.GetWalletsInput) common.Output {
	wallets, err := s.walletRepo.FindByUserID(input.UserID)
	if err != nil {
		return usecase.GetWalletsOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Failed to retrieve wallets: %v", err),
		}
	}

	return usecase.GetWalletsOutput{
		ID:       input.UserID,
		ExitCode: common.Success,
		Message:  "Wallets retrieved successfully",
		Wallets:  wallets,
	}
}