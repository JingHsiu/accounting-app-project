package test

import (
	"testing"

	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

// TestWalletCurrencyRefactoring validates the Currency field removal and getter method
func TestWalletCurrencyRefactoring(t *testing.T) {
	// Test wallet creation
	wallet, err := model.NewWallet("user123", "Test Wallet", model.WalletTypeCash, "USD")
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	// Test Currency() method returns correct currency
	if wallet.Currency() != "USD" {
		t.Errorf("Expected currency 'USD', got '%s'", wallet.Currency())
	}

	// Test that currency comes from Balance
	if wallet.Currency() != wallet.Balance.Currency {
		t.Errorf("Currency() method should return Balance.Currency. Got Currency()='%s', Balance.Currency='%s'", 
			wallet.Currency(), wallet.Balance.Currency)
	}

	// Test currency validation in business methods
	amount, err := model.NewMoney(1000, "EUR") // Different currency
	if err != nil {
		t.Fatalf("Failed to create money: %v", err)
	}

	// Should fail because currency doesn't match
	_, err = wallet.AddExpense(*amount, "cat123", "Test expense", wallet.CreatedAt)
	if err == nil {
		t.Error("Expected error when adding expense with different currency")
	}

	// First add some income so we have balance for expense
	incomeAmount, err := model.NewMoney(1000, "USD")
	if err != nil {
		t.Fatalf("Failed to create income money: %v", err)
	}

	_, err = wallet.AddIncome(*incomeAmount, "income123", "Test income", wallet.CreatedAt)
	if err != nil {
		t.Errorf("Failed to add income: %v", err)
	}

	// Now should work with same currency expense
	usdAmount, err := model.NewMoney(500, "USD")
	if err != nil {
		t.Fatalf("Failed to create USD money: %v", err)
	}

	expense, err := wallet.AddExpense(*usdAmount, "cat123", "Test expense", wallet.CreatedAt)
	if err != nil {
		t.Errorf("Failed to add expense with same currency: %v", err)
	}

	if expense == nil {
		t.Error("Expected expense to be created")
	}

	// Verify wallet currency is still consistent
	if wallet.Currency() != "USD" {
		t.Errorf("Wallet currency should remain 'USD', got '%s'", wallet.Currency())
	}
}

// TestWalletCurrencyFromBalance ensures currency always comes from balance
func TestWalletCurrencyFromBalance(t *testing.T) {
	wallet, err := model.NewWallet("user123", "Test Wallet", model.WalletTypeBank, "EUR")
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	// Currency should be EUR
	if wallet.Currency() != "EUR" {
		t.Errorf("Expected currency 'EUR', got '%s'", wallet.Currency())
	}

	// Add income to change balance but currency should remain same
	income, err := model.NewMoney(2000, "EUR")
	if err != nil {
		t.Fatalf("Failed to create money: %v", err)
	}

	_, err = wallet.AddIncome(*income, "income123", "Test income", wallet.CreatedAt)
	if err != nil {
		t.Errorf("Failed to add income: %v", err)
	}

	// Currency should still be EUR
	if wallet.Currency() != "EUR" {
		t.Errorf("Expected currency to remain 'EUR', got '%s'", wallet.Currency())
	}

	// Balance should have EUR
	if wallet.Balance.Currency != "EUR" {
		t.Errorf("Expected balance currency 'EUR', got '%s'", wallet.Balance.Currency)
	}
}