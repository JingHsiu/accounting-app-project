package test

import (
	"testing"

	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

// TestExpenseCategoryAsAggregateRoot 測試 ExpenseCategory 作為聚合根的行為
func TestExpenseCategoryAsAggregateRoot(t *testing.T) {
	// 1. 創建分類聚合根
	categoryName, err := model.NewCategoryName("餐飲")
	if err != nil {
		t.Fatalf("Failed to create category name: %v", err)
	}

	category, err := model.NewExpenseCategory("user123", *categoryName)
	if err != nil {
		t.Fatalf("Failed to create expense category: %v", err)
	}

	// 2. 透過聚合根新增子分類
	subcategoryName1, err := model.NewCategoryName("早餐")
	if err != nil {
		t.Fatalf("Failed to create subcategory name: %v", err)
	}

	subcategory1, err := category.AddSubcategory(*subcategoryName1)
	if err != nil {
		t.Fatalf("Failed to add subcategory: %v", err)
	}

	if subcategory1.ID == "" {
		t.Error("Subcategory should have an ID")
	}

	// 3. 驗證子分類存在於聚合中
	err = category.ValidateSubcategoryExists(subcategory1.ID)
	if err != nil {
		t.Errorf("Subcategory should exist in the aggregate: %v", err)
	}

	// 4. 測試業務規則：不能有重複名稱的子分類
	_, err = category.AddSubcategory(*subcategoryName1)
	if err == nil {
		t.Error("Should not allow duplicate subcategory names")
	}

	// 5. 新增第二個子分類
	subcategoryName2, err := model.NewCategoryName("午餐")
	if err != nil {
		t.Fatalf("Failed to create subcategory name: %v", err)
	}

	subcategory2, err := category.AddSubcategory(*subcategoryName2)
	if err != nil {
		t.Fatalf("Failed to add second subcategory: %v", err)
	}

	// 6. 驗證聚合狀態
	if !category.HasSubcategories() {
		t.Error("Category should have subcategories")
	}

	subcategories := category.GetAllSubcategories()
	if len(subcategories) != 2 {
		t.Errorf("Expected 2 subcategories, got %d", len(subcategories))
	}

	// 7. 測試移除子分類
	err = category.RemoveSubcategory(subcategory1.ID)
	if err != nil {
		t.Errorf("Failed to remove subcategory: %v", err)
	}

	// 8. 驗證移除後狀態
	err = category.ValidateSubcategoryExists(subcategory1.ID)
	if err == nil {
		t.Error("Subcategory should not exist after removal")
	}

	subcategoriesAfterRemoval := category.GetAllSubcategories()
	if len(subcategoriesAfterRemoval) != 1 {
		t.Errorf("Expected 1 subcategory after removal, got %d", len(subcategoriesAfterRemoval))
	}

	// 9. 測試更新子分類名稱
	newName, err := model.NewCategoryName("晚餐")
	if err != nil {
		t.Fatalf("Failed to create new name: %v", err)
	}

	err = category.UpdateSubcategoryName(subcategory2.ID, *newName)
	if err != nil {
		t.Errorf("Failed to update subcategory name: %v", err)
	}

	// 10. 驗證更新後的子分類
	updatedSubcategory, err := category.GetSubcategory(subcategory2.ID)
	if err != nil {
		t.Errorf("Failed to get updated subcategory: %v", err)
	}

	if !updatedSubcategory.Name.Equals(*newName) {
		t.Errorf("Expected updated name '晚餐', got '%s'", updatedSubcategory.Name.Value)
	}
}

// TestExpenseRecordWithSubcategoryID 測試 ExpenseRecord 使用 SubcategoryID
func TestExpenseRecordWithSubcategoryID(t *testing.T) {
	// 1. 建立子分類 ID (在實際應用中，這會來自 Category 聚合)
	subcategoryID := "subcategory-123"

	// 2. 建立錢包
	wallet, err := model.NewWallet("user123", "測試錢包", model.WalletTypeCash, "USD")
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	// 3. 先新增收入以便有餘額
	incomeAmount, err := model.NewMoney(10000, "USD") // $100.00
	if err != nil {
		t.Fatalf("Failed to create income amount: %v", err)
	}

	_, err = wallet.AddIncome(*incomeAmount, "income-subcategory-123", "測試收入", wallet.CreatedAt)
	if err != nil {
		t.Fatalf("Failed to add income: %v", err)
	}

	// 4. 新增支出記錄
	expenseAmount, err := model.NewMoney(2500, "USD") // $25.00
	if err != nil {
		t.Fatalf("Failed to create expense amount: %v", err)
	}

	expense, err := wallet.AddExpense(*expenseAmount, subcategoryID, "測試支出", wallet.CreatedAt)
	if err != nil {
		t.Fatalf("Failed to add expense: %v", err)
	}

	// 5. 驗證支出記錄
	if expense.SubcategoryID != subcategoryID {
		t.Errorf("Expected subcategory ID '%s', got '%s'", subcategoryID, expense.SubcategoryID)
	}

	if expense.Amount.Amount != 2500 {
		t.Errorf("Expected amount 2500, got %d", expense.Amount.Amount)
	}

	// 6. 驗證錢包餘額更新
	expectedBalance := int64(7500) // $100.00 - $25.00 = $75.00
	if wallet.Balance.Amount != expectedBalance {
		t.Errorf("Expected balance %d, got %d", expectedBalance, wallet.Balance.Amount)
	}
}

// TestCategoryAggregateIntegrity 測試聚合的完整性
func TestCategoryAggregateIntegrity(t *testing.T) {
	// 1. 創建分類聚合
	categoryName, _ := model.NewCategoryName("交通")
	category, _ := model.NewExpenseCategory("user123", *categoryName)

	// 2. 透過聚合根新增子分類
	busName, _ := model.NewCategoryName("公車")
	busSubcategory, _ := category.AddSubcategory(*busName)

	taxiName, _ := model.NewCategoryName("計程車")
	taxiSubcategory, _ := category.AddSubcategory(*taxiName)

	// 3. 驗證聚合的不變量
	// 所有子分類都必須屬於這個聚合
	allSubcategories := category.GetAllSubcategories()
	
	foundBus := false
	foundTaxi := false
	
	for _, sub := range allSubcategories {
		if sub.ID == busSubcategory.ID {
			foundBus = true
		}
		if sub.ID == taxiSubcategory.ID {
			foundTaxi = true
		}
	}

	if !foundBus {
		t.Error("Bus subcategory should be found in the aggregate")
	}
	if !foundTaxi {
		t.Error("Taxi subcategory should be found in the aggregate")
	}

	// 4. 驗證子分類名稱唯一性
	_, err := category.AddSubcategory(*busName)
	if err == nil {
		t.Error("Should not allow duplicate subcategory names within the same aggregate")
	}

	// 5. 驗證聚合邊界
	// 子分類不能獨立存在，必須透過聚合根操作
	// 這在設計上已經通過私有建構函式來保證
	t.Log("Category aggregate integrity test passed")
}