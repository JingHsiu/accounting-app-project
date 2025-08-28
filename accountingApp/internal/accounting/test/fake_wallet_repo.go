package test

import (
	"github.com/JingHsiu/accountingApp/internal/accounting/application/mapper"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

type FakeWalletRepo struct {
	data map[string]*model.Wallet
}

func NewFakeWalletRepo() (*FakeWalletRepo, error) {
	return &FakeWalletRepo{data: make(map[string]*model.Wallet)}, nil
}

func (f *FakeWalletRepo) Save(wallet *model.Wallet) error {
	f.data[wallet.ID] = wallet
	return nil
}

func (f *FakeWalletRepo) FindByID(id string) (*model.Wallet, error) {
	wallet, ok := f.data[id]
	if !ok {
		// Repository pattern - return (nil, nil) for "not found" vs (nil, error) for actual errors
		return nil, nil
	}
	return wallet, nil
}

func (f *FakeWalletRepo) Delete(id string) error {
	delete(f.data, id)
	return nil
}

// 實現WalletRepositoryPeer介面的Bridge Pattern方法
func (f *FakeWalletRepo) SaveData(data mapper.WalletData) error {
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

func (f *FakeWalletRepo) FindDataByID(id string) (*mapper.WalletData, error) {
	wallet, err := f.FindByID(id)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, nil
	}
	// Convert domain model back to data structure
	mapper := mapper.NewWalletMapper()
	data := mapper.ToData(wallet)
	return &data, nil
}

func (f *FakeWalletRepo) DeleteData(id string) error {
	return f.Delete(id)
}

// Add missing methods required by WalletRepository interface
func (f *FakeWalletRepo) FindByIDWithTransactions(id string) (*model.Wallet, error) {
	wallet, err := f.FindByID(id)
	if err != nil {
		return nil, err
	}
	if wallet != nil {
		wallet.SetFullyLoaded(true)
	}
	return wallet, nil
}

func (f *FakeWalletRepo) FindByUserID(userID string) ([]*model.Wallet, error) {
	var wallets []*model.Wallet
	for _, wallet := range f.data {
		if wallet.UserID == userID {
			wallets = append(wallets, wallet)
		}
	}
	return wallets, nil
}

func (f *FakeWalletRepo) FindDataByUserID(userID string) ([]mapper.WalletData, error) {
	wallets, err := f.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var dataList []mapper.WalletData
	walletMapper := mapper.NewWalletMapper()
	for _, wallet := range wallets {
		data := walletMapper.ToData(wallet)
		dataList = append(dataList, data)
	}
	return dataList, nil
}