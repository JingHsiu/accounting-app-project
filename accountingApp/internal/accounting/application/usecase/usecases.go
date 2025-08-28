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