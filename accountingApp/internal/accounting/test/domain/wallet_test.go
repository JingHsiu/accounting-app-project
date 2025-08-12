package domain

import (
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewWallet_Success(t *testing.T) {
	wallet, err := model.NewWallet("user-123", "My Wallet", model.WalletTypeCash, "USD")

	assert.NoError(t, err)
	assert.NotEmpty(t, wallet.ID)
	assert.Equal(t, "user-123", wallet.UserID)
	assert.Equal(t, "My Wallet", wallet.Name)
	assert.Equal(t, model.WalletTypeCash, wallet.Type)
	assert.Equal(t, "USD", wallet.Currency())
	assert.Equal(t, int64(0), wallet.Balance.Amount)
	assert.Equal(t, "USD", wallet.Balance.Currency)
}

func TestNewWallet_InvalidUserID(t *testing.T) {
	wallet, err := model.NewWallet("", "My Wallet", model.WalletTypeCash, "USD")

	assert.Error(t, err)
	assert.Nil(t, wallet)
	assert.Contains(t, err.Error(), "user ID cannot be empty")
}

func TestNewWallet_InvalidName(t *testing.T) {
	wallet, err := model.NewWallet("user-123", "   ", model.WalletTypeCash, "USD")

	assert.Error(t, err)
	assert.Nil(t, wallet)
	assert.Contains(t, err.Error(), "wallet name cannot be empty")
}

func TestWallet_AddExpense_Success(t *testing.T) {
	wallet, _ := model.NewWallet("user-123", "My Wallet", model.WalletTypeCash, "USD")

	initialAmount, _ := model.NewMoney(10000, "USD")
	wallet.Balance = *initialAmount

	expenseAmount, _ := model.NewMoney(2000, "USD")
	expense, err := wallet.AddExpense(*expenseAmount, "cat-123", "Coffee", time.Now())

	assert.NoError(t, err)
	assert.NotNil(t, expense)
	assert.Equal(t, int64(8000), wallet.Balance.Amount)
	assert.Equal(t, wallet.ID, expense.WalletID)
	assert.Equal(t, "cat-123", expense.SubcategoryID)
}

func TestWallet_AddExpense_InsufficientBalance(t *testing.T) {
	wallet, _ := model.NewWallet("user-123", "My Wallet", model.WalletTypeCash, "USD")

	expenseAmount, _ := model.NewMoney(2000, "USD")
	expense, err := wallet.AddExpense(*expenseAmount, "cat-123", "Coffee", time.Now())

	assert.Error(t, err)
	assert.Nil(t, expense)
	assert.Contains(t, err.Error(), "insufficient balance")
}

func TestWallet_AddIncome_Success(t *testing.T) {
	wallet, _ := model.NewWallet("user-123", "My Wallet", model.WalletTypeCash, "USD")

	incomeAmount, _ := model.NewMoney(5000, "USD")
	income, err := wallet.AddIncome(*incomeAmount, "cat-456", "Salary", time.Now())

	assert.NoError(t, err)
	assert.NotNil(t, income)
	assert.Equal(t, int64(5000), wallet.Balance.Amount)
	assert.Equal(t, wallet.ID, income.WalletID)
	assert.Equal(t, "cat-456", income.SubcategoryID)
}

func TestWallet_ProcessOutgoingTransfer_Success(t *testing.T) {
	wallet, _ := model.NewWallet("user-123", "My Wallet", model.WalletTypeCash, "USD")

	initialAmount, _ := model.NewMoney(10000, "USD")
	wallet.Balance = *initialAmount

	transferAmount, _ := model.NewMoney(3000, "USD")
	fee, _ := model.NewMoney(100, "USD")

	err := wallet.ProcessOutgoingTransfer(*transferAmount, *fee)

	assert.NoError(t, err)
	assert.Equal(t, int64(6900), wallet.Balance.Amount) // 10000 - 3000 - 100
}

func TestWallet_ProcessIncomingTransfer_Success(t *testing.T) {
	wallet, _ := model.NewWallet("user-123", "My Wallet", model.WalletTypeCash, "USD")

	transferAmount, _ := model.NewMoney(2000, "USD")
	err := wallet.ProcessIncomingTransfer(*transferAmount)

	assert.NoError(t, err)
	assert.Equal(t, int64(2000), wallet.Balance.Amount)
}
