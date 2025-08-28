package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JingHsiu/accountingApp/internal/accounting/test"
	"github.com/JingHsiu/accountingApp/internal/accounting/adapter/controller"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/query"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
)

func TestGetWalletController_GetWallets_Success(t *testing.T) {
	// Arrange - Use real implementations
	repo, _ := test.NewFakeWalletRepo()
	getWalletsService := query.NewGetWalletsService(repo)
	getWalletService := query.NewGetWalletService(repo)
	ctrl := controller.NewQueryWalletController(getWalletsService, getWalletService)

	// Create test wallets
	createService := command.NewCreateWalletService(repo)
	createService.Execute(usecase.CreateWalletInput{
		UserID:   "test-user",
		Name:     "Test Wallet 1",
		Type:     "CASH",
		Currency: "USD",
	})
	createService.Execute(usecase.CreateWalletInput{
		UserID:   "test-user",
		Name:     "Test Wallet 2", 
		Type:     "BANK",
		Currency: "EUR",
	})

	req := httptest.NewRequest("GET", "/api/v1/wallets?userID=test-user", nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.GetWallets(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["success"] != true {
		t.Errorf("Expected success to be true, got %v", response["success"])
	}

	data := response["data"].(map[string]interface{})
	wallets := data["data"].([]interface{})
	count := data["count"].(float64)

	if count != 2 {
		t.Errorf("Expected count to be 2, got %v", count)
	}

	if len(wallets) != 2 {
		t.Errorf("Expected 2 wallets, got %d", len(wallets))
	}
}

func TestGetWalletController_GetWallets_EmptyResult(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	getWalletsService := query.NewGetWalletsService(repo)
	getWalletService := query.NewGetWalletService(repo)
	ctrl := controller.NewQueryWalletController(getWalletsService, getWalletService)

	req := httptest.NewRequest("GET", "/api/v1/wallets?userID=non-existent-user", nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.GetWallets(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	count := data["count"].(float64)

	if count != 0 {
		t.Errorf("Expected count to be 0, got %v", count)
	}
}

func TestGetWalletController_GetWallets_MissingUserID(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	getWalletsService := query.NewGetWalletsService(repo)
	getWalletService := query.NewGetWalletService(repo)
	ctrl := controller.NewQueryWalletController(getWalletsService, getWalletService)

	req := httptest.NewRequest("GET", "/api/v1/wallets", nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.GetWallets(w, req)

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

func TestGetWalletController_GetWallets_MethodNotAllowed(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	getWalletsService := query.NewGetWalletsService(repo)
	getWalletService := query.NewGetWalletService(repo)
	ctrl := controller.NewQueryWalletController(getWalletsService, getWalletService)

	req := httptest.NewRequest("POST", "/api/v1/wallets?userID=test-user", nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.GetWallets(w, req)

	// Assert
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestGetWalletController_GetWallet_Success(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	getWalletsService := query.NewGetWalletsService(repo)
	getWalletService := query.NewGetWalletService(repo)
	ctrl := controller.NewQueryWalletController(getWalletsService, getWalletService)

	// Create test wallet
	createService := command.NewCreateWalletService(repo)
	createResult := createService.Execute(usecase.CreateWalletInput{
		UserID:   "test-user",
		Name:     "Test Wallet",
		Type:     "CASH",
		Currency: "USD",
	})

	walletID := createResult.GetID()
	req := httptest.NewRequest("GET", "/api/v1/wallets/"+walletID, nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.GetWallet(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["success"] != true {
		t.Errorf("Expected success to be true, got %v", response["success"])
	}

	data := response["data"].(map[string]interface{})
	walletData := data["data"].(map[string]interface{})
	
	if walletData["name"] != "Test Wallet" {
		t.Errorf("Expected wallet name 'Test Wallet', got %v", walletData["name"])
	}
}

func TestGetWalletController_GetWallet_WithTransactions(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	getWalletsService := query.NewGetWalletsService(repo)
	getWalletService := query.NewGetWalletService(repo)
	ctrl := controller.NewQueryWalletController(getWalletsService, getWalletService)

	// Create test wallet
	createService := command.NewCreateWalletService(repo)
	createResult := createService.Execute(usecase.CreateWalletInput{
		UserID:   "test-user",
		Name:     "Test Wallet",
		Type:     "CASH",
		Currency: "USD",
	})

	walletID := createResult.GetID()
	req := httptest.NewRequest("GET", "/api/v1/wallets/"+walletID+"?includeTransactions=true", nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.GetWallet(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	walletData := data["data"].(map[string]interface{})
	
	// Check if transactions field is present (even if empty)
	if _, hasTransactions := walletData["transactions"]; !hasTransactions {
		t.Errorf("Expected transactions field to be present when includeTransactions=true")
	}
}

func TestGetWalletController_GetWallet_NotFound(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	getWalletsService := query.NewGetWalletsService(repo)
	getWalletService := query.NewGetWalletService(repo)
	ctrl := controller.NewQueryWalletController(getWalletsService, getWalletService)

	nonExistentID := "non-existent-wallet-id"
	req := httptest.NewRequest("GET", "/api/v1/wallets/"+nonExistentID, nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.GetWallet(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["success"] != false {
		t.Errorf("Expected success to be false, got %v", response["success"])
	}
}

func TestGetWalletController_GetWallet_InvalidWalletID(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	getWalletsService := query.NewGetWalletsService(repo)
	getWalletService := query.NewGetWalletService(repo)
	ctrl := controller.NewQueryWalletController(getWalletsService, getWalletService)

	req := httptest.NewRequest("GET", "/api/v1/wallets/", nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.GetWallet(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetWalletController_GetWallet_MethodNotAllowed(t *testing.T) {
	// Arrange
	repo, _ := test.NewFakeWalletRepo()
	getWalletsService := query.NewGetWalletsService(repo)
	getWalletService := query.NewGetWalletService(repo)
	ctrl := controller.NewQueryWalletController(getWalletsService, getWalletService)

	req := httptest.NewRequest("DELETE", "/api/v1/wallets/some-id", nil)
	w := httptest.NewRecorder()

	// Act
	ctrl.GetWallet(w, req)

	// Assert
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}