package test

import (
	"fmt"
	"sync"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

// FakeIncomeCategoryRepository 假的收入分類倉庫，用於測試
type FakeIncomeCategoryRepository struct {
	categories map[string]*model.IncomeCategory
	mutex      sync.RWMutex
}

// NewFakeIncomeCategoryRepository 建立新的假倉庫
func NewFakeIncomeCategoryRepository() repository.IncomeCategoryRepository {
	return &FakeIncomeCategoryRepository{
		categories: make(map[string]*model.IncomeCategory),
	}
}

// Save 儲存收入分類聚合
func (r *FakeIncomeCategoryRepository) Save(category *model.IncomeCategory) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if category == nil {
		return fmt.Errorf("category cannot be nil")
	}
	
	// 複製分類以避免外部修改
	categoryData := *category
	r.categories[category.ID] = &categoryData
	return nil
}

// FindByID 根據ID查找收入分類聚合
func (r *FakeIncomeCategoryRepository) FindByID(id string) (*model.IncomeCategory, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}
	
	category, exists := r.categories[id]
	if !exists {
		return nil, nil // Not found
	}
	
	// 返回複製以避免外部修改
	categoryData := *category
	return &categoryData, nil
}

// Delete 根據ID刪除收入分類聚合
func (r *FakeIncomeCategoryRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}
	
	delete(r.categories, id)
	return nil
}

// FindBySubcategoryID 根據子分類ID查找包含它的收入分類聚合
func (r *FakeIncomeCategoryRepository) FindBySubcategoryID(subcategoryID string) (*model.IncomeCategory, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if subcategoryID == "" {
		return nil, fmt.Errorf("subcategory ID cannot be empty")
	}
	
	// 遍歷所有分類尋找包含該子分類的分類
	for _, category := range r.categories {
		for _, subcategory := range category.Subcategories {
			if subcategory.ID == subcategoryID {
				// 返回複製以避免外部修改
				categoryData := *category
				return &categoryData, nil
			}
		}
	}
	
	return nil, nil // Not found
}

// FindByUserID 根據用戶ID查找用戶的所有收入分類聚合
func (r *FakeIncomeCategoryRepository) FindByUserID(userID string) ([]*model.IncomeCategory, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}
	
	var result []*model.IncomeCategory
	
	for _, category := range r.categories {
		if category.UserID == userID {
			// 返回複製以避免外部修改
			categoryData := *category
			result = append(result, &categoryData)
		}
	}
	
	return result, nil
}