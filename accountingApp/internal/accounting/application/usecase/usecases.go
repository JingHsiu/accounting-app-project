package usecase

import (
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
	"time"
)

// =============================================================================
// INPUT/OUTPUT CONTRACTS
// =============================================================================

// Command Inputs
type CreateWalletInput struct {
	UserID         string
	Name           string
	Type           string
	Currency       string
	InitialBalance *int64 // Optional initial balance in cents/smallest currency unit
}

type AddExpenseInput struct {
	WalletID      string
	SubcategoryID string
	Amount        int64
	Currency      string
	Description   string
	Date          time.Time
}

type AddIncomeInput struct {
	WalletID      string
	SubcategoryID string
	Amount        int64
	Currency      string
	Description   string
	Date          time.Time
}

type CreateExpenseCategoryInput struct {
	UserID string
	Name   string
}

type CreateIncomeCategoryInput struct {
	UserID string
	Name   string
}

type UpdateWalletInput struct {
	WalletID string
	Name     *string // Optional - only update if provided
	Type     *string // Optional - only update if provided
	Currency *string // Optional - only update if provided (note: currency changes are complex)
}

type DeleteWalletInput struct {
	WalletID string
}

// Query Inputs
type GetWalletInput struct {
	WalletID            string
	IncludeTransactions bool
}

type GetWalletBalanceInput struct {
	WalletID string
}

type GetWalletsInput struct {
	UserID string
}

type GetExpenseCategoriesInput struct {
	UserID string
}

type GetIncomeCategoriesInput struct {
	UserID string
}

type GetIncomesInput struct {
	UserID       string
	WalletID     *string // Optional filter
	CategoryID   *string // Optional filter
	StartDate    *time.Time // Optional date range filter
	EndDate      *time.Time // Optional date range filter
	MinAmount    *int64  // Optional amount range filter (in cents)
	MaxAmount    *int64  // Optional amount range filter (in cents)
	Description  *string // Optional description search filter
}

type GetExpensesInput struct {
	UserID       string
	WalletID     *string // Optional filter
	CategoryID   *string // Optional filter
	StartDate    *time.Time // Optional date range filter
	EndDate      *time.Time // Optional date range filter
	MinAmount    *int64  // Optional amount range filter (in cents)
	MaxAmount    *int64  // Optional amount range filter (in cents)
	Description  *string // Optional description search filter
}

// Query Outputs (specialized outputs for queries that return data)
type GetWalletOutput struct {
	ID       string          `json:"id"`
	ExitCode common.ExitCode `json:"exit_code"`
	Message  string          `json:"message"`
	Wallet   *model.Wallet   `json:"wallet,omitempty"`
}

func (o GetWalletOutput) GetID() string                { return o.ID }
func (o GetWalletOutput) GetExitCode() common.ExitCode { return o.ExitCode }
func (o GetWalletOutput) GetMessage() string           { return o.Message }

type GetWalletBalanceOutput struct {
	ID       string          `json:"id"`
	ExitCode common.ExitCode `json:"exit_code"`
	Message  string          `json:"message"`
	Balance  string          `json:"balance,omitempty"`
	Currency string          `json:"currency,omitempty"`
}

func (o GetWalletBalanceOutput) GetID() string                { return o.ID }
func (o GetWalletBalanceOutput) GetExitCode() common.ExitCode { return o.ExitCode }
func (o GetWalletBalanceOutput) GetMessage() string           { return o.Message }

type GetWalletsOutput struct {
	ID       string          `json:"id"`
	ExitCode common.ExitCode `json:"exit_code"`
	Message  string          `json:"message"`
	Wallets  []*model.Wallet `json:"wallets,omitempty"`
}

func (o GetWalletsOutput) GetID() string                { return o.ID }
func (o GetWalletsOutput) GetExitCode() common.ExitCode { return o.ExitCode }
func (o GetWalletsOutput) GetMessage() string           { return o.Message }

// Category structure for API responses
type CategoryData struct {
	ID            string                   `json:"id"`
	Name          string                   `json:"name"`
	Type          string                   `json:"type"` // "expense" or "income"
	CreatedAt     string                   `json:"created_at"`
	UpdatedAt     string                   `json:"updated_at"`
	Subcategories []SubcategoryData        `json:"subcategories,omitempty"`
}

type SubcategoryData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Income record structure for API responses
type IncomeRecordData struct {
	ID            string `json:"id"`
	WalletID      string `json:"wallet_id"`
	SubcategoryID string `json:"subcategory_id"`
	Amount        struct {
		Amount   int64  `json:"amount"`   // Amount in cents
		Currency string `json:"currency"`
	} `json:"amount"`
	Description string `json:"description"`
	Date        string `json:"date"`        // ISO format
	CreatedAt   string `json:"created_at"`  // ISO format
}

// Expense record structure for API responses
type ExpenseRecordData struct {
	ID            string `json:"id"`
	WalletID      string `json:"wallet_id"`
	SubcategoryID string `json:"subcategory_id"`
	Amount        struct {
		Amount   int64  `json:"amount"`   // Amount in cents
		Currency string `json:"currency"`
	} `json:"amount"`
	Description string `json:"description"`
	Date        string `json:"date"`        // ISO format
	CreatedAt   string `json:"created_at"`  // ISO format
}

type GetExpenseCategoriesOutput struct {
	ID         string          `json:"id"`
	ExitCode   common.ExitCode `json:"exit_code"`
	Message    string          `json:"message"`
	Categories []CategoryData  `json:"categories,omitempty"`
}

func (o GetExpenseCategoriesOutput) GetID() string                { return o.ID }
func (o GetExpenseCategoriesOutput) GetExitCode() common.ExitCode { return o.ExitCode }
func (o GetExpenseCategoriesOutput) GetMessage() string           { return o.Message }

type GetIncomeCategoriesOutput struct {
	ID         string          `json:"id"`
	ExitCode   common.ExitCode `json:"exit_code"`
	Message    string          `json:"message"`
	Categories []CategoryData  `json:"categories,omitempty"`
}

func (o GetIncomeCategoriesOutput) GetID() string                { return o.ID }
func (o GetIncomeCategoriesOutput) GetExitCode() common.ExitCode { return o.ExitCode }
func (o GetIncomeCategoriesOutput) GetMessage() string           { return o.Message }

type GetIncomesOutput struct {
	ID      string             `json:"id"`
	ExitCode common.ExitCode   `json:"exit_code"`
	Message string            `json:"message"`
	Data    []IncomeRecordData `json:"data,omitempty"`
	Count   int               `json:"count"`
}

func (o GetIncomesOutput) GetID() string                { return o.ID }
func (o GetIncomesOutput) GetExitCode() common.ExitCode { return o.ExitCode }
func (o GetIncomesOutput) GetMessage() string           { return o.Message }

type GetExpensesOutput struct {
	ID      string              `json:"id"`
	ExitCode common.ExitCode    `json:"exit_code"`
	Message string             `json:"message"`
	Data    []ExpenseRecordData `json:"data,omitempty"`
	Count   int                `json:"count"`
}

func (o GetExpensesOutput) GetID() string                { return o.ID }
func (o GetExpensesOutput) GetExitCode() common.ExitCode { return o.ExitCode }
func (o GetExpensesOutput) GetMessage() string           { return o.Message }

// =============================================================================
// USE CASE INTERFACES
// =============================================================================

// Command Use Case Interfaces

// CreateWalletUseCase defines the interface for creating wallets
type CreateWalletUseCase interface {
	Execute(input CreateWalletInput) common.Output
}

// AddExpenseUseCase defines the interface for adding expenses
type AddExpenseUseCase interface {
	Execute(input AddExpenseInput) common.Output
}

// AddIncomeUseCase defines the interface for adding income
type AddIncomeUseCase interface {
	Execute(input AddIncomeInput) common.Output
}

// CreateExpenseCategoryUseCase defines the interface for creating expense categories
type CreateExpenseCategoryUseCase interface {
	Execute(input CreateExpenseCategoryInput) common.Output
}

// CreateIncomeCategoryUseCase defines the interface for creating income categories
type CreateIncomeCategoryUseCase interface {
	Execute(input CreateIncomeCategoryInput) common.Output
}

// UpdateWalletUseCase defines the interface for updating wallet information
type UpdateWalletUseCase interface {
	Execute(input UpdateWalletInput) common.Output
}

// DeleteWalletUseCase defines the interface for deleting wallets
type DeleteWalletUseCase interface {
	Execute(input DeleteWalletInput) common.Output
}

// Query Use Case Interfaces

// GetWalletBalanceUseCase defines the interface for querying wallet balance
type GetWalletBalanceUseCase interface {
	Execute(input GetWalletBalanceInput) common.Output
}

// GetWalletsUseCase defines the interface for querying user's wallets
type GetWalletsUseCase interface {
	Execute(input GetWalletsInput) common.Output
}

// GetWalletUseCase defines the interface for querying a single wallet
type GetWalletUseCase interface {
	Execute(input GetWalletInput) common.Output
}

// GetExpenseCategoriesUseCase defines the interface for querying user's expense categories
type GetExpenseCategoriesUseCase interface {
	Execute(input GetExpenseCategoriesInput) common.Output
}

// GetIncomeCategoriesUseCase defines the interface for querying user's income categories
type GetIncomeCategoriesUseCase interface {
	Execute(input GetIncomeCategoriesInput) common.Output
}

// GetIncomesUseCase defines the interface for querying income records
type GetIncomesUseCase interface {
	Execute(input GetIncomesInput) common.Output
}

// GetExpensesUseCase defines the interface for querying expense records
type GetExpensesUseCase interface {
	Execute(input GetExpensesInput) common.Output
}