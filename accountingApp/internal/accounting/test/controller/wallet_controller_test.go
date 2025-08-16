package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JingHsiu/accountingApp/internal/accounting/adapter/controller"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/query"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

// Mock implementations for testing

type MockCreateWalletUseCase struct {
	output common.Output
}

func (m *MockCreateWalletUseCase) Execute(input command.CreateWalletInput) common.Output {
	return m.output
}

type MockAddExpenseUseCase struct {
	output common.Output
}

func (m *MockAddExpenseUseCase) Execute(input command.AddExpenseInput) common.Output {
	return m.output
}

type MockAddIncomeUseCase struct {
	output common.Output
}

func (m *MockAddIncomeUseCase) Execute(input command.AddIncomeInput) common.Output {
	return m.output
}

type MockGetWalletBalanceUseCase struct {
	output common.Output
}

func (m *MockGetWalletBalanceUseCase) Execute(input query.GetWalletBalanceInput) common.Output {
	return m.output
}

type MockWalletRepository struct {
	wallets          map[string]*model.Wallet
	findByUserIDFunc func(userID string) ([]*model.Wallet, error)
	saveFunc         func(wallet *model.Wallet) error
	deleteFunc       func(id string) error
}

func NewMockWalletRepository() *MockWalletRepository {
	return &MockWalletRepository{
		wallets: make(map[string]*model.Wallet),
	}
}

func (m *MockWalletRepository) Save(wallet *model.Wallet) error {
	if m.saveFunc != nil {
		return m.saveFunc(wallet)
	}
	m.wallets[wallet.ID] = wallet
	return nil
}

func (m *MockWalletRepository) FindByID(id string) (*model.Wallet, error) {
	return m.wallets[id], nil
}

func (m *MockWalletRepository) FindByIDWithTransactions(id string) (*model.Wallet, error) {
	wallet := m.wallets[id]
	if wallet != nil {
		wallet.SetFullyLoaded(true)
	}
	return wallet, nil
}

func (m *MockWalletRepository) FindByUserID(userID string) ([]*model.Wallet, error) {
	if m.findByUserIDFunc != nil {
		return m.findByUserIDFunc(userID)
	}
	
	var wallets []*model.Wallet
	for _, wallet := range m.wallets {
		if wallet.UserID == userID {
			wallets = append(wallets, wallet)
		}
	}
	return wallets, nil
}

func (m *MockWalletRepository) Delete(id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	delete(m.wallets, id)
	return nil
}

func TestWalletController_CreateWallet(t *testing.T) {
	// Arrange
	mockCreateWalletUseCase := &MockCreateWalletUseCase{
		output: common.UseCaseOutput{
			ID:       "test-wallet-id",
			ExitCode: 0,
			Message:  "Wallet created successfully",
		},
	}
	
	mockRepo := NewMockWalletRepository()
	
	controller := controller.NewWalletController(
		mockCreateWalletUseCase,
		&MockAddExpenseUseCase{},
		&MockAddIncomeUseCase{},
		&MockGetWalletBalanceUseCase{},
		mockRepo,
	)

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
	controller.CreateWallet(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["success"] != true {
		t.Errorf("Expected success to be true, got %v", response["success"])
	}
	
	if response["id"] != "test-wallet-id" {
		t.Errorf("Expected id to be 'test-wallet-id', got %v", response["id"])
	}
}

func TestWalletController_GetWallets(t *testing.T) {
	// Arrange
	testWallet, _ := model.NewWallet("test-user", "Test Wallet", model.WalletTypeCash, "USD")
	
	mockRepo := NewMockWalletRepository()
	mockRepo.wallets["wallet-1"] = testWallet
	mockRepo.findByUserIDFunc = func(userID string) ([]*model.Wallet, error) {
		if userID == "test-user" {
			return []*model.Wallet{testWallet}, nil
		}
		return []*model.Wallet{}, nil
	}
	
	controller := controller.NewWalletController(
		&MockCreateWalletUseCase{},
		&MockAddExpenseUseCase{},
		&MockAddIncomeUseCase{},
		&MockGetWalletBalanceUseCase{},
		mockRepo,
	)

	req := httptest.NewRequest("GET", "/api/v1/wallets?userID=test-user", nil)
	w := httptest.NewRecorder()

	// Act
	controller.GetWallets(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["success"] != true {
		t.Errorf("Expected success to be true, got %v", response["success"])
	}
	
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Error("Expected data to be a map")
		return
	}
	
	wallets, ok := data["data"].([]interface{})
	if !ok {
		t.Error("Expected data.data to be an array")
		return
	}
	
	if len(wallets) != 1 {
		t.Errorf("Expected 1 wallet, got %d", len(wallets))
	}
}

func TestWalletController_GetWallet(t *testing.T) {
	// Arrange
	testWallet, _ := model.NewWallet("test-user", "Test Wallet", model.WalletTypeCash, "USD")
	
	mockRepo := NewMockWalletRepository()
	mockRepo.wallets["test-wallet-id"] = testWallet
	
	controller := controller.NewWalletController(
		&MockCreateWalletUseCase{},
		&MockAddExpenseUseCase{},
		&MockAddIncomeUseCase{},
		&MockGetWalletBalanceUseCase{},
		mockRepo,
	)

	req := httptest.NewRequest("GET", "/api/v1/wallets/test-wallet-id", nil)
	w := httptest.NewRecorder()

	// Act
	controller.GetWallet(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["success"] != true {
		t.Errorf("Expected success to be true, got %v", response["success"])
	}
	
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Error("Expected data to be a map")
		return
	}
	
	walletData, ok := data["data"].(map[string]interface{})
	if !ok {
		t.Error("Expected wallet data to be a map")
		return
	}
	
	if walletData["id"] != testWallet.ID {
		t.Errorf("Expected wallet ID %s, got %v", testWallet.ID, walletData["id"])
	}
}

func TestWalletController_GetWallet_NotFound(t *testing.T) {
	// Arrange
	mockRepo := NewMockWalletRepository()
	
	controller := controller.NewWalletController(
		&MockCreateWalletUseCase{},
		&MockAddExpenseUseCase{},
		&MockAddIncomeUseCase{},
		&MockGetWalletBalanceUseCase{},
		mockRepo,
	)

	req := httptest.NewRequest("GET", "/api/v1/wallets/non-existent-id", nil)
	w := httptest.NewRecorder()

	// Act
	controller.GetWallet(w, req)

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

func TestWalletController_DeleteWallet(t *testing.T) {
	// Arrange
	testWallet, _ := model.NewWallet("test-user", "Test Wallet", model.WalletTypeCash, "USD")
	
	mockRepo := NewMockWalletRepository()
	mockRepo.wallets["test-wallet-id"] = testWallet
	
	controller := controller.NewWalletController(
		&MockCreateWalletUseCase{},
		&MockAddExpenseUseCase{},
		&MockAddIncomeUseCase{},
		&MockGetWalletBalanceUseCase{},
		mockRepo,
	)

	req := httptest.NewRequest("DELETE", "/api/v1/wallets/test-wallet-id", nil)
	w := httptest.NewRecorder()

	// Act
	controller.DeleteWallet(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["success"] != true {
		t.Errorf("Expected success to be true, got %v", response["success"])
	}
	
	// Verify wallet was deleted
	if _, exists := mockRepo.wallets["test-wallet-id"]; exists {
		t.Error("Expected wallet to be deleted")
	}
}

func TestWalletController_UpdateWallet(t *testing.T) {
	// Arrange
	testWallet, _ := model.NewWallet("test-user", "Test Wallet", model.WalletTypeCash, "USD")
	
	mockRepo := NewMockWalletRepository()
	mockRepo.wallets["test-wallet-id"] = testWallet
	
	controller := controller.NewWalletController(
		&MockCreateWalletUseCase{},
		&MockAddExpenseUseCase{},
		&MockAddIncomeUseCase{},
		&MockGetWalletBalanceUseCase{},
		mockRepo,
	)

	requestBody := map[string]interface{}{
		"name": "Updated Wallet Name",
		"type": "BANK",
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/api/v1/wallets/test-wallet-id", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.UpdateWallet(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["success"] != true {
		t.Errorf("Expected success to be true, got %v", response["success"])
	}
}

func TestWalletController_GetWallets_MissingUserID(t *testing.T) {
	// Arrange
	mockRepo := NewMockWalletRepository()
	
	controller := controller.NewWalletController(
		&MockCreateWalletUseCase{},
		&MockAddExpenseUseCase{},
		&MockAddIncomeUseCase{},
		&MockGetWalletBalanceUseCase{},
		mockRepo,
	)

	req := httptest.NewRequest("GET", "/api/v1/wallets", nil) // Missing userID parameter
	w := httptest.NewRecorder()

	// Act
	controller.GetWallets(w, req)

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