package repository

import (
	"github.com/JingHsiu/accountingApp/internal/accounting/application/mapper"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

// WalletRepositoryImpl Layer 2 (Application) 錢包儲存庫實現
type WalletRepositoryImpl struct {
	peer   WalletRepositoryPeer // 橋接到Layer 3的實現
	mapper *mapper.WalletMapper // AggregateMapper：Domain ↔ Data轉換
}

// NewWalletRepositoryImpl 創建錢包儲存庫實現
func NewWalletRepositoryImpl(peer WalletRepositoryPeer) WalletRepository {
	return &WalletRepositoryImpl{
		peer:   peer,
		mapper: mapper.NewWalletMapper(),
	}
}

// Save 儲存錢包Domain Model (狀態源模式)
// 使用AggregateMapper轉換聚合狀態，透過peer介面橋接到AggregateStore
func (r *WalletRepositoryImpl) Save(wallet *model.Wallet) error {
	// 使用AggregateMapper: Domain Aggregate → AggregateData
	aggregateData := r.mapper.ToData(wallet)

	// 透過peer介面橋接到Layer 3 → Layer 4的AggregateStore
	err := r.peer.Save(aggregateData)
	if err != nil {
		// TODO: 包裝為Repository專用異常
		return err
	}

	// 清除領域事件（如果將來添加事件源）
	// wallet.ClearDomainEvents()

	return nil
}

// FindByID 根據ID查找錢包Domain Model
func (r *WalletRepositoryImpl) FindByID(id string) (*model.Wallet, error) {
	// 透過peer介面從AggregateStore取得聚合狀態
	aggregateData, err := r.peer.FindByID(id)
	if err != nil {
		return nil, err
	}

	if aggregateData == nil {
		return nil, nil // 聚合不存在
	}

	// 使用AggregateMapper: AggregateData → Domain Aggregate
	return r.mapper.ToDomain(*aggregateData)
}

// Delete 刪除錢包
func (r *WalletRepositoryImpl) Delete(id string) error {
	// 透過peer介面橋接到AggregateStore刪除聚合狀態
	return r.peer.Delete(id)
}

// FindByIDWithTransactions 根據ID查找錢包及所有交易記錄 (載入完整聚合)
func (r *WalletRepositoryImpl) FindByIDWithTransactions(id string) (*model.Wallet, error) {
	// 透過peer介面載入完整聚合（包含所有子實體）
	aggregateData, err := r.peer.FindByIDWithChildEntities(id)
	if err != nil {
		return nil, err
	}

	if aggregateData == nil {
		return nil, nil // 聚合不存在
	}

	// 使用AggregateMapper: AggregateData → Domain Aggregate (含子實體)
	return r.mapper.ToDomain(*aggregateData)
}

// FindByUserID 根據UserID查找用戶的所有錢包
func (r *WalletRepositoryImpl) FindByUserID(userID string) ([]*model.Wallet, error) {
	// 透過peer介面從AggregateStore取得聚合狀態列表
	aggregateDataList, err := r.peer.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 使用AggregateMapper批量轉換：AggregateData → Domain Aggregate
	wallets := make([]*model.Wallet, len(aggregateDataList))
	for i, aggregateData := range aggregateDataList {
		wallet, err := r.mapper.ToDomain(aggregateData)
		if err != nil {
			return nil, err
		}
		wallets[i] = wallet
	}

	return wallets, nil
}

// 注意：移除了直接實現WalletRepositoryPeer介面的方法
// Repository Impl (Layer 2) 只應該通過peer介面與Layer 3溝通
// 避免破壞分層架構的依賴規則
