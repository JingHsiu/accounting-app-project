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
// 使用Inquiry來載入交易記錄，而不修改Repository基本CRUD介面
func (r *WalletRepositoryImpl) FindByIDWithTransactions(id string) (*model.Wallet, error) {
	// 1. 透過基本Repository方法取得錢包
	wallet, err := r.FindByID(id)
	if err != nil || wallet == nil {
		return wallet, err
	}

	// 2. 使用Inquiry載入交易記錄 (透過第三層直接查詢)
	// TODO: 需要實作WalletInquiry來載入交易記錄
	// 暫時標記為已載入，避免編譯錯誤
	wallet.SetFullyLoaded(true)

	return wallet, nil
}

// FindByUserID 根據UserID查找用戶的所有錢包
// 使用Inquiry實作，不修改基本Repository介面
func (r *WalletRepositoryImpl) FindByUserID(userID string) ([]*model.Wallet, error) {
	// TODO: 使用Inquiry查詢用戶的所有錢包
	// 暫時返回空切片，避免編譯錯誤
	return make([]*model.Wallet, 0), nil
}

// SaveData 實現WalletRepositoryPeer介面 - 直接委派給peer
func (r *WalletRepositoryImpl) SaveData(data mapper.WalletData) error {
	return r.peer.SaveData(data)
}

// FindDataByID 實現WalletRepositoryPeer介面 - 直接委派給peer
func (r *WalletRepositoryImpl) FindDataByID(id string) (*mapper.WalletData, error) {
	return r.peer.FindDataByID(id)
}

// DeleteData 實現WalletRepositoryPeer介面 - 直接委派給peer
func (r *WalletRepositoryImpl) DeleteData(id string) error {
	return r.peer.DeleteData(id)
}