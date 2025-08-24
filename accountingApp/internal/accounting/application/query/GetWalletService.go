package query

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

type GetWalletInput struct {
	WalletID           string
	IncludeTransactions bool
}

type GetWalletOutput struct {
	ID       string          `json:"id"`
	ExitCode common.ExitCode `json:"exit_code"`
	Message  string          `json:"message"`
	Wallet   *model.Wallet   `json:"wallet,omitempty"`
}

func (o GetWalletOutput) GetID() string                { return o.ID }
func (o GetWalletOutput) GetExitCode() common.ExitCode { return o.ExitCode }
func (o GetWalletOutput) GetMessage() string           { return o.Message }

type GetWalletService struct {
	walletRepo repository.WalletRepository
}

func NewGetWalletService(walletRepo repository.WalletRepository) *GetWalletService {
	return &GetWalletService{walletRepo: walletRepo}
}

func (s *GetWalletService) Execute(input GetWalletInput) common.Output {
	var wallet *model.Wallet
	var err error

	if input.IncludeTransactions {
		wallet, err = s.walletRepo.FindByIDWithTransactions(input.WalletID)
	} else {
		wallet, err = s.walletRepo.FindByID(input.WalletID)
	}

	if err != nil {
		return GetWalletOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("Failed to retrieve wallet: %v", err),
		}
	}

	if wallet == nil {
		return GetWalletOutput{
			ExitCode: common.Failure,
			Message:  "Wallet not found",
		}
	}

	return GetWalletOutput{
		ID:       wallet.ID,
		ExitCode: common.Success,
		Message:  "Wallet retrieved successfully",
		Wallet:   wallet,
	}
}