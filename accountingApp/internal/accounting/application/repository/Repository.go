package repository

import (
	"github.com/JingHsiu/accountingApp/internal/accounting/application/mapper"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

// Repository 通用儲存庫介面
// 遵循簡化原則，只包含基本的CRUD操作
type Repository[T any] interface {
	// Save 儲存或更新聚合
	Save(entity T) error

	// FindByID 根據ID查找聚合
	FindByID(id string) (T, error)

	// Delete 根據ID刪除聚合
	Delete(id string) error
}

// WalletRepositoryPeer Layer 3 (Adapter) 橋接介面 (Bridge Pattern)
// 使用AggregateStore抽象，不直接依賴具體的數據庫實現
type WalletRepositoryPeer interface {
	// Save 儲存錢包聚合狀態
	Save(data mapper.WalletData) error

	// FindByID 根據ID查找錢包聚合狀態（僅載入基本資料）
	FindByID(id string) (*mapper.WalletData, error)

	// FindByIDWithChildEntities 根據ID查找錢包聚合狀態並完整載入所有子實體
	FindByIDWithChildEntities(id string) (*mapper.WalletData, error)

	// FindByUserID 根據UserID查找用戶的所有錢包聚合狀態（僅載入基本資料）
	FindByUserID(userID string) ([]mapper.WalletData, error)

	// Delete 根據ID刪除錢包聚合狀態
	Delete(id string) error

	// Note: Use FindByID() for existence checks - returns nil if not found
}

// WalletRepository 錢包專用儲存庫介面 (第二層)
// 基於Aggregate Root的Repository設計
type WalletRepository interface {
	// 基本CRUD操作
	Save(wallet *model.Wallet) error
	FindByID(id string) (*model.Wallet, error)
	Delete(id string) error

	// 必要的Domain查詢
	FindByIDWithTransactions(id string) (*model.Wallet, error) // 載入完整聚合
	FindByUserID(userID string) ([]*model.Wallet, error)       // 用戶的所有錢包
}

// ExpenseCategoryRepositoryPeer 支出分類第二層儲存實現的橋接介面
type ExpenseCategoryRepositoryPeer interface {
	// SaveData 儲存支出分類資料結構
	SaveData(data mapper.ExpenseCategoryData) error

	// FindDataByID 根據ID查找支出分類資料結構
	FindDataByID(id string) (*mapper.ExpenseCategoryData, error)

	// FindDataBySubcategoryID 根據子分類ID查找支出分類資料結構
	FindDataBySubcategoryID(subcategoryID string) (*mapper.ExpenseCategoryData, error)

	// DeleteData 根據ID刪除支出分類資料
	DeleteData(id string) error
}

// ExpenseCategoryRepository 支出分類專用儲存庫介面
type ExpenseCategoryRepository interface {
	// 基本CRUD操作
	Save(category *model.ExpenseCategory) error
	FindByID(id string) (*model.ExpenseCategory, error)
	Delete(id string) error

	// 必要的Domain查詢
	FindBySubcategoryID(subcategoryID string) (*model.ExpenseCategory, error) // 透過子分類找父分類
	FindByUserID(userID string) ([]*model.ExpenseCategory, error)             // 用戶的所有分類
}

// IncomeCategoryRepositoryPeer 收入分類第二層儲存實現的橋接介面
type IncomeCategoryRepositoryPeer interface {
	// SaveData 儲存收入分類資料結構
	SaveData(data mapper.IncomeCategoryData) error

	// FindDataByID 根據ID查找收入分類資料結構
	FindDataByID(id string) (*mapper.IncomeCategoryData, error)

	// FindDataBySubcategoryID 根據子分類ID查找收入分類資料結構
	FindDataBySubcategoryID(subcategoryID string) (*mapper.IncomeCategoryData, error)

	// DeleteData 根據ID刪除收入分類資料
	DeleteData(id string) error
}

// IncomeCategoryRepository 收入分類專用儲存庫介面
type IncomeCategoryRepository interface {
	// 基本CRUD操作
	Save(category *model.IncomeCategory) error
	FindByID(id string) (*model.IncomeCategory, error)
	Delete(id string) error

	// 必要的Domain查詢
	FindBySubcategoryID(subcategoryID string) (*model.IncomeCategory, error) // 透過子分類找父分類
	FindByUserID(userID string) ([]*model.IncomeCategory, error)             // 用戶的所有分類
}
