package command

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
)

type DeleteWalletInput struct {
	WalletID string
}

type DeleteWalletService struct {
	repo repository.WalletRepository
}

func NewDeleteWalletService(repo repository.WalletRepository) *DeleteWalletService {
	return &DeleteWalletService{repo: repo}
}

func (s *DeleteWalletService) Execute(input DeleteWalletInput) common.Output {
	// Check if wallet exists
	wallet, err := s.repo.FindByID(input.WalletID)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Failed to retrieve wallet: %v", err),
		}
	}

	if wallet == nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  "Wallet not found",
		}
	}

	// Delete the wallet
	if err := s.repo.Delete(input.WalletID); err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Failed to delete wallet: %v", err),
		}
	}

	return common.UseCaseOutput{
		ID:       input.WalletID,
		ExitCode: common.Success,
		Message:  "Wallet deleted successfully",
	}
}