package repository

import (
	"github.com/JingHsiu/accountingApp/internal/accounting/adapter/store"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/mapper"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
)

// PgWalletRepositoryPeerAdapter Layer 3 (Adapter) 實現
// 使用QueryAggregateStore抽象，遵循正確分層：Peer (Layer 3) → AggregateStore (Layer 4)
type PgWalletRepositoryPeerAdapter struct {
	walletStore store.QueryAggregateStore[mapper.WalletData] // QueryAggregateStore抽象
}

// NewPgWalletRepositoryPeerAdapter 創建PostgreSQL錢包儲存實現
// 接受QueryAggregateStore，遵循依賴反轉原則
func NewPgWalletRepositoryPeerAdapter(walletStore store.QueryAggregateStore[mapper.WalletData]) repository.WalletRepositoryPeer {
	return &PgWalletRepositoryPeerAdapter{
		walletStore: walletStore,
	}
}

// Save 儲存錢包聚合狀態 (實現WalletRepositoryPeer介面)
func (p *PgWalletRepositoryPeerAdapter) Save(data mapper.WalletData) error {
	return p.walletStore.Save(data)
}

// FindByID 根據ID查找錢包聚合狀態 (實現WalletRepositoryPeer介面)
func (p *PgWalletRepositoryPeerAdapter) FindByID(id string) (*mapper.WalletData, error) {
	return p.walletStore.FindByID(id)
}

// Delete 根據ID刪除錢包聚合狀態 (實現WalletRepositoryPeer介面)
func (p *PgWalletRepositoryPeerAdapter) Delete(id string) error {
	return p.walletStore.Delete(id)
}

// FindByUserID 根據UserID查找用戶的所有錢包聚合狀態 (實現WalletRepositoryPeer介面)
func (p *PgWalletRepositoryPeerAdapter) FindByUserID(userID string) ([]mapper.WalletData, error) {
	// 使用QueryAggregateStore的FindBy方法查詢用戶的所有錢包
	criteria := map[string]interface{}{
		"user_id": userID,
	}
	
	return p.walletStore.FindBy(criteria)
}
