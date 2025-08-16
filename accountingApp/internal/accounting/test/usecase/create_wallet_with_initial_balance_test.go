package usecase

import (
	"testing"

	"github.com/JingHsiu/accountingApp/internal/accounting/adapter"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/stretchr/testify/assert"
)

func TestCreateWalletWithInitialBalance(t *testing.T) {
	// Arrange
	repo, _ := adapter.NewFakeWalletRepo()
	service := command.NewCreateWalletService(repo)

	testCases := []struct {
		name           string
		input          command.CreateWalletInput
		expectedAmount int64
		shouldSucceed  bool
	}{
		{
			name: "Create wallet with zero initial balance",
			input: command.CreateWalletInput{
				UserID:         "user123",
				Name:           "Test Wallet",
				Type:           "CASH",
				Currency:       "TWD",
				InitialBalance: nil, // Should default to 0
			},
			expectedAmount: 0,
			shouldSucceed:  true,
		},
		{
			name: "Create wallet with positive initial balance",
			input: command.CreateWalletInput{
				UserID:         "user123",
				Name:           "Test Wallet",
				Type:           "BANK",
				Currency:       "USD",
				InitialBalance: int64Ptr(150000), // $1500.00
			},
			expectedAmount: 150000,
			shouldSucceed:  true,
		},
		{
			name: "Create wallet with negative initial balance should fail",
			input: command.CreateWalletInput{
				UserID:         "user123",
				Name:           "Test Wallet",
				Type:           "CASH",
				Currency:       "TWD",
				InitialBalance: int64Ptr(-10000),
			},
			expectedAmount: 0,
			shouldSucceed:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := service.Execute(tc.input)

			// Assert
			if tc.shouldSucceed {
				assert.Equal(t, common.Success, result.GetExitCode(), "Expected success but got: %s", result.GetMessage())

				// Verify wallet was created with correct initial balance
				walletID := result.GetID()
				assert.NotEmpty(t, walletID, "Expected wallet ID to be returned")

				// Retrieve the created wallet
				wallet, err := repo.FindByID(walletID)
				assert.NoError(t, err, "Failed to retrieve created wallet")
				assert.NotNil(t, wallet, "Created wallet not found")

				assert.Equal(t, tc.expectedAmount, wallet.Balance.Amount, "Initial balance mismatch")
				assert.Equal(t, tc.input.Currency, wallet.Balance.Currency, "Currency mismatch")
			} else {
				assert.NotEqual(t, common.Success, result.GetExitCode(), "Expected failure, but got success")
			}
		})
	}
}

func TestCreateWalletWithDifferentTypes(t *testing.T) {
	// Test all wallet types
	repo, _ := adapter.NewFakeWalletRepo()
	service := command.NewCreateWalletService(repo)

	walletTypes := []string{"CASH", "BANK", "CREDIT", "INVESTMENT"}
	initialBalance := int64(100000) // $1000.00

	for _, walletType := range walletTypes {
		t.Run("Create_"+walletType+"_wallet", func(t *testing.T) {
			input := command.CreateWalletInput{
				UserID:         "user123",
				Name:           walletType + " Wallet",
				Type:           walletType,
				Currency:       "USD",
				InitialBalance: &initialBalance,
			}

			result := service.Execute(input)

			assert.Equal(t, common.Success, result.GetExitCode(), "Failed to create %s wallet: %s", walletType, result.GetMessage())

			// Verify wallet was created correctly
			wallet, err := repo.FindByID(result.GetID())
			assert.NoError(t, err, "Failed to retrieve %s wallet", walletType)
			assert.NotNil(t, wallet, "Wallet not found")

			assert.Equal(t, walletType, string(wallet.Type), "Wallet type mismatch")
			assert.Equal(t, initialBalance, wallet.Balance.Amount, "Initial balance mismatch")
		})
	}
}

// Helper function to create pointer to int64
func int64Ptr(val int64) *int64 {
	return &val
}