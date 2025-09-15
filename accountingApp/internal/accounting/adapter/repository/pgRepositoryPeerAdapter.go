package repository

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/adapter/store"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/mapper"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/frameworks/database"
)

// PgWalletRepositoryPeerAdapter Layer 3 (Adapter) 實現
// 使用QueryAggregateStore抽象，遵循正確分層：Peer (Layer 3) → AggregateStore (Layer 4)
type PgWalletRepositoryPeerAdapter struct {
	walletStore      store.QueryAggregateStore[mapper.WalletData] // QueryAggregateStore抽象
	dbClient        database.DatabaseClient                       // 直接數據庫訪問用於複雜事務
	incomeStore     store.BatchAggregateStore[mapper.IncomeRecordData]
	expenseStore    store.BatchAggregateStore[mapper.ExpenseRecordData]
	transferStore   store.BatchAggregateStore[mapper.TransferData]
}

// NewPgWalletRepositoryPeerAdapter 創建PostgreSQL錢包儲存實現
// 接受QueryAggregateStore，遵循依賴反轉原則
func NewPgWalletRepositoryPeerAdapter(
	walletStore store.QueryAggregateStore[mapper.WalletData],
	dbClient database.DatabaseClient,
	incomeStore store.BatchAggregateStore[mapper.IncomeRecordData],
	expenseStore store.BatchAggregateStore[mapper.ExpenseRecordData],
	transferStore store.BatchAggregateStore[mapper.TransferData],
) repository.WalletRepositoryPeer {
	return &PgWalletRepositoryPeerAdapter{
		walletStore:   walletStore,
		dbClient:     dbClient,
		incomeStore:  incomeStore,
		expenseStore: expenseStore,
		transferStore: transferStore,
	}
}

// Save 儲存錢包聚合狀態 (實現WalletRepositoryPeer介面)
func (p *PgWalletRepositoryPeerAdapter) Save(data mapper.WalletData) error {
	// 開始交易以確保聚合完整性
	return p.saveWithTransaction(data)
}

// saveWithTransaction 在交易中保存完整聚合
// 實現完整的DDD聚合持久化模式，確保事務完整性
func (p *PgWalletRepositoryPeerAdapter) saveWithTransaction(data mapper.WalletData) error {
	// 開始數據庫事務
	tx, err := p.dbClient.BeginTx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. 保存錢包主體實體
	err = p.saveWalletInTransaction(tx, data)
	if err != nil {
		return fmt.Errorf("failed to save wallet: %w", err)
	}

	// 2. 保存子實體 - 收入記錄
	if len(data.IncomeRecords) > 0 {
		err = p.saveIncomeRecords(tx, data.IncomeRecords)
		if err != nil {
			return fmt.Errorf("failed to save income records: %w", err)
		}
	}

	// 3. 保存子實體 - 支出記錄
	if len(data.ExpenseRecords) > 0 {
		err = p.saveExpenseRecords(tx, data.ExpenseRecords)
		if err != nil {
			return fmt.Errorf("failed to save expense records: %w", err)
		}
	}

	// 4. 保存子實體 - 轉帳記錄
	if len(data.Transfers) > 0 {
		err = p.saveTransfers(tx, data.Transfers)
		if err != nil {
			return fmt.Errorf("failed to save transfers: %w", err)
		}
	}

	// 提交事務
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// FindByID 根據ID查找錢包聚合狀態 (實現WalletRepositoryPeer介面)
// 只加載錢包基本資料，不加載子實體
func (p *PgWalletRepositoryPeerAdapter) FindByID(id string) (*mapper.WalletData, error) {
	walletData, err := p.walletStore.FindByID(id)
	if err != nil || walletData == nil {
		return walletData, err
	}

	// 設置為未完全加載狀態
	walletData.IsFullyLoaded = false
	return walletData, nil
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
	
	wallets, err := p.walletStore.FindBy(criteria)
	if err != nil {
		return nil, err
	}

	// 設置所有錢包為未完全加載狀態
	for i := range wallets {
		wallets[i].IsFullyLoaded = false
	}

	return wallets, nil
}

// FindByIDWithChildEntities 根據ID查找錢包聚合狀態並完整載入所有子實體
// 實現完整的聚合重建模式
func (p *PgWalletRepositoryPeerAdapter) FindByIDWithChildEntities(id string) (*mapper.WalletData, error) {
	// 1. 首先獲取錢包基本資料
	walletData, err := p.walletStore.FindByID(id)
	if err != nil || walletData == nil {
		return walletData, err
	}

	// 2. 載入所有子實體
	err = p.loadChildEntities(walletData)
	if err != nil {
		return nil, fmt.Errorf("failed to load child entities for wallet %s: %w", id, err)
	}

	// 3. 標記為完全載入
	walletData.IsFullyLoaded = true
	return walletData, nil
}

// loadChildEntities 載入錢包的所有子實體
func (p *PgWalletRepositoryPeerAdapter) loadChildEntities(walletData *mapper.WalletData) error {
	// 載入收入記錄
	incomeRecords, err := p.loadIncomeRecords(walletData.ID)
	if err != nil {
		return fmt.Errorf("failed to load income records: %w", err)
	}
	walletData.IncomeRecords = incomeRecords

	// 載入支出記錄
	expenseRecords, err := p.loadExpenseRecords(walletData.ID)
	if err != nil {
		return fmt.Errorf("failed to load expense records: %w", err)
	}
	walletData.ExpenseRecords = expenseRecords

	// 載入轉帳記錄
	transfers, err := p.loadTransfers(walletData.ID)
	if err != nil {
		return fmt.Errorf("failed to load transfers: %w", err)
	}
	walletData.Transfers = transfers

	return nil
}

// saveWalletInTransaction 在事務中保存錢包主體實體
func (p *PgWalletRepositoryPeerAdapter) saveWalletInTransaction(tx database.Transaction, data mapper.WalletData) error {
	query := `
		INSERT INTO wallets (
			id, user_id, name, type, currency, 
			balance_amount, balance_currency, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			type = EXCLUDED.type,
			currency = EXCLUDED.currency,
			balance_amount = EXCLUDED.balance_amount,
			balance_currency = EXCLUDED.balance_currency,
			updated_at = EXCLUDED.updated_at
	`
	
	_, err := tx.Exec(query,
		data.ID, data.UserID, data.Name, data.Type, data.Currency,
		data.BalanceAmount, data.BalanceCurrency, data.CreatedAt, data.UpdatedAt)
	
	return err
}

// saveIncomeRecords 在事務中批次保存收入記錄
func (p *PgWalletRepositoryPeerAdapter) saveIncomeRecords(tx database.Transaction, records []mapper.IncomeRecordData) error {
	if len(records) == 0 {
		return nil
	}

	// Use individual upserts to preserve existing records and only insert new ones
	// This prevents overwriting existing income records when adding new ones
	query := `
		INSERT INTO income_records (
			id, wallet_id, category_id, amount, currency, description, date, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO NOTHING
	`

	for _, record := range records {
		_, err := tx.Exec(query,
			record.ID, record.WalletID, record.SubcategoryID, record.Amount,
			record.Currency, record.Description, record.Date, record.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to save income record %s: %w", record.ID, err)
		}
	}

	return nil
}

// saveExpenseRecords 在事務中批次保存支出記錄
func (p *PgWalletRepositoryPeerAdapter) saveExpenseRecords(tx database.Transaction, records []mapper.ExpenseRecordData) error {
	if len(records) == 0 {
		return nil
	}

	// 先清除該錢包的現有支出記錄
	walletID := records[0].WalletID
	deleteQuery := "DELETE FROM expense_records WHERE wallet_id = $1"
	_, err := tx.Exec(deleteQuery, walletID)
	if err != nil {
		return fmt.Errorf("failed to delete existing expense records: %w", err)
	}

	// 批次插入新記錄
	query := `
		INSERT INTO expense_records (
			id, wallet_id, category_id, amount, currency, description, date, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	for _, record := range records {
		_, err = tx.Exec(query,
			record.ID, record.WalletID, record.SubcategoryID, record.Amount,
			record.Currency, record.Description, record.Date, record.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to save expense record %s: %w", record.ID, err)
		}
	}

	return nil
}

// saveTransfers 在事務中批次保存轉帳記錄
func (p *PgWalletRepositoryPeerAdapter) saveTransfers(tx database.Transaction, transfers []mapper.TransferData) error {
	if len(transfers) == 0 {
		return nil
	}

	// 清除相關的轉帳記錄（FROM 或 TO 此錢包的轉帳）
	walletID := transfers[0].FromWalletID // 假設所有轉帳都是從同一個錢包
	deleteQuery := "DELETE FROM transfers WHERE from_wallet_id = $1 OR to_wallet_id = $1"
	_, err := tx.Exec(deleteQuery, walletID)
	if err != nil {
		return fmt.Errorf("failed to delete existing transfers: %w", err)
	}

	// 批次插入新記錄
	query := `
		INSERT INTO transfers (
			id, from_wallet_id, to_wallet_id, amount, currency, 
			fee_amount, fee_currency, description, date, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	for _, transfer := range transfers {
		_, err = tx.Exec(query,
			transfer.ID, transfer.FromWalletID, transfer.ToWalletID,
			transfer.Amount, transfer.Currency, transfer.Fee, transfer.Currency,
			transfer.Description, transfer.Date, transfer.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to save transfer %s: %w", transfer.ID, err)
		}
	}

	return nil
}

// loadIncomeRecords 載入特定錢包的所有收入記錄
func (p *PgWalletRepositoryPeerAdapter) loadIncomeRecords(walletID string) ([]mapper.IncomeRecordData, error) {
	query := `
		SELECT id, wallet_id, category_id, amount, currency, description, date, created_at
		FROM income_records
		WHERE wallet_id = $1
		ORDER BY date DESC, created_at DESC
	`

	rows, err := p.dbClient.Query(query, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to query income records: %w", err)
	}
	defer rows.Close()

	var records []mapper.IncomeRecordData
	for rows.Next() {
		var record mapper.IncomeRecordData
		err = rows.Scan(
			&record.ID, &record.WalletID, &record.SubcategoryID,
			&record.Amount, &record.Currency, &record.Description,
			&record.Date, &record.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan income record: %w", err)
		}
		records = append(records, record)
	}

	return records, nil
}

// loadExpenseRecords 載入特定錢包的所有支出記錄
func (p *PgWalletRepositoryPeerAdapter) loadExpenseRecords(walletID string) ([]mapper.ExpenseRecordData, error) {
	query := `
		SELECT id, wallet_id, category_id, amount, currency, description, date, created_at
		FROM expense_records
		WHERE wallet_id = $1
		ORDER BY date DESC, created_at DESC
	`

	rows, err := p.dbClient.Query(query, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to query expense records: %w", err)
	}
	defer rows.Close()

	var records []mapper.ExpenseRecordData
	for rows.Next() {
		var record mapper.ExpenseRecordData
		err = rows.Scan(
			&record.ID, &record.WalletID, &record.SubcategoryID,
			&record.Amount, &record.Currency, &record.Description,
			&record.Date, &record.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expense record: %w", err)
		}
		records = append(records, record)
	}

	return records, nil
}

// loadTransfers 載入特定錢包相關的所有轉帳記錄
func (p *PgWalletRepositoryPeerAdapter) loadTransfers(walletID string) ([]mapper.TransferData, error) {
	query := `
		SELECT id, from_wallet_id, to_wallet_id, amount, currency, 
			   fee_amount as fee, description, date, created_at
		FROM transfers
		WHERE from_wallet_id = $1 OR to_wallet_id = $1
		ORDER BY date DESC, created_at DESC
	`

	rows, err := p.dbClient.Query(query, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transfers: %w", err)
	}
	defer rows.Close()

	var transfers []mapper.TransferData
	for rows.Next() {
		var transfer mapper.TransferData
		err = rows.Scan(
			&transfer.ID, &transfer.FromWalletID, &transfer.ToWalletID,
			&transfer.Amount, &transfer.Currency, &transfer.Fee,
			&transfer.Description, &transfer.Date, &transfer.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transfer: %w", err)
		}
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}
