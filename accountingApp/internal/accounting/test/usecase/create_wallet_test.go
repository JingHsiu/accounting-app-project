package usecase

import (
	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"testing"

	"github.com/JingHsiu/accountingApp/internal/accounting/adapter"
	"github.com/stretchr/testify/assert"
)

func Test_CreateWallet_Success(t *testing.T) {
	repo, _ := adapter.NewFakeWalletRepo()
	usecase := command.NewCreateWalletService(repo)
	input := command.CreateWalletInput{
		UserID:   "user-123",
		Name:     "My Wallet",
		Type:     "CASH",
		Currency: "USD",
	}

	output := usecase.Execute(input)
	assert.Equal(t, common.Success, output.GetExitCode())

	saved, err := repo.FindByID(output.GetID())
	assert.NoError(t, err)
	assert.NotNil(t, saved)
	assert.Equal(t, "My Wallet", saved.Name)

}
