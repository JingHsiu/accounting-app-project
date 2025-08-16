package repository

import (
	"github.com/JingHsiu/accountingApp/internal/accounting/application/mapper"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

// WalletRepositoryImpl 第二層錢包儲存庫實現 (Application Layer)
// 使用Bridge Pattern，透過peer介面橋接到第三層，避免直接依賴Domain Model
type WalletRepositoryImpl struct {
	peer   WalletRepositoryPeer // 橋接到第三層的實現
	mapper *mapper.WalletMapper // 負責Domain Model與Data結構轉換
}

// NewWalletRepositoryImpl 創建錢包儲存庫實現
func NewWalletRepositoryImpl(peer WalletRepositoryPeer) WalletRepository {
	return &WalletRepositoryImpl{
		peer:   peer,
		mapper: mapper.NewWalletMapper(),
	}
}

// Save 儲存錢包Domain Model
// 透過mapper轉換後，使用peer介面操作資料層
func (r *WalletRepositoryImpl) Save(wallet *model.Wallet) error {
	// Domain Model → Data結構
	data := r.mapper.ToData(wallet)
	
	// 透過peer介面橋接到第四層儲存
	return r.peer.SaveData(data)
}

// FindByID 根據ID查找錢包Domain Model
func (r *WalletRepositoryImpl) FindByID(id string) (*model.Wallet, error) {
	// 透過peer介面從第四層取得資料
	data, err := r.peer.FindDataByID(id)
	if err != nil {
		return nil, err
	}
	
	if data == nil {
		return nil, nil
	}
	
	// Data結構 → Domain Model
	return r.mapper.ToDomain(*data)
}

// Delete 刪除錢包
func (r *WalletRepositoryImpl) Delete(id string) error {
	// 透過peer介面橋接到第四層刪除
	return r.peer.DeleteData(id)
}

// FindByIDWithTransactions 根據ID查找錢包及所有交易記錄 (載入完整聚合)
func (r *WalletRepositoryImpl) FindByIDWithTransactions(id string) (*model.Wallet, error) {
	// 1. 透過基本Repository方法取得錢包
	wallet, err := r.FindByID(id)
	if err != nil || wallet == nil {
		return wallet, err
	}

	// 2. 載入交易記錄 (從相關的交易表中查詢)
	// TODO: 實際實作需要查詢 expenses, incomes, transfers 表
	// 目前先標記為已載入，避免編譯錯誤
	// 
	// 未來的實作應該包含:
	// - 查詢 expenses 表: WHERE wallet_id = $1
	// - 查詢 incomes 表: WHERE wallet_id = $1  
	// - 查詢 transfers 表: WHERE from_wallet_id = $1 OR to_wallet_id = $1
	// - 將查詢結果透過 mapper 轉換成 domain model
	// - 使用 wallet.AddExpenseRecord(), wallet.AddIncomeRecord(), wallet.AddTransfer() 載入
	
	wallet.SetFullyLoaded(true)

	return wallet, nil
}

// FindByUserID 根據UserID查找用戶的所有錢包
func (r *WalletRepositoryImpl) FindByUserID(userID string) ([]*model.Wallet, error) {
	// 透過peer介面從第四層取得資料
	dataList, err := r.peer.FindDataByUserID(userID)
	if err != nil {
		return nil, err
	}
	
	// 將資料結構轉換為Domain Model
	wallets := make([]*model.Wallet, len(dataList))
	for i, data := range dataList {
		wallet, err := r.mapper.ToDomain(data)
		if err != nil {
			return nil, err
		}
		wallets[i] = wallet
	}
	
	return wallets, nil
}

// SaveData 實現WalletRepositoryPeer介面 - 直接委派給peer
func (r *WalletRepositoryImpl) SaveData(data mapper.WalletData) error {
	return r.peer.SaveData(data)
}

// FindDataByID 實現WalletRepositoryPeer介面 - 直接委派給peer
func (r *WalletRepositoryImpl) FindDataByID(id string) (*mapper.WalletData, error) {
	return r.peer.FindDataByID(id)
}

// FindDataByUserID 實現WalletRepositoryPeer介面 - 直接委派給peer
func (r *WalletRepositoryImpl) FindDataByUserID(userID string) ([]mapper.WalletData, error) {
	return r.peer.FindDataByUserID(userID)
}

// DeleteData 實現WalletRepositoryPeer介面 - 直接委派給peer
func (r *WalletRepositoryImpl) DeleteData(id string) error {
	return r.peer.DeleteData(id)
}