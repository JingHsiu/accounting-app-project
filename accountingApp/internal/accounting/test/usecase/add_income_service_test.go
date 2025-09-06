package usecase

import (
	"testing"
	"time"

	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
	"github.com/JingHsiu/accountingApp/internal/accounting/test"
	"github.com/stretchr/testify/assert"
)

// Test helper functions
func createTestWalletInRepo(repo *test.FakeWalletRepo, userID, currency string, initialBalance int64) *model.Wallet {
	wallet, _ := model.NewWalletWithInitialBalance(userID, "Test Wallet", model.WalletTypeCash, currency, initialBalance)
	repo.Save(wallet)
	return wallet
}

func createTestIncomeCategoryWithSubcategory(repo repository.IncomeCategoryRepository, userID, categoryName, subcategoryName string) (*model.IncomeCategory, *model.IncomeSubcategory) {
	catName, _ := model.NewCategoryName(categoryName)
	category, _ := model.NewIncomeCategory(userID, *catName)
	
	subCatName, _ := model.NewCategoryName(subcategoryName)
	subcategory, _ := category.AddSubcategory(*subCatName)
	
	repo.Save(category)
	return category, subcategory
}

func createAddIncomeInput(walletID, subcategoryID string, amount int64, currency, description string) usecase.AddIncomeInput {
	return usecase.AddIncomeInput{
		WalletID:      walletID,
		SubcategoryID: subcategoryID,
		Amount:        amount,
		Currency:      currency,
		Description:   description,
		Date:          time.Now(),
	}
}

// SUCCESS SCENARIOS

func Test_AddIncomeService_Success_BasicIncome(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 1000)
	_, subcategory := createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Salary", "Monthly Salary")

	input := createAddIncomeInput(wallet.ID, subcategory.ID, 500, "USD", "Monthly salary payment")

	// Act
	output := service.Execute(input)

	// Assert
	assert.Equal(t, common.Success, output.GetExitCode())
	assert.NotEmpty(t, output.GetID(), "Income ID should be generated")
	assert.Empty(t, output.GetMessage(), "Success should have no error message")

	// Verify wallet balance updated correctly
	updatedWallet, err := walletRepo.FindByID(wallet.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedWallet)
	assert.Equal(t, int64(1500), updatedWallet.Balance.Amount, "Balance should increase by 500")
	assert.Equal(t, "USD", updatedWallet.Balance.Currency)
}

func Test_AddIncomeService_Success_ZeroInitialBalance(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 0)
	_, subcategory := createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Freelance", "Project Payment")

	input := createAddIncomeInput(wallet.ID, subcategory.ID, 1200, "USD", "Project completion payment")

	// Act
	output := service.Execute(input)

	// Assert
	assert.Equal(t, common.Success, output.GetExitCode())
	
	// Verify wallet balance updated correctly
	updatedWallet, err := walletRepo.FindByID(wallet.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1200), updatedWallet.Balance.Amount, "Balance should be exactly 1200")
}

func Test_AddIncomeService_Success_LargeAmount(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 50000)
	_, subcategory := createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Investment", "Stock Dividend")

	input := createAddIncomeInput(wallet.ID, subcategory.ID, 1000000, "USD", "Large stock dividend")

	// Act
	output := service.Execute(input)

	// Assert
	assert.Equal(t, common.Success, output.GetExitCode())
	
	// Verify wallet balance updated correctly
	updatedWallet, err := walletRepo.FindByID(wallet.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1050000), updatedWallet.Balance.Amount, "Balance should handle large amounts")
}

func Test_AddIncomeService_Success_WithoutCategoryRepo(t *testing.T) {
	// Arrange - Test scenario where categoryRepo is nil (skips validation)
	walletRepo, _ := test.NewFakeWalletRepo()
	service := command.NewAddIncomeService(walletRepo, nil)

	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 1000)
	input := createAddIncomeInput(wallet.ID, "any-subcategory-id", 300, "USD", "Income without category validation")

	// Act
	output := service.Execute(input)

	// Assert
	assert.Equal(t, common.Success, output.GetExitCode())
	
	// Verify wallet balance updated correctly
	updatedWallet, err := walletRepo.FindByID(wallet.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1300), updatedWallet.Balance.Amount)
}

// FAILURE SCENARIOS - WALLET ISSUES

func Test_AddIncomeService_Failure_WalletNotFound(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	input := createAddIncomeInput("nonexistent-wallet-id", "subcategory-id", 100, "USD", "Test income")

	// Act
	output := service.Execute(input)

	// Assert
	assert.Equal(t, common.Failure, output.GetExitCode())
	assert.Contains(t, output.GetMessage(), "Wallet not found")
	assert.Empty(t, output.GetID(), "No income ID should be generated on failure")
}

func Test_AddIncomeService_Failure_EmptyWalletID(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	input := createAddIncomeInput("", "subcategory-id", 100, "USD", "Test income")

	// Act
	output := service.Execute(input)

	// Assert
	assert.Equal(t, common.Failure, output.GetExitCode())
	assert.Contains(t, output.GetMessage(), "Wallet not found")
}

// FAILURE SCENARIOS - CATEGORY ISSUES

func Test_AddIncomeService_Failure_SubcategoryNotFound(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 1000)
	input := createAddIncomeInput(wallet.ID, "nonexistent-subcategory-id", 100, "USD", "Test income")

	// Act
	output := service.Execute(input)

	// Assert
	assert.Equal(t, common.Failure, output.GetExitCode())
	assert.Contains(t, output.GetMessage(), "Subcategory not found in any category")
	assert.Empty(t, output.GetID())
}

func Test_AddIncomeService_Failure_InvalidSubcategoryID(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 1000)
	_, _ = createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Salary", "Monthly Salary")
	
	// Try to use a subcategory ID that doesn't exist in this category
	input := createAddIncomeInput(wallet.ID, "invalid-subcategory-id", 100, "USD", "Test income")

	// Act
	output := service.Execute(input)

	// Assert
	assert.Equal(t, common.Failure, output.GetExitCode())
	assert.Contains(t, output.GetMessage(), "Subcategory not found")
}

func Test_AddIncomeService_Failure_EmptySubcategoryID(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 1000)
	input := createAddIncomeInput(wallet.ID, "", 100, "USD", "Test income")

	// Act
	output := service.Execute(input)

	// Assert
	assert.Equal(t, common.Failure, output.GetExitCode())
	assert.Contains(t, output.GetMessage(), "Failed to find category for subcategory")
}

// FAILURE SCENARIOS - AMOUNT AND CURRENCY ISSUES

func Test_AddIncomeService_Failure_NegativeAmount(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 1000)
	_, subcategory := createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Salary", "Monthly Salary")

	input := createAddIncomeInput(wallet.ID, subcategory.ID, -100, "USD", "Negative amount test")

	// Act
	output := service.Execute(input)

	// Assert
	assert.Equal(t, common.Failure, output.GetExitCode())
	assert.Contains(t, output.GetMessage(), "Invalid amount")
	assert.Contains(t, output.GetMessage(), "amount cannot be negative")
	assert.Empty(t, output.GetID())
}

func Test_AddIncomeService_Failure_ZeroAmount(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 1000)
	_, subcategory := createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Salary", "Monthly Salary")

	input := createAddIncomeInput(wallet.ID, subcategory.ID, 0, "USD", "Zero amount test")

	// Act
	output := service.Execute(input)

	// Assert
	assert.Equal(t, common.Failure, output.GetExitCode())
	assert.Contains(t, output.GetMessage(), "Adding income failed")
	assert.Contains(t, output.GetMessage(), "income amount must be positive")
}

func Test_AddIncomeService_Failure_EmptyCurrency(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 1000)
	_, subcategory := createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Salary", "Monthly Salary")

	input := createAddIncomeInput(wallet.ID, subcategory.ID, 100, "", "Empty currency test")

	// Act
	output := service.Execute(input)

	// Assert
	assert.Equal(t, common.Failure, output.GetExitCode())
	assert.Contains(t, output.GetMessage(), "Invalid amount")
	assert.Contains(t, output.GetMessage(), "currency cannot be empty")
}

func Test_AddIncomeService_Failure_InvalidCurrencyLength(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 1000)
	_, subcategory := createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Salary", "Monthly Salary")

	testCases := []struct {
		currency string
		name     string
	}{
		{"US", "too short currency"},
		{"USDD", "too long currency"},
		{"U", "single character currency"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := createAddIncomeInput(wallet.ID, subcategory.ID, 100, tc.currency, "Invalid currency test")

			// Act
			output := service.Execute(input)

			// Assert
			assert.Equal(t, common.Failure, output.GetExitCode())
			assert.Contains(t, output.GetMessage(), "Invalid amount")
			assert.Contains(t, output.GetMessage(), "currency must be 3 characters")
		})
	}
}

func Test_AddIncomeService_Failure_CurrencyMismatch(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	// Create wallet with USD currency
	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 1000)
	_, subcategory := createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Salary", "Monthly Salary")

	// Try to add income with EUR currency
	input := createAddIncomeInput(wallet.ID, subcategory.ID, 100, "EUR", "Currency mismatch test")

	// Act
	output := service.Execute(input)

	// Assert
	assert.Equal(t, common.Failure, output.GetExitCode())
	assert.Contains(t, output.GetMessage(), "Adding income failed")
	assert.Contains(t, output.GetMessage(), "currency EUR does not match wallet currency USD")
}

// EDGE CASES

func Test_AddIncomeService_EdgeCase_VeryLongDescription(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 1000)
	_, subcategory := createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Salary", "Monthly Salary")

	longDescription := "This is a very long description that might test the limits of the system. " +
		"It contains many words and characters to see how the system handles long text input. " +
		"We want to ensure that even with very long descriptions, the income can still be added successfully " +
		"without any issues or truncation problems that might affect the core functionality."

	input := createAddIncomeInput(wallet.ID, subcategory.ID, 500, "USD", longDescription)

	// Act
	output := service.Execute(input)

	// Assert
	assert.Equal(t, common.Success, output.GetExitCode())
	
	// Verify wallet balance updated correctly
	updatedWallet, err := walletRepo.FindByID(wallet.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1500), updatedWallet.Balance.Amount)
}

func Test_AddIncomeService_EdgeCase_EmptyDescription(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 1000)
	_, subcategory := createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Salary", "Monthly Salary")

	input := createAddIncomeInput(wallet.ID, subcategory.ID, 500, "USD", "")

	// Act
	output := service.Execute(input)

	// Assert
	assert.Equal(t, common.Success, output.GetExitCode())
	
	// Verify wallet balance updated correctly
	updatedWallet, err := walletRepo.FindByID(wallet.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1500), updatedWallet.Balance.Amount)
}

func Test_AddIncomeService_EdgeCase_MultipleIncomes_SameWallet(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 1000)
	_, subcategory1 := createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Salary", "Monthly Salary")
	_, subcategory2 := createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Freelance", "Project Work")

	// Act - Add multiple incomes
	input1 := createAddIncomeInput(wallet.ID, subcategory1.ID, 2000, "USD", "First income")
	output1 := service.Execute(input1)

	input2 := createAddIncomeInput(wallet.ID, subcategory2.ID, 1500, "USD", "Second income")
	output2 := service.Execute(input2)

	input3 := createAddIncomeInput(wallet.ID, subcategory1.ID, 800, "USD", "Third income")
	output3 := service.Execute(input3)

	// Assert
	assert.Equal(t, common.Success, output1.GetExitCode())
	assert.Equal(t, common.Success, output2.GetExitCode())
	assert.Equal(t, common.Success, output3.GetExitCode())
	
	// Verify final wallet balance (1000 + 2000 + 1500 + 800 = 5300)
	updatedWallet, err := walletRepo.FindByID(wallet.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(5300), updatedWallet.Balance.Amount)
}

// INTEGRATION TESTS

func Test_AddIncomeService_Integration_CompleteWorkflow(t *testing.T) {
	// Arrange
	walletRepo, _ := test.NewFakeWalletRepo()
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	service := command.NewAddIncomeService(walletRepo, categoryRepo)

	// Create user's wallet and categories
	wallet := createTestWalletInRepo(walletRepo, "user-123", "USD", 5000)
	_, salarySubcat := createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Employment", "Full Time Salary")
	_, bonusSubcat := createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Employment", "Performance Bonus")
	_, freelanceSubcat := createTestIncomeCategoryWithSubcategory(categoryRepo, "user-123", "Freelance", "Web Development")

	initialBalance := wallet.Balance.Amount

	// Act - Add various types of income
	incomes := []struct {
		subcategoryID string
		amount        int64
		description   string
	}{
		{salarySubcat.ID, 4000, "Monthly salary - January"},
		{bonusSubcat.ID, 1000, "Q4 performance bonus"},
		{freelanceSubcat.ID, 2500, "Website development project"},
		{salarySubcat.ID, 4000, "Monthly salary - February"},
	}

	var totalIncomeAdded int64
	for _, income := range incomes {
		input := createAddIncomeInput(wallet.ID, income.subcategoryID, income.amount, "USD", income.description)
		output := service.Execute(input)
		
		assert.Equal(t, common.Success, output.GetExitCode(), "All incomes should be added successfully")
		assert.NotEmpty(t, output.GetID(), "Each income should generate an ID")
		
		totalIncomeAdded += income.amount
	}

	// Assert - Verify final state
	finalWallet, err := walletRepo.FindByID(wallet.ID)
	assert.NoError(t, err)
	assert.NotNil(t, finalWallet)
	
	expectedFinalBalance := initialBalance + totalIncomeAdded
	assert.Equal(t, expectedFinalBalance, finalWallet.Balance.Amount, 
		"Final balance should equal initial balance plus all added incomes")
	
	assert.Equal(t, "USD", finalWallet.Balance.Currency, "Currency should remain unchanged")
}