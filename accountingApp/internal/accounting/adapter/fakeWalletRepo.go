package adapter

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/mapper"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

type fakeWalletRepo struct {
	data map[string]*model.Wallet
}

func NewFakeWalletRepo() (*fakeWalletRepo, error) {
	return &fakeWalletRepo{data: make(map[string]*model.Wallet)}, nil
}

func (f *fakeWalletRepo) Save(wallet *model.Wallet) error {
	f.data[wallet.ID] = wallet
	return nil
}

func (f *fakeWalletRepo) FindByID(id string) (*model.Wallet, error) {
	wallet, ok := f.data[id]
	if !ok {
		return nil, fmt.Errorf("wallet not found")
	}
	return wallet, nil
}

func (f *fakeWalletRepo) Delete(id string) error {
	delete(f.data, id)
	return nil
}

// 實現WalletRepositoryPeer介面的Bridge Pattern方法
func (f *fakeWalletRepo) SaveData(data mapper.WalletData) error {
	// Convert data back to domain model for storage
	money, _ := model.NewMoney(data.BalanceAmount, data.BalanceCurrency)
	wallet := &model.Wallet{
		ID:        data.ID,
		UserID:    data.UserID,
		Name:      data.Name,
		Type:      model.WalletType(data.Type),
		Balance:   *money,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
	return f.Save(wallet)
}

func (f *fakeWalletRepo) FindDataByID(id string) (*mapper.WalletData, error) {
	wallet, err := f.FindByID(id)
	if err != nil {
		return nil, err
	}
	// Convert domain model back to data structure
	mapper := mapper.NewWalletMapper()
	data := mapper.ToData(wallet)
	return &data, nil
}

func (f *fakeWalletRepo) DeleteData(id string) error {
	return f.Delete(id)
}
