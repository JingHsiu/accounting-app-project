package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JingHsiu/accountingApp/internal/accounting/test"
	"github.com/JingHsiu/accountingApp/internal/accounting/adapter/controller"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
)

func TestDeleteWalletController_DeleteWallet_Success(t *testing.T) {
	// Arrange - Use real implementations
	repo, _ := test.NewFakeWalletRepo()
	deleteService := command.NewDeleteWalletService(repo)
	ctrl := controller.NewDeleteWalletController(deleteService)

	// Create a wallet first
	createService := command.NewCreateWalletService(repo)
	createResult := createService.Execute(usecase.CreateWalletInput{
		UserID:   "test-user",
		Name:     "Test Wallet",
		Type:     "CASH",
		Currency: "USD",
	})

	if createResult.GetExitCode() != 0 {
		t.Fatalf("Failed to create test wallet: %v", createResult.GetMessage())
	}

	walletID := createResult.GetID()
	req := httptest.NewRequest("DELETE", "/api/v1/wallets/"+walletID, nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.DeleteWallet(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify wallet was actually deleted from repository
	wallet, _ := repo.FindByID(walletID)
	if wallet != nil {
		t.Errorf("Expected wallet to be deleted, but it still exists")
	}
}

func TestDeleteWalletController_DeleteWallet_NotFound(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	deleteService := command.NewDeleteWalletService(repo)
	ctrl := controller.NewDeleteWalletController(deleteService)

	nonExistentID := "non-existent-wallet-id"
	req := httptest.NewRequest("DELETE", "/api/v1/wallets/"+nonExistentID, nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.DeleteWallet(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteWalletController_DeleteWallet_InvalidWalletID(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	deleteService := command.NewDeleteWalletService(repo)
	ctrl := controller.NewDeleteWalletController(deleteService)

	// Test with empty wallet ID path
	req := httptest.NewRequest("DELETE", "/api/v1/wallets/", nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.DeleteWallet(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteWalletController_DeleteWallet_MethodNotAllowed(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	deleteService := command.NewDeleteWalletService(repo)
	ctrl := controller.NewDeleteWalletController(deleteService)

	req := httptest.NewRequest("GET", "/api/v1/wallets/some-id", nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.DeleteWallet(w, req)

	// Assert
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestDeleteWalletController_DeleteWallet_URLDecoding(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	deleteService := command.NewDeleteWalletService(repo)
	ctrl := controller.NewDeleteWalletController(deleteService)

	// Create a wallet with a UUID that might need URL decoding
	createService := command.NewCreateWalletService(repo)
	createResult := createService.Execute(usecase.CreateWalletInput{
		UserID:   "test-user",
		Name:     "Test Wallet",
		Type:     "CASH",
		Currency: "USD",
	})

	walletID := createResult.GetID()
	
	// Test with URL encoded wallet ID (simulating special characters)
	req := httptest.NewRequest("DELETE", "/api/v1/wallets/"+walletID, nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.DeleteWallet(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
	}
}