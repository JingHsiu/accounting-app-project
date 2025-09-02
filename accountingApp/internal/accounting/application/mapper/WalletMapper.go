package mapper

import (
	"time"
	"github.com/JingHsiu/accountingApp/internal/accounting/adapter/store"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

// WalletData Wallet的持久化資料結構
type WalletData struct {
	ID              string    `db:"id"`
	UserID          string    `db:"user_id"`
	Name            string    `db:"name"`
	Type            string    `db:"type"`
	Currency        string    `db:"currency"`
	BalanceAmount   int64     `db:"balance_amount"`
	BalanceCurrency string    `db:"balance_currency"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

func (wd WalletData) GetID() string {
	return wd.ID
}

// WalletMapper Wallet聚合的資料轉換器
type WalletMapper struct{}

func NewWalletMapper() *WalletMapper {
	return &WalletMapper{}
}

// ToData 將Wallet Domain Model轉換為WalletData
func (m *WalletMapper) ToData(wallet *model.Wallet) WalletData {
	return WalletData{
		ID:              wallet.ID,
		UserID:          wallet.UserID,
		Name:            wallet.Name,
		Type:            string(wallet.Type),
		Currency:        wallet.Currency(),
		BalanceAmount:   wallet.Balance.Amount,
		BalanceCurrency: wallet.Balance.Currency,
		CreatedAt:       wallet.CreatedAt,
		UpdatedAt:       wallet.UpdatedAt,
	}
}

// ToDomain 將WalletData轉換為Wallet Domain Model
func (m *WalletMapper) ToDomain(data WalletData) (*model.Wallet, error) {
	walletType, err := model.ParseWalletType(data.Type)
	if err != nil {
		return nil, err
	}
	
	balance, err := model.NewMoney(data.BalanceAmount, data.BalanceCurrency)
	if err != nil {
		return nil, err
	}
	
	return &model.Wallet{
		ID:        data.ID,
		UserID:    data.UserID,
		Name:      data.Name,
		Type:      walletType,
		Balance:   *balance,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}, nil
}

// 確保WalletData實現AggregateData介面
var _ store.AggregateData = (*WalletData)(nil)

// 確保WalletMapper實現Mapper介面和AggregateMapper介面
var _ Mapper[*model.Wallet, WalletData] = (*WalletMapper)(nil)
var _ store.AggregateMapper[*model.Wallet, WalletData] = (*WalletMapper)(nil)