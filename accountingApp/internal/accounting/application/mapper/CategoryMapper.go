package mapper

import (
	"time"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

// ExpenseCategoryData 支出分類的持久化資料結構
type ExpenseCategoryData struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (ecd ExpenseCategoryData) GetID() string {
	return ecd.ID
}

// ExpenseSubcategoryData 支出子分類的持久化資料結構
type ExpenseSubcategoryData struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	ParentID string `db:"parent_id"`
}

// IncomeCategoryData 收入分類的持久化資料結構
type IncomeCategoryData struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (icd IncomeCategoryData) GetID() string {
	return icd.ID
}

// IncomeSubcategoryData 收入子分類的持久化資料結構
type IncomeSubcategoryData struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	ParentID string `db:"parent_id"`
}

// ExpenseCategoryMapper 支出分類聚合的資料轉換器
type ExpenseCategoryMapper struct{}

func NewExpenseCategoryMapper() *ExpenseCategoryMapper {
	return &ExpenseCategoryMapper{}
}

// ToData 將ExpenseCategory Domain Model轉換為ExpenseCategoryData
func (m *ExpenseCategoryMapper) ToData(category *model.ExpenseCategory) ExpenseCategoryData {
	return ExpenseCategoryData{
		ID:        category.ID,
		UserID:    category.UserID,
		Name:      category.Name.Value,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

// ToDomain 將ExpenseCategoryData轉換為ExpenseCategory Domain Model
func (m *ExpenseCategoryMapper) ToDomain(data ExpenseCategoryData) (*model.ExpenseCategory, error) {
	categoryName, err := model.NewCategoryName(data.Name)
	if err != nil {
		return nil, err
	}
	
	return &model.ExpenseCategory{
		ID:            data.ID,
		UserID:        data.UserID,
		Name:          *categoryName,
		Subcategories: make([]model.ExpenseSubcategory, 0), // 子分類需要另外查詢
		CreatedAt:     data.CreatedAt,
		UpdatedAt:     data.UpdatedAt,
	}, nil
}

// ToSubcategoryData 將ExpenseSubcategory轉換為ExpenseSubcategoryData
// parentID 需要從聚合根傳入，因為 Subcategory 不再持有 ParentID
func (m *ExpenseCategoryMapper) ToSubcategoryData(subcategory model.ExpenseSubcategory, parentID string) ExpenseSubcategoryData {
	return ExpenseSubcategoryData{
		ID:       subcategory.ID,
		Name:     subcategory.Name.Value,
		ParentID: parentID,
	}
}

// ToSubcategoryDomain 將ExpenseSubcategoryData轉換為ExpenseSubcategory
// ParentID 不再需要，因為 Subcategory 屬於 Category 聚合內部
func (m *ExpenseCategoryMapper) ToSubcategoryDomain(data ExpenseSubcategoryData) (*model.ExpenseSubcategory, error) {
	categoryName, err := model.NewCategoryName(data.Name)
	if err != nil {
		return nil, err
	}
	
	return &model.ExpenseSubcategory{
		ID:   data.ID,
		Name: *categoryName,
	}, nil
}

// IncomeCategoryMapper 收入分類聚合的資料轉換器
type IncomeCategoryMapper struct{}

func NewIncomeCategoryMapper() *IncomeCategoryMapper {
	return &IncomeCategoryMapper{}
}

// ToData 將IncomeCategory Domain Model轉換為IncomeCategoryData
func (m *IncomeCategoryMapper) ToData(category *model.IncomeCategory) IncomeCategoryData {
	return IncomeCategoryData{
		ID:        category.ID,
		UserID:    category.UserID,
		Name:      category.Name.Value,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

// ToDomain 將IncomeCategoryData轉換為IncomeCategory Domain Model
func (m *IncomeCategoryMapper) ToDomain(data IncomeCategoryData) (*model.IncomeCategory, error) {
	categoryName, err := model.NewCategoryName(data.Name)
	if err != nil {
		return nil, err
	}
	
	return &model.IncomeCategory{
		ID:            data.ID,
		UserID:        data.UserID,
		Name:          *categoryName,
		Subcategories: make([]model.IncomeSubcategory, 0), // 子分類需要另外查詢
		CreatedAt:     data.CreatedAt,
		UpdatedAt:     data.UpdatedAt,
	}, nil
}

// ToSubcategoryData 將IncomeSubcategory轉換為IncomeSubcategoryData
// parentID 需要從聚合根傳入，因為 Subcategory 不再持有 ParentID
func (m *IncomeCategoryMapper) ToSubcategoryData(subcategory model.IncomeSubcategory, parentID string) IncomeSubcategoryData {
	return IncomeSubcategoryData{
		ID:       subcategory.ID,
		Name:     subcategory.Name.Value,
		ParentID: parentID,
	}
}

// ToSubcategoryDomain 將IncomeSubcategoryData轉換為IncomeSubcategory
// ParentID 不再需要，因為 Subcategory 屬於 Category 聚合內部
func (m *IncomeCategoryMapper) ToSubcategoryDomain(data IncomeSubcategoryData) (*model.IncomeSubcategory, error) {
	categoryName, err := model.NewCategoryName(data.Name)
	if err != nil {
		return nil, err
	}
	
	return &model.IncomeSubcategory{
		ID:   data.ID,
		Name: *categoryName,
	}, nil
}

// 確保Mapper實現介面
var _ Mapper[*model.ExpenseCategory, ExpenseCategoryData] = (*ExpenseCategoryMapper)(nil)
var _ Mapper[*model.IncomeCategory, IncomeCategoryData] = (*IncomeCategoryMapper)(nil)