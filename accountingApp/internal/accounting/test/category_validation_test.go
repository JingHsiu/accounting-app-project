package test

import (
	"testing"
	"time"

	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/mapper"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

// MockExpenseCategoryRepository 模擬支出分類儲存庫
type MockExpenseCategoryRepository struct {
	categories map[string]*model.ExpenseCategory
	subcategoryToCategory map[string]string
}

func NewMockExpenseCategoryRepository() *MockExpenseCategoryRepository {
	return &MockExpenseCategoryRepository{
		categories: make(map[string]*model.ExpenseCategory),
		subcategoryToCategory: make(map[string]string),
	}
}

func (m *MockExpenseCategoryRepository) Save(category *model.ExpenseCategory) error {
	m.categories[category.ID] = category
	// 更新子分類到分類的映射
	for _, sub := range category.GetAllSubcategories() {
		m.subcategoryToCategory[sub.ID] = category.ID
	}
	return nil
}

func (m *MockExpenseCategoryRepository) FindByID(id string) (*model.ExpenseCategory, error) {
	if category, exists := m.categories[id]; exists {
		return category, nil
	}
	return nil, nil
}

func (m *MockExpenseCategoryRepository) FindBySubcategoryID(subcategoryID string) (*model.ExpenseCategory, error) {
	if categoryID, exists := m.subcategoryToCategory[subcategoryID]; exists {
		return m.categories[categoryID], nil
	}
	return nil, nil
}

func (m *MockExpenseCategoryRepository) Delete(id string) error {
	delete(m.categories, id)
	return nil
}

// 實現WalletRepositoryPeer介面的Bridge Pattern方法
func (m *MockExpenseCategoryRepository) SaveData(data mapper.ExpenseCategoryData) error {
	// Mock implementation - convert data back to domain model for storage
	categoryName, _ := model.NewCategoryName(data.Name)
	category := &model.ExpenseCategory{
		ID:            data.ID,
		UserID:        data.UserID,
		Name:          *categoryName,
		Subcategories: make([]model.ExpenseSubcategory, 0),
		CreatedAt:     data.CreatedAt,
		UpdatedAt:     data.UpdatedAt,
	}
	return m.Save(category)
}

func (m *MockExpenseCategoryRepository) FindDataByID(id string) (*mapper.ExpenseCategoryData, error) {
	category, err := m.FindByID(id)
	if err != nil || category == nil {
		return nil, err
	}
	// Convert domain model back to data structure
	mapper := mapper.NewExpenseCategoryMapper()
	data := mapper.ToData(category)
	return &data, nil
}

func (m *MockExpenseCategoryRepository) FindDataBySubcategoryID(subcategoryID string) (*mapper.ExpenseCategoryData, error) {
	category, err := m.FindBySubcategoryID(subcategoryID)
	if err != nil || category == nil {
		return nil, err
	}
	// Convert domain model back to data structure
	mapper := mapper.NewExpenseCategoryMapper()
	data := mapper.ToData(category)
	return &data, nil
}

func (m *MockExpenseCategoryRepository) DeleteData(id string) error {
	return m.Delete(id)
}

// MockWalletRepository 模擬錢包儲存庫
type MockWalletRepository struct {
	wallets map[string]*model.Wallet
}

func NewMockWalletRepository() *MockWalletRepository {
	return &MockWalletRepository{
		wallets: make(map[string]*model.Wallet),
	}
}

func (m *MockWalletRepository) Save(wallet *model.Wallet) error {
	m.wallets[wallet.ID] = wallet
	return nil
}

func (m *MockWalletRepository) FindByID(id string) (*model.Wallet, error) {
	if wallet, exists := m.wallets[id]; exists {
		return wallet, nil
	}
	return nil, nil
}

func (m *MockWalletRepository) Delete(id string) error {
	delete(m.wallets, id)
	return nil
}

// 實現WalletRepositoryPeer介面的Bridge Pattern方法
func (m *MockWalletRepository) SaveData(data mapper.WalletData) error {
	// Mock implementation - convert data back to domain model for storage
	money, _ := model.NewMoney(data.BalanceAmount, data.BalanceCurrency)
	wallet := &model.Wallet{
		ID:        data.ID,
		UserID:    data.UserID,
		Name:      data.Name,
		Type:      model.WalletType(data.Type),
		Balance:   *money,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
	return m.Save(wallet)
}

func (m *MockWalletRepository) FindDataByID(id string) (*mapper.WalletData, error) {
	wallet, err := m.FindByID(id)
	if err != nil || wallet == nil {
		return nil, err
	}
	// Convert domain model back to data structure
	mapper := mapper.NewWalletMapper()
	data := mapper.ToData(wallet)
	return &data, nil
}

func (m *MockWalletRepository) DeleteData(id string) error {
	return m.Delete(id)
}

// TestAddExpenseWithValidation 測試新增支出時的分類驗證
func TestAddExpenseWithValidation(t *testing.T) {
	// 設置測試資料
	categoryRepo := NewMockExpenseCategoryRepository()
	walletRepo := NewMockWalletRepository()

	// 1. 建立分類聚合和子分類
	categoryName, err := model.NewCategoryName("餐飲")
	if err != nil {
		t.Fatalf("Failed to create category name: %v", err)
	}

	category, err := model.NewExpenseCategory("user123", *categoryName)
	if err != nil {
		t.Fatalf("Failed to create expense category: %v", err)
	}

	subcategoryName, err := model.NewCategoryName("早餐")
	if err != nil {
		t.Fatalf("Failed to create subcategory name: %v", err)
	}

	subcategory, err := category.AddSubcategory(*subcategoryName)
	if err != nil {
		t.Fatalf("Failed to add subcategory: %v", err)
	}

	// 儲存分類聚合
	err = categoryRepo.Save(category)
	if err != nil {
		t.Fatalf("Failed to save category: %v", err)
	}

	// 2. 建立錢包
	wallet, err := model.NewWallet("user123", "測試錢包", model.WalletTypeCash, "USD")
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	// 先新增收入以便有餘額
	incomeAmount, err := model.NewMoney(10000, "USD") // $100.00
	if err != nil {
		t.Fatalf("Failed to create income amount: %v", err)
	}

	_, err = wallet.AddIncome(*incomeAmount, "income-subcategory-123", "測試收入", wallet.CreatedAt)
	if err != nil {
		t.Fatalf("Failed to add income: %v", err)
	}

	// 儲存錢包
	err = walletRepo.Save(wallet)
	if err != nil {
		t.Fatalf("Failed to save wallet: %v", err)
	}

	// 3. 建立服務
	addExpenseService := command.NewAddExpenseService(walletRepo, categoryRepo)

	// 4. 測試有效的子分類ID
	validInput := command.AddExpenseInput{
		WalletID:      wallet.ID,
		SubcategoryID: subcategory.ID,
		Amount:        2500, // $25.00
		Currency:      "USD",
		Description:   "測試支出",
		Date:          time.Now(),
	}

	result := addExpenseService.Execute(validInput)
	if result.GetExitCode() != common.Success {
		t.Errorf("Expected success for valid subcategory, got: %s", result.GetMessage())
	}

	expenseID := result.GetID()
	if expenseID == "" {
		t.Error("Expected expense ID to be returned")
	}

	// 5. 測試無效的子分類ID
	invalidInput := command.AddExpenseInput{
		WalletID:      wallet.ID,
		SubcategoryID: "invalid-subcategory-id",
		Amount:        1000,
		Currency:      "USD",
		Description:   "無效測試",
		Date:          time.Now(),
	}

	result = addExpenseService.Execute(invalidInput)
	if result.GetExitCode() != common.Failure {
		t.Error("Expected failure for invalid subcategory")
	}

	expectedMessage := "Subcategory not found in any category"
	if result.GetMessage() != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, result.GetMessage())
	}

	// 6. 驗證錢包餘額正確更新 (只有有效交易)
	updatedWallet, err := walletRepo.FindByID(wallet.ID)
	if err != nil {
		t.Fatalf("Failed to find updated wallet: %v", err)
	}

	expectedBalance := int64(7500) // $100.00 - $25.00 = $75.00
	if updatedWallet.Balance.Amount != expectedBalance {
		t.Errorf("Expected balance %d, got %d", expectedBalance, updatedWallet.Balance.Amount)
	}
}

// TestAddExpenseValidationIntegration 測試完整的聚合根驗證整合
func TestAddExpenseValidationIntegration(t *testing.T) {
	// 設置
	categoryRepo := NewMockExpenseCategoryRepository()
	walletRepo := NewMockWalletRepository()

	// 建立多個分類聚合
	transportName, _ := model.NewCategoryName("交通")
	transportCategory, _ := model.NewExpenseCategory("user123", *transportName)
	
	busName, _ := model.NewCategoryName("公車")
	busSubcategory, _ := transportCategory.AddSubcategory(*busName)
	
	taxiName, _ := model.NewCategoryName("計程車")
	taxiSubcategory, _ := transportCategory.AddSubcategory(*taxiName)

	// 儲存交通分類
	categoryRepo.Save(transportCategory)

	// 建立第二個分類聚合
	foodName, _ := model.NewCategoryName("餐飲")
	foodCategory, _ := model.NewExpenseCategory("user123", *foodName)
	
	lunchName, _ := model.NewCategoryName("午餐")
	lunchSubcategory, _ := foodCategory.AddSubcategory(*lunchName)
	
	// 儲存餐飲分類
	categoryRepo.Save(foodCategory)

	// 建立錢包並新增收入
	wallet, _ := model.NewWallet("user123", "整合測試錢包", model.WalletTypeCash, "USD")
	incomeAmount, _ := model.NewMoney(50000, "USD") // $500.00
	wallet.AddIncome(*incomeAmount, "income-sub-123", "測試收入", wallet.CreatedAt)
	walletRepo.Save(wallet)

	// 建立服務
	service := command.NewAddExpenseService(walletRepo, categoryRepo)

	// 測試案例：不同分類的子分類都應該可以正確驗證
	testCases := []struct {
		name          string
		subcategoryID string
		amount        int64
		shouldSucceed bool
	}{
		{"公車支出", busSubcategory.ID, 500, true},
		{"計程車支出", taxiSubcategory.ID, 1500, true},
		{"午餐支出", lunchSubcategory.ID, 2000, true},
		{"無效子分類", "invalid-id", 1000, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := command.AddExpenseInput{
				WalletID:      wallet.ID,
				SubcategoryID: tc.subcategoryID,
				Amount:        tc.amount,
				Currency:      "USD",
				Description:   tc.name,
				Date:          time.Now(),
			}

			result := service.Execute(input)
			
			if tc.shouldSucceed {
				if result.GetExitCode() != common.Success {
					t.Errorf("Expected success for %s, got: %s", tc.name, result.GetMessage())
				}
			} else {
				if result.GetExitCode() != common.Failure {
					t.Errorf("Expected failure for %s, but got success", tc.name)
				}
			}
		})
	}

	// 驗證最終餘額
	finalWallet, _ := walletRepo.FindByID(wallet.ID)
	expectedBalance := int64(46000) // $500 - $5 - $15 - $20 - $10(失敗) = $450
	if finalWallet.Balance.Amount != expectedBalance {
		t.Errorf("Expected final balance %d, got %d", expectedBalance, finalWallet.Balance.Amount)
	}
}