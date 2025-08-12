package command

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
	"time"
)

type ProcessTransferInput struct {
	FromWalletID string    // 來源錢包ID
	ToWalletID   string    // 目標錢包ID
	Amount       int64     // 轉帳金額 (cents)
	Currency     string    // 貨幣
	Fee          int64     // 手續費 (cents)
	Description  string    // 描述
	Date         time.Time // 轉帳日期
}

// ProcessTransferService - 只依賴Repository
type ProcessTransferService struct {
	walletRepo repository.WalletRepository
}

func NewProcessTransferService(walletRepo repository.WalletRepository) *ProcessTransferService {
	return &ProcessTransferService{
		walletRepo: walletRepo,
	}
}

func (s *ProcessTransferService) Execute(input ProcessTransferInput) common.Output {
	// 1. 取得兩個錢包 (載入完整聚合)
	fromWallet, err := s.walletRepo.FindByIDWithTransactions(input.FromWalletID)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("from wallet not found: %v", err),
		}
	}

	toWallet, err := s.walletRepo.FindByIDWithTransactions(input.ToWalletID)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("to wallet not found: %v", err),
		}
	}

	// 2. 建立金額物件
	amount, err := model.NewMoney(input.Amount, input.Currency)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("invalid amount: %v", err),
		}
	}

	fee, err := model.NewMoney(input.Fee, input.Currency)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("invalid fee: %v", err),
		}
	}

	// 3. 透過Domain Model處理轉帳 (雙邊操作)
	err = fromWallet.ProcessOutgoingTransfer(*amount, *fee)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("transfer failed: %v", err),
		}
	}

	err = toWallet.ProcessIncomingTransfer(*amount)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("transfer failed: %v", err),
		}
	}

	// 4. 建立轉帳記錄 (在來源錢包中)
	transfer, err := fromWallet.CreateTransfer(input.ToWalletID, *amount, *fee, input.Description, input.Date)
	if err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("failed to create transfer record: %v", err),
		}
	}

	// 5. 儲存兩個錢包 (應該在同一個資料庫交易中)
	// TODO: 實作交易管理 (Transaction Manager)
	if err := s.walletRepo.Save(fromWallet); err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("failed to save from wallet: %v", err),
		}
	}

	if err := s.walletRepo.Save(toWallet); err != nil {
		return common.UseCaseOutput{
			ExitCode: common.Failure,
			Message:  fmt.Sprintf("failed to save to wallet: %v", err),
		}
	}

	return common.UseCaseOutput{
		ID:       transfer.ID,
		ExitCode: common.Success,
		Message:  "Transfer processed successfully",
	}
}