package command

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

type UpdateWalletService struct {
	repo repository.WalletRepository
}

func NewUpdateWalletService(repo repository.WalletRepository) *UpdateWalletService {
	return &UpdateWalletService{repo: repo}
}

func (s *UpdateWalletService) Execute(input usecase.UpdateWalletInput) common.Output {
	// Get existing wallet
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

	// Update wallet properties using domain model methods
	updated := false

	if input.Name != nil && *input.Name != wallet.Name {
		if err := wallet.UpdateName(*input.Name); err != nil {
			return common.UseCaseOutput{
				ExitCode: common.Failure,
				Message:  fmt.Sprintf("Invalid wallet name: %v", err),
			}
		}
		updated = true
	}

	if input.Type != nil && *input.Type != string(wallet.Type) {
		walletType, err := model.ParseWalletType(*input.Type)
		if err != nil {
			return common.UseCaseOutput{
				ExitCode: common.Failure,
				Message:  fmt.Sprintf("Invalid wallet type: %v", err),
			}
		}
		if err := wallet.UpdateType(walletType); err != nil {
			return common.UseCaseOutput{
				ExitCode: common.Failure,
				Message:  fmt.Sprintf("Failed to update wallet type: %v", err),
			}
		}
		updated = true
	}

	// Note: Currency update is intentionally excluded as it would require complex balance conversion

	if updated {
		if err := s.repo.Save(wallet); err != nil {
			return common.UseCaseOutput{
				ExitCode: common.Failure,
				Message:  fmt.Sprintf("Failed to update wallet: %v", err),
			}
		}
	}

	return common.UseCaseOutput{
		ID:       wallet.ID,
		ExitCode: common.Success,
		Message:  "Wallet updated successfully",
	}
}