package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ExpenseRecord struct {
	ID            string
	WalletID      string
	SubcategoryID string    // 指向 ExpenseSubcategory.ID
	Amount        Money
	Description   string
	Date          time.Time
	CreatedAt     time.Time
}

func NewExpenseRecord(walletID, subcategoryID string, amount Money, description string, date time.Time) (*ExpenseRecord, error) {
	if walletID == "" {
		return nil, errors.New("wallet ID cannot be empty")
	}
	if subcategoryID == "" {
		return nil, errors.New("subcategory ID cannot be empty")
	}
	if amount.Amount <= 0 {
		return nil, errors.New("expense amount must be positive")
	}

	return &ExpenseRecord{
		ID:            uuid.NewString(),
		WalletID:      walletID,
		SubcategoryID: subcategoryID,
		Amount:        amount,
		Description:   description,
		Date:          date,
		CreatedAt:     time.Now(),
	}, nil
}

type IncomeRecord struct {
	ID            string
	WalletID      string
	SubcategoryID string    // 指向 IncomeSubcategory.ID
	Amount        Money
	Description   string
	Date          time.Time
	CreatedAt     time.Time
}

func NewIncomeRecord(walletID, subcategoryID string, amount Money, description string, date time.Time) (*IncomeRecord, error) {
	if walletID == "" {
		return nil, errors.New("wallet ID cannot be empty")
	}
	if subcategoryID == "" {
		return nil, errors.New("subcategory ID cannot be empty")
	}
	if amount.Amount <= 0 {
		return nil, errors.New("income amount must be positive")
	}

	return &IncomeRecord{
		ID:            uuid.NewString(),
		WalletID:      walletID,
		SubcategoryID: subcategoryID,
		Amount:        amount,
		Description:   description,
		Date:          date,
		CreatedAt:     time.Now(),
	}, nil
}

type Transfer struct {
	ID           string
	FromWalletID string
	ToWalletID   string
	Amount       Money
	Fee          Money
	Description  string
	Date         time.Time
	CreatedAt    time.Time
}

func NewTransfer(fromWalletID, toWalletID string, amount Money, fee Money, description string, date time.Time) (*Transfer, error) {
	if fromWalletID == "" {
		return nil, errors.New("from wallet ID cannot be empty")
	}
	if toWalletID == "" {
		return nil, errors.New("to wallet ID cannot be empty")
	}
	if fromWalletID == toWalletID {
		return nil, errors.New("cannot transfer to the same wallet")
	}
	if amount.Amount <= 0 {
		return nil, errors.New("transfer amount must be positive")
	}
	if fee.Amount < 0 {
		return nil, errors.New("transfer fee cannot be negative")
	}

	return &Transfer{
		ID:           uuid.NewString(),
		FromWalletID: fromWalletID,
		ToWalletID:   toWalletID,
		Amount:       amount,
		Fee:          fee,
		Description:  description,
		Date:         date,
		CreatedAt:    time.Now(),
	}, nil
}
