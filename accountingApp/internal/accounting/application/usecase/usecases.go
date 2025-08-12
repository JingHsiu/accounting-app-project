package usecase

import (
	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/query"
)

// Command Use Case Interfaces

// CreateWalletUseCase defines the interface for creating wallets
type CreateWalletUseCase interface {
	Execute(input command.CreateWalletInput) common.Output
}

// AddExpenseUseCase defines the interface for adding expenses
type AddExpenseUseCase interface {
	Execute(input command.AddExpenseInput) common.Output
}

// AddIncomeUseCase defines the interface for adding income
type AddIncomeUseCase interface {
	Execute(input command.AddIncomeInput) common.Output
}

// CreateExpenseCategoryUseCase defines the interface for creating expense categories
type CreateExpenseCategoryUseCase interface {
	Execute(input command.CreateExpenseCategoryInput) common.Output
}

// CreateIncomeCategoryUseCase defines the interface for creating income categories
type CreateIncomeCategoryUseCase interface {
	Execute(input command.CreateIncomeCategoryInput) common.Output
}

// Query Use Case Interfaces

// GetWalletBalanceUseCase defines the interface for querying wallet balance
type GetWalletBalanceUseCase interface {
	Execute(input query.GetWalletBalanceInput) common.Output
}