package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// IncomeSubcategory 是 IncomeCategory 聚合內的實體 (Entity)
// 不能獨立存在，只能通過 IncomeCategory 聚合根操作
type IncomeSubcategory struct {
	ID   string
	Name CategoryName
	// 移除 ParentID - Subcategory 屬於 Category 聚合內部，不需要向上引用
}

// IncomeSubcategory 不提供公開建構函式
// 只能通過 IncomeCategory.AddSubcategory() 建立
func newIncomeSubcategory(name CategoryName) *IncomeSubcategory {
	return &IncomeSubcategory{
		ID:   uuid.NewString(),
		Name: name,
	}
}

type IncomeCategory struct {
	ID            string
	UserID        string
	Name          CategoryName
	Subcategories []IncomeSubcategory
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewIncomeCategory(userID string, name CategoryName) (*IncomeCategory, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	
	now := time.Now()
	return &IncomeCategory{
		ID:            uuid.NewString(),
		UserID:        userID,
		Name:          name,
		Subcategories: make([]IncomeSubcategory, 0),
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// AddSubcategory 透過聚合根新增子分類
func (ic *IncomeCategory) AddSubcategory(name CategoryName) (*IncomeSubcategory, error) {
	// 業務規則：檢查名稱不能重複
	for _, existing := range ic.Subcategories {
		if existing.Name.Equals(name) {
			return nil, errors.New("subcategory with this name already exists")
		}
	}

	// 透過聚合根建立新的子分類
	subcategory := newIncomeSubcategory(name)
	ic.Subcategories = append(ic.Subcategories, *subcategory)
	ic.UpdatedAt = time.Now()
	
	return subcategory, nil
}

func (ic *IncomeCategory) RemoveSubcategory(subcategoryID string) error {
	for i, sub := range ic.Subcategories {
		if sub.ID == subcategoryID {
			ic.Subcategories = append(ic.Subcategories[:i], ic.Subcategories[i+1:]...)
			ic.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("subcategory not found")
}

// GetSubcategory 取得子分類
func (ic *IncomeCategory) GetSubcategory(subcategoryID string) (*IncomeSubcategory, error) {
	for _, sub := range ic.Subcategories {
		if sub.ID == subcategoryID {
			return &sub, nil
		}
	}
	return nil, errors.New("subcategory not found")
}

// GetAllSubcategories 取得所有子分類
func (ic *IncomeCategory) GetAllSubcategories() []IncomeSubcategory {
	return ic.Subcategories
}

// HasSubcategories 檢查是否有子分類
func (ic *IncomeCategory) HasSubcategories() bool {
	return len(ic.Subcategories) > 0
}

// ValidateSubcategoryExists 驗證子分類是否存在於此聚合中
func (ic *IncomeCategory) ValidateSubcategoryExists(subcategoryID string) error {
	_, err := ic.GetSubcategory(subcategoryID)
	return err
}

// UpdateSubcategoryName 更新子分類名稱
func (ic *IncomeCategory) UpdateSubcategoryName(subcategoryID string, newName CategoryName) error {
	// 檢查新名稱不與其他子分類重複
	for _, existing := range ic.Subcategories {
		if existing.ID != subcategoryID && existing.Name.Equals(newName) {
			return errors.New("subcategory with this name already exists")
		}
	}
	
	// 找到並更新
	for i, sub := range ic.Subcategories {
		if sub.ID == subcategoryID {
			ic.Subcategories[i].Name = newName
			ic.UpdatedAt = time.Now()
			return nil
		}
	}
	
	return errors.New("subcategory not found")
}

