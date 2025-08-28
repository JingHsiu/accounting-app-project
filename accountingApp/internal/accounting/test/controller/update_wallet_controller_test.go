package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JingHsiu/accountingApp/internal/accounting/test"
	"github.com/JingHsiu/accountingApp/internal/accounting/adapter/controller"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

func TestUpdateWalletController_UpdateWallet_Success_Name(t *testing.T) {
	// Arrange - Use real implementations
	repo, _ := test.NewFakeWalletRepo()
	updateService := command.NewUpdateWalletService(repo)
	ctrl := controller.NewUpdateWalletController(updateService)

	// Create a wallet first
	createService := command.NewCreateWalletService(repo)
	createResult := createService.Execute(usecase.CreateWalletInput{
		UserID:   "test-user",
		Name:     "Original Wallet",
		Type:     "CASH",
		Currency: "USD",
	})

	walletID := createResult.GetID()

	requestBody := map[string]interface{}{
		"name": "Updated Wallet Name",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/api/v1/wallets/"+walletID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ctrl.UpdateWallet(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["success"] != true {
		t.Errorf("Expected success to be true, got %v", response["success"])
	}

	// Verify wallet was actually updated
	wallet, _ := repo.FindByID(walletID)
	if wallet.Name != "Updated Wallet Name" {
		t.Errorf("Expected wallet name to be 'Updated Wallet Name', got %s", wallet.Name)
	}
}

func TestUpdateWalletController_UpdateWallet_Success_Type(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	updateService := command.NewUpdateWalletService(repo)
	ctrl := controller.NewUpdateWalletController(updateService)

	// Create a wallet first
	createService := command.NewCreateWalletService(repo)
	createResult := createService.Execute(usecase.CreateWalletInput{
		UserID:   "test-user",
		Name:     "Test Wallet",
		Type:     "CASH",
		Currency: "USD",
	})

	walletID := createResult.GetID()

	requestBody := map[string]interface{}{
		"type": "BANK",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/api/v1/wallets/"+walletID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ctrl.UpdateWallet(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify wallet type was updated
	wallet, _ := repo.FindByID(walletID)
	if wallet.Type != model.WalletTypeBank {
		t.Errorf("Expected wallet type to be BANK_ACCOUNT, got %v", wallet.Type)
	}
}

func TestUpdateWalletController_UpdateWallet_Success_Multiple_Fields(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	updateService := command.NewUpdateWalletService(repo)
	ctrl := controller.NewUpdateWalletController(updateService)

	// Create a wallet first
	createService := command.NewCreateWalletService(repo)
	createResult := createService.Execute(usecase.CreateWalletInput{
		UserID:   "test-user",
		Name:     "Original Wallet",
		Type:     "CASH",
		Currency: "USD",
	})

	walletID := createResult.GetID()

	requestBody := map[string]interface{}{
		"name": "Updated Wallet",
		"type": "CREDIT",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/api/v1/wallets/"+walletID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ctrl.UpdateWallet(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify both fields were updated
	wallet, _ := repo.FindByID(walletID)
	if wallet.Name != "Updated Wallet" {
		t.Errorf("Expected wallet name to be 'Updated Wallet', got %s", wallet.Name)
	}
	if wallet.Type != model.WalletTypeCredit {
		t.Errorf("Expected wallet type to be CREDIT_CARD, got %v", wallet.Type)
	}
}

func TestUpdateWalletController_UpdateWallet_NotFound(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	updateService := command.NewUpdateWalletService(repo)
	ctrl := controller.NewUpdateWalletController(updateService)

	nonExistentID := "non-existent-wallet-id"
	requestBody := map[string]interface{}{
		"name": "New Name",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/api/v1/wallets/"+nonExistentID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ctrl.UpdateWallet(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateWalletController_UpdateWallet_InvalidWalletType(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	updateService := command.NewUpdateWalletService(repo)
	ctrl := controller.NewUpdateWalletController(updateService)

	// Create a wallet first
	createService := command.NewCreateWalletService(repo)
	createResult := createService.Execute(usecase.CreateWalletInput{
		UserID:   "test-user",
		Name:     "Test Wallet",
		Type:     "CASH",
		Currency: "USD",
	})

	walletID := createResult.GetID()

	requestBody := map[string]interface{}{
		"type": "INVALID_TYPE",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/api/v1/wallets/"+walletID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ctrl.UpdateWallet(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["success"] != false {
		t.Errorf("Expected success to be false, got %v", response["success"])
	}
}

func TestUpdateWalletController_UpdateWallet_EmptyName(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	updateService := command.NewUpdateWalletService(repo)
	ctrl := controller.NewUpdateWalletController(updateService)

	// Create a wallet first
	createService := command.NewCreateWalletService(repo)
	createResult := createService.Execute(usecase.CreateWalletInput{
		UserID:   "test-user",
		Name:     "Original Name",
		Type:     "CASH",
		Currency: "USD",
	})

	walletID := createResult.GetID()

	requestBody := map[string]interface{}{
		"name": "",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/api/v1/wallets/"+walletID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ctrl.UpdateWallet(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateWalletController_UpdateWallet_InvalidWalletID(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	updateService := command.NewUpdateWalletService(repo)
	ctrl := controller.NewUpdateWalletController(updateService)

	requestBody := map[string]interface{}{
		"name": "New Name",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/api/v1/wallets/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ctrl.UpdateWallet(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateWalletController_UpdateWallet_InvalidJSON(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	updateService := command.NewUpdateWalletService(repo)
	ctrl := controller.NewUpdateWalletController(updateService)

	// Create a wallet first
	createService := command.NewCreateWalletService(repo)
	createResult := createService.Execute(usecase.CreateWalletInput{
		UserID:   "test-user",
		Name:     "Test Wallet",
		Type:     "CASH",
		Currency: "USD",
	})

	walletID := createResult.GetID()

	req := httptest.NewRequest("PUT", "/api/v1/wallets/"+walletID, bytes.NewBufferString("{invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ctrl.UpdateWallet(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateWalletController_UpdateWallet_MethodNotAllowed(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	updateService := command.NewUpdateWalletService(repo)
	ctrl := controller.NewUpdateWalletController(updateService)

	req := httptest.NewRequest("GET", "/api/v1/wallets/some-id", nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.UpdateWallet(w, req)

	// Assert
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestUpdateWalletController_UpdateWallet_NoChanges(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	updateService := command.NewUpdateWalletService(repo)
	ctrl := controller.NewUpdateWalletController(updateService)

	// Create a wallet first
	createService := command.NewCreateWalletService(repo)
	createResult := createService.Execute(usecase.CreateWalletInput{
		UserID:   "test-user",
		Name:     "Test Wallet",
		Type:     "CASH",
		Currency: "USD",
	})

	walletID := createResult.GetID()

	// Send empty update (no fields to update)
	requestBody := map[string]interface{}{}

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/api/v1/wallets/"+walletID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ctrl.UpdateWallet(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify wallet remains unchanged
	wallet, _ := repo.FindByID(walletID)
	if wallet.Name != "Test Wallet" {
		t.Errorf("Expected wallet name to remain 'Test Wallet', got %s", wallet.Name)
	}
}