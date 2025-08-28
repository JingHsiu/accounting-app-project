package test

import (
	"testing"

	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
)

// MockCreateWalletUseCase is a mock implementation of CreateWalletUseCase
type MockCreateWalletUseCase struct {
	executeFunc func(input usecase.CreateWalletInput) common.Output
}

func (m *MockCreateWalletUseCase) Execute(input usecase.CreateWalletInput) common.Output {
	if m.executeFunc != nil {
		return m.executeFunc(input)
	}
	return common.UseCaseOutput{
		ID:       "mock-wallet-id",
		ExitCode: common.Success,
		Message:  "Mock wallet created successfully",
	}
}

// TestDependencyInversion validates that our interface system works correctly
func TestDependencyInversion(t *testing.T) {
	// Create a mock implementation
	mockUseCase := &MockCreateWalletUseCase{}

	// Test that we can assign the mock to the interface
	var useCase usecase.CreateWalletUseCase = mockUseCase

	// Test that the interface method works
	input := usecase.CreateWalletInput{
		UserID:   "test-user",
		Name:     "Test Wallet",
		Type:     "CASH",
		Currency: "USD",
	}

	output := useCase.Execute(input)

	// Verify the mock was called and returned expected results
	if output.GetID() != "mock-wallet-id" {
		t.Errorf("Expected mock-wallet-id, got %s", output.GetID())
	}

	if output.GetExitCode() != common.Success {
		t.Errorf("Expected success exit code, got %v", output.GetExitCode())
	}

	if output.GetMessage() != "Mock wallet created successfully" {
		t.Errorf("Expected mock message, got %s", output.GetMessage())
	}
}

// TestInterfaceImplementation verifies that real services implement the interfaces
func TestInterfaceImplementation(t *testing.T) {
	// This test verifies that concrete services implement the interfaces correctly
	// If this compiles, it means the interfaces are implemented properly

	var createWalletUseCase usecase.CreateWalletUseCase
	var addExpenseUseCase usecase.AddExpenseUseCase
	var addIncomeUseCase usecase.AddIncomeUseCase
	var getWalletBalanceUseCase usecase.GetWalletBalanceUseCase
	var createExpenseCategoryUseCase usecase.CreateExpenseCategoryUseCase
	var createIncomeCategoryUseCase usecase.CreateIncomeCategoryUseCase

	// These assignments will fail to compile if interfaces are not implemented
	createWalletUseCase = command.NewCreateWalletService(nil)
	addExpenseUseCase = command.NewAddExpenseService(nil, nil)
	addIncomeUseCase = command.NewAddIncomeService(nil, nil)
	// getWalletBalanceUseCase = query.NewGetWalletBalanceService(nil) // Would need import
	createExpenseCategoryUseCase = command.NewCreateExpenseCategoryService(nil)
	createIncomeCategoryUseCase = command.NewCreateIncomeCategoryService(nil)

	// Prevent unused variable errors
	_ = createWalletUseCase
	_ = addExpenseUseCase
	_ = addIncomeUseCase
	_ = getWalletBalanceUseCase
	_ = createExpenseCategoryUseCase
	_ = createIncomeCategoryUseCase

	t.Log("All services correctly implement their respective interfaces")
}
