package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type WalletType string

const (
	WalletTypeCash       WalletType = "CASH"
	WalletTypeBank       WalletType = "BANK"
	WalletTypeCredit     WalletType = "CREDIT"
	WalletTypeInvestment WalletType = "INVESTMENT"
)

func ParseWalletType(s string) (WalletType, error) {
	switch WalletType(s) {
	case WalletTypeCash, WalletTypeBank, WalletTypeCredit, WalletTypeInvestment:
		return WalletType(s), nil
	default:
		return "", fmt.Errorf("invalid wallet type: %s", s)
	}
}

type Wallet struct {
	ID        string
	UserID    string
	Name      string
	Type      WalletType
	Balance   Money
	CreatedAt time.Time
	UpdatedAt time.Time
	
	// 內部Entities - 聚合邊界內的所有交易記錄
	expenseRecords []ExpenseRecord
	incomeRecords  []IncomeRecord
	transfers      []Transfer
	
	// 載入狀態標記
	isFullyLoaded bool // 標記是否已載入所有交易記錄
}

func NewWallet(userID, name string, walletType WalletType, currency string) (*Wallet, error) {
	return NewWalletWithInitialBalance(userID, name, walletType, currency, 0)
}

func NewWalletWithInitialBalance(userID, name string, walletType WalletType, currency string, initialBalanceAmount int64) (*Wallet, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("wallet name cannot be empty")
	}
	if currency == "" || len(currency) != 3 {
		return nil, errors.New("currency must be 3 characters (ISO 4217)")
	}
	if initialBalanceAmount < 0 {
		return nil, errors.New("initial balance cannot be negative")
	}

	initialBalance, err := NewMoney(initialBalanceAmount, currency)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Wallet{
		ID:              uuid.NewString(),
		UserID:          userID,
		Name:            strings.TrimSpace(name),
		Type:            walletType,
		Balance:         *initialBalance,
		CreatedAt:       now,
		UpdatedAt:       now,
		expenseRecords:  make([]ExpenseRecord, 0),
		incomeRecords:   make([]IncomeRecord, 0),
		transfers:       make([]Transfer, 0),
		isFullyLoaded:   false,
	}, nil
}

// The Currency returns the currency of the wallet's balance
func (w *Wallet) Currency() string {
	return w.Balance.Currency
}

// Domain Model方法 - 透過聚合根獲取資訊
func (w *Wallet) GetExpenseRecords() []ExpenseRecord {
	return w.expenseRecords
}

func (w *Wallet) GetIncomeRecords() []IncomeRecord {
	return w.incomeRecords
}

func (w *Wallet) GetTransfers() []Transfer {
	return w.transfers
}

func (w *Wallet) IsFullyLoaded() bool {
	return w.isFullyLoaded
}

func (w *Wallet) SetFullyLoaded(loaded bool) {
	w.isFullyLoaded = loaded
}

func (w *Wallet) AddExpenseRecord(record ExpenseRecord) {
	w.expenseRecords = append(w.expenseRecords, record)
}

func (w *Wallet) AddIncomeRecord(record IncomeRecord) {
	w.incomeRecords = append(w.incomeRecords, record)
}

func (w *Wallet) AddTransfer(transfer Transfer) {
	w.transfers = append(w.transfers, transfer)
}

func (w *Wallet) AddExpense(amount Money, subcategoryID, description string, date time.Time) (*ExpenseRecord, error) {
	if amount.Currency != w.Currency() {
		return nil, fmt.Errorf("expense currency %s does not match wallet currency %s", amount.Currency, w.Currency())
	}

	newBalance, err := w.Balance.Subtract(amount)
	if err != nil {
		return nil, fmt.Errorf("insufficient balance: %w", err)
	}

	expense, err := NewExpenseRecord(w.ID, subcategoryID, amount, description, date)
	if err != nil {
		return nil, err
	}

	w.Balance = *newBalance
	w.expenseRecords = append(w.expenseRecords, *expense)
	w.UpdatedAt = time.Now()
	return expense, nil
}

func (w *Wallet) AddIncome(amount Money, subcategoryID, description string, date time.Time) (*IncomeRecord, error) {
	if amount.Currency != w.Currency() {
		return nil, fmt.Errorf("income currency %s does not match wallet currency %s", amount.Currency, w.Currency())
	}

	newBalance, err := w.Balance.Add(amount)
	if err != nil {
		return nil, err
	}

	income, err := NewIncomeRecord(w.ID, subcategoryID, amount, description, date)
	if err != nil {
		return nil, err
	}

	w.Balance = *newBalance
	w.incomeRecords = append(w.incomeRecords, *income)
	w.UpdatedAt = time.Now()
	return income, nil
}

func (w *Wallet) CanTransfer(amount Money) error {
	if amount.Currency != w.Currency() {
		return fmt.Errorf("transfer currency %s does not match wallet currency %s", amount.Currency, w.Currency())
	}

	_, err := w.Balance.Subtract(amount)
	return err
}

func (w *Wallet) ProcessOutgoingTransfer(amount Money, fee Money) error {
	if amount.Currency != w.Currency() {
		return fmt.Errorf("transfer currency %s does not match wallet currency %s", amount.Currency, w.Currency())
	}
	if fee.Currency != w.Currency() {
		return fmt.Errorf("fee currency %s does not match wallet currency %s", fee.Currency, w.Currency())
	}

	totalAmount, err := amount.Add(fee)
	if err != nil {
		return err
	}

	newBalance, err := w.Balance.Subtract(*totalAmount)
	if err != nil {
		return fmt.Errorf("insufficient balance for transfer: %w", err)
	}

	w.Balance = *newBalance
	w.UpdatedAt = time.Now()
	return nil
}

// CreateTransfer 建立轉帳記錄
func (w *Wallet) CreateTransfer(toWalletID string, amount Money, fee Money, description string, date time.Time) (*Transfer, error) {
	transfer, err := NewTransfer(w.ID, toWalletID, amount, fee, description, date)
	if err != nil {
		return nil, err
	}
	
	w.transfers = append(w.transfers, *transfer)
	return transfer, nil
}

func (w *Wallet) ProcessIncomingTransfer(amount Money) error {
	if amount.Currency != w.Currency() {
		return fmt.Errorf("transfer currency %s does not match wallet currency %s", amount.Currency, w.Currency())
	}

	newBalance, err := w.Balance.Add(amount)
	if err != nil {
		return err
	}

	w.Balance = *newBalance
	w.UpdatedAt = time.Now()
	return nil
}

// Transaction 統一交易記錄介面
type Transaction struct {
	Type   string      // "expense", "income", "transfer"
	Record interface{} // ExpenseRecord, IncomeRecord, or Transfer
}

// 聚合內部查詢方法 - 避免外部Inquiry
func (w *Wallet) GetTransactionHistory(from, to time.Time) []Transaction {
	var transactions []Transaction
	
	// 收集所有交易記錄
	for _, expense := range w.expenseRecords {
		if expense.Date.After(from) && expense.Date.Before(to) {
			transactions = append(transactions, Transaction{
				Type:   "expense",
				Record: expense,
			})
		}
	}
	
	for _, income := range w.incomeRecords {
		if income.Date.After(from) && income.Date.Before(to) {
			transactions = append(transactions, Transaction{
				Type:   "income",
				Record: income,
			})
		}
	}
	
	for _, transfer := range w.transfers {
		if transfer.Date.After(from) && transfer.Date.Before(to) {
			transactions = append(transactions, Transaction{
				Type:   "transfer",
				Record: transfer,
			})
		}
	}
	
	return transactions
}

func (w *Wallet) GetMonthlyTotal(year int, month int) (expenses Money, incomes Money) {
	// 在Domain Model內計算月度總計，無需外部查詢
	for _, expense := range w.expenseRecords {
		if expense.Date.Year() == year && int(expense.Date.Month()) == month {
			if expenses.Currency == "" {
				expenses = expense.Amount
			} else {
				result, _ := expenses.Add(expense.Amount)
				if result != nil {
					expenses = *result
				}
			}
		}
	}
	
	for _, income := range w.incomeRecords {
		if income.Date.Year() == year && int(income.Date.Month()) == month {
			if incomes.Currency == "" {
				incomes = income.Amount
			} else {
				result, _ := incomes.Add(income.Amount)
				if result != nil {
					incomes = *result
				}
			}
		}
	}
	
	return expenses, incomes
}
