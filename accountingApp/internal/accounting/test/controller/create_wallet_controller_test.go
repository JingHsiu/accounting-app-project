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
)

func TestCreateWalletController_CreateWallet_Success(t *testing.T) {
	// Arrange - Use real implementations
	repo, _ := test.NewFakeWalletRepo()
	service := command.NewCreateWalletService(repo)
	ctrl := controller.NewCreateWalletController(service)

	requestBody := map[string]interface{}{
		"user_id":  "test-user",
		"name":     "Test Wallet",
		"type":     "CASH",
		"currency": "USD",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/wallets", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ctrl.CreateWallet(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["success"] != true {
		t.Errorf("Expected success to be true, got %v", response["success"])
	}

	// Verify wallet was actually created in repository
	wallets, _ := repo.FindByUserID("test-user")
	if len(wallets) != 1 {
		t.Errorf("Expected 1 wallet to be created, got %d", len(wallets))
	}
}

func TestCreateWalletController_CreateWallet_WithInitialBalance(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	service := command.NewCreateWalletService(repo)
	ctrl := controller.NewCreateWalletController(service)

	initialBalance := int64(10000) // 100.00 in cents
	requestBody := map[string]interface{}{
		"user_id":        "test-user",
		"name":           "Test Wallet",
		"type":           "CASH",
		"currency":       "USD",
		"initialBalance": initialBalance,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/wallets", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ctrl.CreateWallet(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify wallet has correct initial balance
	wallets, _ := repo.FindByUserID("test-user")
	if len(wallets) != 1 {
		t.Errorf("Expected 1 wallet, got %d", len(wallets))
		return
	}
	
	wallet := wallets[0]
	if wallet.Balance.Amount != initialBalance {
		t.Errorf("Expected balance %d, got %d", initialBalance, wallet.Balance.Amount)
	}
}

func TestCreateWalletController_CreateWallet_MissingUserID(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	service := command.NewCreateWalletService(repo)
	ctrl := controller.NewCreateWalletController(service)

	requestBody := map[string]interface{}{
		"name":     "Test Wallet",
		"type":     "CASH",
		"currency": "USD",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/wallets", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ctrl.CreateWallet(w, req)

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

func TestCreateWalletController_CreateWallet_InvalidWalletType(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	service := command.NewCreateWalletService(repo)
	ctrl := controller.NewCreateWalletController(service)

	requestBody := map[string]interface{}{
		"user_id":  "test-user",
		"name":     "Test Wallet",
		"type":     "INVALID_TYPE",
		"currency": "USD",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/wallets", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ctrl.CreateWallet(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["success"] != false {
		t.Errorf("Expected success to be false, got %v", response["success"])
	}

	// Verify no wallet was created
	wallets, _ := repo.FindByUserID("test-user")
	if len(wallets) != 0 {
		t.Errorf("Expected 0 wallets, got %d", len(wallets))
	}
}

func TestCreateWalletController_CreateWallet_InvalidJSON(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	service := command.NewCreateWalletService(repo)
	ctrl := controller.NewCreateWalletController(service)

	req := httptest.NewRequest("POST", "/api/v1/wallets", bytes.NewBufferString("{invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ctrl.CreateWallet(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateWalletController_CreateWallet_MethodNotAllowed(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	service := command.NewCreateWalletService(repo)
	ctrl := controller.NewCreateWalletController(service)

	req := httptest.NewRequest("GET", "/api/v1/wallets", nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.CreateWallet(w, req)

	// Assert
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}