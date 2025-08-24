package integration

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
	"github.com/JingHsiu/accountingApp/internal/accounting/frameworks/web"
)

// API Response structures for validation
type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type WalletResponse struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Currency    string                 `json:"currency"`
	Balance     map[string]interface{} `json:"balance"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
	IsFullyLoaded bool                 `json:"is_fully_loaded"`
	Transactions  []interface{}        `json:"transactions,omitempty"`
}

type WalletListResponse struct {
	Data  []WalletResponse `json:"data"`
	Count int              `json:"count"`
}

// Mock implementations for full API testing
type FullMockCreateWalletUseCase struct {
	wallets map[string]*model.Wallet
}

func NewFullMockCreateWalletUseCase() *FullMockCreateWalletUseCase {
	return &FullMockCreateWalletUseCase{
		wallets: make(map[string]*model.Wallet),
	}
}

func (m *FullMockCreateWalletUseCase) Execute(input command.CreateWalletInput) common.Output {
	wallet, err := model.NewWallet(input.UserID, input.Name, model.WalletType(input.Type), input.Currency)
	if err != nil {
		return common.UseCaseOutput{
			ID:       "",
			ExitCode: 1,
			Message:  err.Error(),
		}
	}
	
	m.wallets[wallet.ID] = wallet
	
	return common.UseCaseOutput{
		ID:       wallet.ID,
		ExitCode: 0,
		Message:  "Wallet created successfully",
	}
}

type FullMockWalletRepository struct {
	wallets map[string]*model.Wallet
}

func NewFullMockWalletRepository() *FullMockWalletRepository {
	return &FullMockWalletRepository{
		wallets: make(map[string]*model.Wallet),
	}
}

func (m *FullMockWalletRepository) Save(wallet *model.Wallet) error {
	m.wallets[wallet.ID] = wallet
	return nil
}

func (m *FullMockWalletRepository) FindByID(id string) (*model.Wallet, error) {
	return m.wallets[id], nil
}

func (m *FullMockWalletRepository) FindByIDWithTransactions(id string) (*model.Wallet, error) {
	wallet := m.wallets[id]
	if wallet != nil {
		wallet.SetFullyLoaded(true)
	}
	return wallet, nil
}

func (m *FullMockWalletRepository) FindByUserID(userID string) ([]*model.Wallet, error) {
	var wallets []*model.Wallet
	for _, wallet := range m.wallets {
		if wallet.UserID == userID {
			wallets = append(wallets, wallet)
		}
	}
	return wallets, nil
}

func (m *FullMockWalletRepository) Delete(id string) error {
	delete(m.wallets, id)
	return nil
}

func setupTestServer() (*httptest.Server, *FullMockWalletRepository) {
	// Create mock repository
	mockRepo := NewFullMockWalletRepository()
	
	// Create mock use cases
	mockCreateWalletUseCase := NewFullMockCreateWalletUseCase()
	mockGetWalletBalanceUseCase := &MockGetWalletBalanceUseCase{}
	
	// Create specialized controllers
	createWalletController := controller.NewCreateWalletController(mockCreateWalletUseCase)
	queryWalletController := controller.NewQueryWalletController(mockRepo)
	updateWalletController := controller.NewUpdateWalletController(mockRepo)
	deleteWalletController := controller.NewDeleteWalletController(mockRepo)
	
	// Create legacy wallet controller for transaction operations
	walletController := controller.NewWalletController(
		mockCreateWalletUseCase,
		&MockAddExpenseUseCase{},
		&MockAddIncomeUseCase{},
		mockGetWalletBalanceUseCase,
		mockRepo,
	)
	
	// Create category controller with mocks
	mockCreateExpenseCategoryUseCase := &MockCreateExpenseCategoryUseCase{}
	mockCreateIncomeCategoryUseCase := &MockCreateIncomeCategoryUseCase{}
	categoryController := controller.NewCategoryController(
		mockCreateExpenseCategoryUseCase,
		mockCreateIncomeCategoryUseCase,
	)
	
	// Create router with specialized controllers
	router := web.NewRouter(
		createWalletController,
		queryWalletController,
		updateWalletController,
		deleteWalletController,
		walletController, // For transaction and balance operations
		categoryController,
	)
	handler := router.SetupRoutes()
	
	// Create test server
	server := httptest.NewServer(handler)
	
	return server, mockRepo
}

// Additional mock implementations
type MockAddExpenseUseCase struct{}
func (m *MockAddExpenseUseCase) Execute(input command.AddExpenseInput) common.Output {
	return common.UseCaseOutput{ID: "expense-id", ExitCode: 0, Message: "Success"}
}

type MockAddIncomeUseCase struct{}
func (m *MockAddIncomeUseCase) Execute(input command.AddIncomeInput) common.Output {
	return common.UseCaseOutput{ID: "income-id", ExitCode: 0, Message: "Success"}
}

type MockGetWalletBalanceUseCase struct{}
func (m *MockGetWalletBalanceUseCase) Execute(input query.GetWalletBalanceInput) common.Output {
	return common.UseCaseOutput{
		ID:       input.WalletID,
		ExitCode: 0,
		Message:  "Success",
	}
}

// Mock use cases for CategoryController
type MockCreateExpenseCategoryUseCase struct{}
func (m *MockCreateExpenseCategoryUseCase) Execute(input command.CreateExpenseCategoryInput) common.Output {
	return common.UseCaseOutput{ID: "category-id", ExitCode: 0, Message: "Success"}
}

type MockCreateIncomeCategoryUseCase struct{}
func (m *MockCreateIncomeCategoryUseCase) Execute(input command.CreateIncomeCategoryInput) common.Output {
	return common.UseCaseOutput{ID: "income-category-id", ExitCode: 0, Message: "Success"}
}

func TestWalletAPI_CreateWallet(t *testing.T) {
	// Arrange
	server, _ := setupTestServer()
	defer server.Close()

	requestBody := map[string]interface{}{
		"user_id":  "test-user",
		"name":     "Test Wallet",
		"type":     "CASH",
		"currency": "USD",
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	
	// Act
	resp, err := http.Post(server.URL+"/api/v1/wallets", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response ApiResponse
	json.NewDecoder(resp.Body).Decode(&response)
	
	if !response.Success {
		t.Errorf("Expected success to be true, got %v", response.Success)
	}
	
	// Check response structure matches frontend expectations
	// CreateWallet returns {id, success, message} directly
	// Parse the response body manually since it's not wrapped in data
	resp.Body.Close()
	resp2, _ := http.Post(server.URL+"/api/v1/wallets", "application/json", bytes.NewBuffer(jsonBody))
	defer resp2.Body.Close()
	
	var directResponse map[string]interface{}
	json.NewDecoder(resp2.Body).Decode(&directResponse)
	
	if _, exists := directResponse["id"]; !exists {
		t.Error("Expected response to contain 'id' field")
	}
	
	if success, ok := directResponse["success"].(bool); !ok || !success {
		t.Error("Expected success to be true in direct response")
	}
}

func TestWalletAPI_GetWallets(t *testing.T) {
	// Arrange
	server, mockRepo := setupTestServer()
	defer server.Close()
	
	// Add test data
	testWallet, _ := model.NewWallet("test-user", "Test Wallet", model.WalletTypeCash, "USD")
	mockRepo.Save(testWallet)
	
	// Act
	resp, err := http.Get(server.URL + "/api/v1/wallets?userID=test-user")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response ApiResponse
	json.NewDecoder(resp.Body).Decode(&response)
	
	if !response.Success {
		t.Errorf("Expected success to be true, got %v", response.Success)
	}
	
	// Validate response structure
	data, ok := response.Data.(map[string]interface{})
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
	
	// Validate wallet structure
	walletData, ok := wallets[0].(map[string]interface{})
	if !ok {
		t.Error("Expected wallet to be a map")
		return
	}
	
	// Check required fields for frontend
	requiredFields := []string{"id", "user_id", "name", "type", "currency", "balance", "created_at", "updated_at", "is_fully_loaded"}
	for _, field := range requiredFields {
		if _, exists := walletData[field]; !exists {
			t.Errorf("Expected wallet to have field '%s'", field)
		}
	}
}

func TestWalletAPI_GetWallet(t *testing.T) {
	// Arrange
	server, mockRepo := setupTestServer()
	defer server.Close()
	
	// Add test data
	testWallet, _ := model.NewWallet("test-user", "Test Wallet", model.WalletTypeCash, "USD")
	mockRepo.Save(testWallet)
	
	// Act
	resp, err := http.Get(server.URL + "/api/v1/wallets/" + testWallet.ID)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response ApiResponse
	json.NewDecoder(resp.Body).Decode(&response)
	
	if !response.Success {
		t.Errorf("Expected success to be true, got %v", response.Success)
	}
}

func TestWalletAPI_GetWallet_WithTransactions(t *testing.T) {
	// Arrange
	server, mockRepo := setupTestServer()
	defer server.Close()
	
	// Add test data
	testWallet, _ := model.NewWallet("test-user", "Test Wallet", model.WalletTypeCash, "USD")
	mockRepo.Save(testWallet)
	
	// Act - Request wallet with transactions
	resp, err := http.Get(server.URL + "/api/v1/wallets/" + testWallet.ID + "?includeTransactions=true")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response ApiResponse
	json.NewDecoder(resp.Body).Decode(&response)
	
	if !response.Success {
		t.Errorf("Expected success to be true, got %v", response.Success)
	}
	
	// Check that wallet is marked as fully loaded
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Error("Expected data to be a map")
		return
	}
	
	walletData, ok := data["data"].(map[string]interface{})
	if !ok {
		t.Error("Expected wallet data to be a map")
		return
	}
	
	if isFullyLoaded, ok := walletData["is_fully_loaded"].(bool); !ok || !isFullyLoaded {
		t.Error("Expected wallet to be marked as fully loaded when includeTransactions=true")
	}
}

func TestWalletAPI_UpdateWallet(t *testing.T) {
	// Arrange
	server, mockRepo := setupTestServer()
	defer server.Close()
	
	// Add test data
	testWallet, _ := model.NewWallet("test-user", "Test Wallet", model.WalletTypeCash, "USD")
	mockRepo.Save(testWallet)
	
	requestBody := map[string]interface{}{
		"name": "Updated Wallet Name",
		"type": "BANK",
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	
	// Act
	req, _ := http.NewRequest("PUT", server.URL+"/api/v1/wallets/"+testWallet.ID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response ApiResponse
	json.NewDecoder(resp.Body).Decode(&response)
	
	if !response.Success {
		t.Errorf("Expected success to be true, got %v", response.Success)
	}
}

func TestWalletAPI_DeleteWallet(t *testing.T) {
	// Arrange
	server, mockRepo := setupTestServer()
	defer server.Close()
	
	// Add test data
	testWallet, _ := model.NewWallet("test-user", "Test Wallet", model.WalletTypeCash, "USD")
	mockRepo.Save(testWallet)
	
	// Act
	req, _ := http.NewRequest("DELETE", server.URL+"/api/v1/wallets/"+testWallet.ID, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response ApiResponse
	json.NewDecoder(resp.Body).Decode(&response)
	
	if !response.Success {
		t.Errorf("Expected success to be true, got %v", response.Success)
	}
	
	// Verify wallet was deleted
	if wallet, _ := mockRepo.FindByID(testWallet.ID); wallet != nil {
		t.Error("Expected wallet to be deleted")
	}
}

func TestWalletAPI_GetWallets_MissingUserID(t *testing.T) {
	// Arrange
	server, _ := setupTestServer()
	defer server.Close()
	
	// Act
	resp, err := http.Get(server.URL + "/api/v1/wallets") // Missing userID parameter
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	var response ApiResponse
	json.NewDecoder(resp.Body).Decode(&response)
	
	if response.Success {
		t.Errorf("Expected success to be false, got %v", response.Success)
	}
	
	if response.Error == "" {
		t.Error("Expected error message to be provided")
	}
}

func TestWalletAPI_GetWallet_NotFound(t *testing.T) {
	// Arrange
	server, _ := setupTestServer()
	defer server.Close()
	
	// Act
	resp, err := http.Get(server.URL + "/api/v1/wallets/non-existent-id")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}

	var response ApiResponse
	json.NewDecoder(resp.Body).Decode(&response)
	
	if response.Success {
		t.Errorf("Expected success to be false, got %v", response.Success)
	}
}