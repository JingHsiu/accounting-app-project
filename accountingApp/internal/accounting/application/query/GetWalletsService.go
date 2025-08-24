package query

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

type GetWalletsInput struct {
	UserID string
}

type GetWalletsOutput struct {
	ID       string          `json:"id"`
	ExitCode common.ExitCode `json:"exit_code"`
	Message  string          `json:"message"`
	Wallets  []*model.Wallet `json:"wallets,omitempty"`
}

func (o GetWalletsOutput) GetID() string                { return o.ID }
func (o GetWalletsOutput) GetExitCode() common.ExitCode { return o.ExitCode }
func (o GetWalletsOutput) GetMessage() string           { return o.Message }

type GetWalletsService struct {
	walletRepo repository.WalletRepository
}

func NewGetWalletsService(walletRepo repository.WalletRepository) *GetWalletsService {
	return &GetWalletsService{walletRepo: walletRepo}
}

func (s *GetWalletsService) Execute(input GetWalletsInput) common.Output {
	wallets, err := s.walletRepo.FindByUserID(input.UserID)
	if err != nil {
		return GetWalletsOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Failed to retrieve wallets: %v", err),
		}
	}

	return GetWalletsOutput{
		ID:       input.UserID,
		ExitCode: common.Success,
		Message:  "Wallets retrieved successfully",
		Wallets:  wallets,
	}
}