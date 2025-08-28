package usecase

import (
	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
	"testing"

	"github.com/JingHsiu/accountingApp/internal/accounting/test"
	"github.com/stretchr/testify/assert"
)

func Test_CreateWallet_Success(t *testing.T) {
	repo, _ := test.NewFakeWalletRepo()
	service := command.NewCreateWalletService(repo)
	input := usecase.CreateWalletInput{
		UserID:   "user-123",
		Name:     "My Wallet",
		Type:     "CASH",
		Currency: "USD",
	}

	output := service.Execute(input)
	assert.Equal(t, common.Success, output.GetExitCode())

	saved, err := repo.FindByID(output.GetID())
	assert.NoError(t, err)
	assert.NotNil(t, saved)
	assert.Equal(t, "My Wallet", saved.Name)

}
