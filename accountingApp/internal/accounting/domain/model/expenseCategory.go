package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ExpenseSubcategory 是 Category 聚合內的實體 (Entity)
// 不能獨立存在，只能通過 ExpenseCategory 聚合根操作
type ExpenseSubcategory struct {
	ID   string
	Name CategoryName
	// 移除 ParentID - Subcategory 屬於 Category 聚合內部，不需要向上引用
}

// ExpenseSubcategory 不提供公開建構函式
// 只能通過 ExpenseCategory.AddSubcategory() 建立
func newExpenseSubcategory(name CategoryName) *ExpenseSubcategory {
	return &ExpenseSubcategory{
		ID:   uuid.NewString(),
		Name: name,
	}
}

type ExpenseCategory struct {
	ID            string
	UserID        string
	Name          CategoryName
	Subcategories []ExpenseSubcategory
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewExpenseCategory(userID string, name CategoryName) (*ExpenseCategory, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	now := time.Now()
	return &ExpenseCategory{
		ID:            uuid.NewString(),
		UserID:        userID,
		Name:          name,
		Subcategories: make([]ExpenseSubcategory, 0),
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// AddSubcategory 透過聚合根新增子分類
func (ec *ExpenseCategory) AddSubcategory(name CategoryName) (*ExpenseSubcategory, error) {
	// 業務規則：檢查名稱不能重複
	for _, existing := range ec.Subcategories {
		if existing.Name.Equals(name) {
			return nil, errors.New("subcategory with this name already exists")
		}
	}

	// 透過聚合根建立新的子分類
	subcategory := newExpenseSubcategory(name)
	ec.Subcategories = append(ec.Subcategories, *subcategory)
	ec.UpdatedAt = time.Now()

	return subcategory, nil
}

func (ec *ExpenseCategory) RemoveSubcategory(subcategoryID string) error {
	for i, sub := range ec.Subcategories {
		if sub.ID == subcategoryID {
			ec.Subcategories = append(ec.Subcategories[:i], ec.Subcategories[i+1:]...)
			ec.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("subcategory not found")
}

func (ec *ExpenseCategory) GetSubcategory(subcategoryID string) (*ExpenseSubcategory, error) {
	for _, sub := range ec.Subcategories {
		if sub.ID == subcategoryID {
			return &sub, nil
		}
	}
	return nil, errors.New("subcategory not found")
}

// GetAllSubcategories 取得所有子分類
func (ec *ExpenseCategory) GetAllSubcategories() []ExpenseSubcategory {
	return ec.Subcategories
}

// HasSubcategories 檢查是否有子分類
func (ec *ExpenseCategory) HasSubcategories() bool {
	return len(ec.Subcategories) > 0
}

// ValidateSubcategoryExists 驗證子分類是否存在於此聚合中
func (ec *ExpenseCategory) ValidateSubcategoryExists(subcategoryID string) error {
	_, err := ec.GetSubcategory(subcategoryID)
	return err
}

// UpdateSubcategoryName 更新子分類名稱
func (ec *ExpenseCategory) UpdateSubcategoryName(subcategoryID string, newName CategoryName) error {
	// 檢查新名稱不與其他子分類重複
	for _, existing := range ec.Subcategories {
		if existing.ID != subcategoryID && existing.Name.Equals(newName) {
			return errors.New("subcategory with this name already exists")
		}
	}
	
	// 找到並更新
	for i, sub := range ec.Subcategories {
		if sub.ID == subcategoryID {
			ec.Subcategories[i].Name = newName
			ec.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("subcategory not found")
}
