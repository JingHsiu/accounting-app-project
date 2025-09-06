package repository

import (
	"fmt"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/mapper"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

// IncomeCategoryRepositoryImpl 收入分類倉庫實作
type IncomeCategoryRepositoryImpl struct {
	peer   IncomeCategoryRepositoryPeer
	mapper *mapper.IncomeCategoryMapper
}

// NewIncomeCategoryRepositoryImpl 建立新的收入分類倉庫實作
func NewIncomeCategoryRepositoryImpl(peer IncomeCategoryRepositoryPeer) IncomeCategoryRepository {
	return &IncomeCategoryRepositoryImpl{
		peer:   peer,
		mapper: mapper.NewIncomeCategoryMapper(),
	}
}

// Save 儲存收入分類聚合
func (r *IncomeCategoryRepositoryImpl) Save(category *model.IncomeCategory) error {
	if category == nil {
		return fmt.Errorf("category cannot be nil")
	}

	// 轉換Domain Model為Data Model
	data := r.mapper.ToData(category)
	
	// 透過Peer儲存資料
	return r.peer.SaveData(data)
}

// FindByID 根據ID查找收入分類聚合
func (r *IncomeCategoryRepositoryImpl) FindByID(id string) (*model.IncomeCategory, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	// 透過Peer查找資料
	data, err := r.peer.FindDataByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find income category by ID: %w", err)
	}
	if data == nil {
		return nil, nil // Not found
	}

	// 轉換Data Model為Domain Model
	return r.mapper.ToDomain(*data)
}

// Delete 根據ID刪除收入分類聚合
func (r *IncomeCategoryRepositoryImpl) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	return r.peer.DeleteData(id)
}

// FindBySubcategoryID 根據子分類ID查找包含它的收入分類聚合
func (r *IncomeCategoryRepositoryImpl) FindBySubcategoryID(subcategoryID string) (*model.IncomeCategory, error) {
	if subcategoryID == "" {
		return nil, fmt.Errorf("subcategory ID cannot be empty")
	}

	// 透過Peer查找包含該子分類的分類資料
	data, err := r.peer.FindDataBySubcategoryID(subcategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to find income category by subcategory ID: %w", err)
	}
	if data == nil {
		return nil, nil // Not found
	}

	// 轉換Data Model為Domain Model
	return r.mapper.ToDomain(*data)
}

// FindByUserID 根據用戶ID查找用戶的所有收入分類聚合
func (r *IncomeCategoryRepositoryImpl) FindByUserID(userID string) ([]*model.IncomeCategory, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	// 臨時實現：由於Peer層尚未支持FindDataByUserID，先返回預設分類
	// 根據default_categories.sql，預設收入分類的ID範圍是 default-income-1 到 default-income-4
	defaultCategoryIDs := []string{
		"default-income-1", // 薪資
		"default-income-2", // 投資  
		"default-income-3", // 副業
		"default-income-4", // 其他收入
	}

	var categories []*model.IncomeCategory
	for _, categoryID := range defaultCategoryIDs {
		data, err := r.peer.FindDataByID(categoryID)
		if err != nil {
			// 如果單個分類找不到，記錄但繼續處理其他分類
			continue
		}
		if data != nil {
			category, err := r.mapper.ToDomain(*data)
			if err != nil {
				// 轉換失敗，記錄但繼續處理其他分類
				continue
			}
			categories = append(categories, category)
		}
	}

	return categories, nil
}