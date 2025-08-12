package database

import (
	"database/sql"
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/mapper"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	_ "github.com/lib/pq"
)

// PostgresWalletRepositoryPeer 第三層PostgreSQL錢包儲存實現
// 只實現WalletRepositoryPeer介面，不接觸Domain Model，符合Clean Architecture依賴規則
type PostgresWalletRepositoryPeer struct {
	db *sql.DB
}

// NewPostgresWalletRepositoryPeer 創建PostgreSQL錢包儲存實現
func NewPostgresWalletRepositoryPeer(db *sql.DB) repository.WalletRepositoryPeer {
	return &PostgresWalletRepositoryPeer{
		db: db,
	}
}

// SaveData 儲存錢包資料結構 (實現WalletRepositoryPeer介面)
func (r *PostgresWalletRepositoryPeer) SaveData(data mapper.WalletData) error {
	query := `
		INSERT INTO wallets (id, user_id, name, type, currency, balance_amount, balance_currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			type = EXCLUDED.type,
			currency = EXCLUDED.currency,
			balance_amount = EXCLUDED.balance_amount,
			balance_currency = EXCLUDED.balance_currency,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.Exec(query,
		data.ID, data.UserID, data.Name, data.Type, data.Currency,
		data.BalanceAmount, data.BalanceCurrency, data.CreatedAt, data.UpdatedAt)
	print(err)
	return err
}

// FindDataByID 根據ID查找錢包資料結構 (實現WalletRepositoryPeer介面)
func (r *PostgresWalletRepositoryPeer) FindDataByID(id string) (*mapper.WalletData, error) {
	query := `
		SELECT id, user_id, name, type, currency, balance_amount, balance_currency, created_at, updated_at
		FROM wallets WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	var data mapper.WalletData
	err := row.Scan(&data.ID, &data.UserID, &data.Name, &data.Type, &data.Currency,
		&data.BalanceAmount, &data.BalanceCurrency, &data.CreatedAt, &data.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 未找到記錄
		}
		return nil, err
	}

	return &data, nil
}

//
//func (r *PostgresWalletRepository) FindByUserID(userID string) ([]*model.Wallet, error) {
//	query := `
//		SELECT id, user_id, name, type, currency, balance_amount, balance_currency, created_at, updated_at
//		FROM wallets WHERE user_id = $1 ORDER BY created_at DESC
//	`
//
//	rows, err := r.db.Query(query, userID)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var wallets []*model.Wallet
//	for rows.Next() {
//		var data mapper.WalletData
//		err := rows.Scan(&data.ID, &data.UserID, &data.Name, &data.Type, &data.Currency,
//			&data.BalanceAmount, &data.BalanceCurrency, &data.CreatedAt, &data.UpdatedAt)
//		if err != nil {
//			return nil, err
//		}
//
//		wallet, err := r.mapper.ToDomain(data)
//		if err != nil {
//			return nil, err
//		}
//
//		wallets = append(wallets, wallet)
//	}
//
//	return wallets, nil
//}

// DeleteData 根據ID刪除錢包資料 (實現WalletRepositoryPeer介面)
func (r *PostgresWalletRepositoryPeer) DeleteData(id string) error {
	query := `DELETE FROM wallets WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("wallet with id %s not found", id)
	}

	return nil
}
