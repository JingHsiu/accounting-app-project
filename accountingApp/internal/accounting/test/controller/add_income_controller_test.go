package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JingHsiu/accountingApp/internal/accounting/test"
	"github.com/JingHsiu/accountingApp/internal/accounting/adapter/controller"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

// setupAddIncomeController creates a controller with repositories and test data
func setupAddIncomeController(t *testing.T) (*controller.AddIncomeController, string, string) {
	// Setup repositories
	walletRepo, err := test.NewFakeWalletRepo()
	if err != nil {
		t.Fatalf("Failed to create wallet repo: %v", err)
	}
	
	categoryRepo := test.NewFakeIncomeCategoryRepository()
	
	// Create use case service
	addIncomeService := command.NewAddIncomeService(walletRepo, categoryRepo)
	ctrl := controller.NewAddIncomeController(addIncomeService)
	
	// Create test wallet
	createWalletService := command.NewCreateWalletService(walletRepo)
	walletResult := createWalletService.Execute(usecase.CreateWalletInput{
		UserID:   "test-user",
		Name:     "Test Wallet",
		Type:     "CASH",
		Currency: "USD",
	})
	walletID := walletResult.GetID()
	
	// Create test income category with subcategory
	categoryName, _ := model.NewCategoryName("Test Category")
	category, _ := model.NewIncomeCategory("test-user", *categoryName)
	
	// Add subcategory to the category
	subcategoryName, _ := model.NewCategoryName("Test Subcategory")
	subcategory, _ := category.AddSubcategory(*subcategoryName)
	subcategoryID := subcategory.ID
	
	// Save category to repository
	categoryRepo.Save(category)
	
	return ctrl, walletID, subcategoryID
}

// TestAddIncomeController_AddIncome_Success tests successful income addition
func TestAddIncomeController_AddIncome_Success(t *testing.T) {
	// Arrange
	ctrl, walletID, subcategoryID := setupAddIncomeController(t)
	
	requestBody := map[string]interface{}{
		"wallet_id":      walletID,
		"subcategory_id": subcategoryID,
		"amount":         10000, // $100.00 in cents
		"currency":       "USD",
		"description":    "Test income",
		"date":          time.Now().Format(time.RFC3339),
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	// Act
	ctrl.AddIncome(w, req)
	
	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
	}
	
	// Verify response structure
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	// Check response fields
	if response["success"] != true {
		t.Errorf("Expected success to be true, got %v", response["success"])
	}
	
	if response["id"] == nil || response["id"] == "" {
		t.Errorf("Expected non-empty id, got %v", response["id"])
	}
	
	if response["message"] != "" {
		t.Errorf("Expected empty message on success, got %v", response["message"])
	}
	
	// Verify Content-Type header
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}
}

// TestAddIncomeController_AddIncome_Success_WithDescription tests successful income with description
func TestAddIncomeController_AddIncome_Success_WithDescription(t *testing.T) {
	// Arrange
	ctrl, walletID, subcategoryID := setupAddIncomeController(t)
	
	requestBody := map[string]interface{}{
		"wallet_id":      walletID,
		"subcategory_id": subcategoryID,
		"amount":         50000, // $500.00 in cents
		"currency":       "USD",
		"description":    "Monthly salary bonus",
		"date":          time.Now().Format(time.RFC3339),
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	// Act
	ctrl.AddIncome(w, req)
	
	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["success"] != true {
		t.Errorf("Expected success to be true, got %v", response["success"])
	}
}

// TestAddIncomeController_AddIncome_MethodNotAllowed tests invalid HTTP method
func TestAddIncomeController_AddIncome_MethodNotAllowed(t *testing.T) {
	// Arrange
	ctrl, _, _ := setupAddIncomeController(t)
	
	// Test various invalid methods
	methods := []string{"GET", "PUT", "DELETE", "PATCH"}
	
	for _, method := range methods {
		t.Run("Method_"+method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/v1/incomes", nil)
			w := httptest.NewRecorder()
			
			// Act
			ctrl.AddIncome(w, req)
			
			// Assert
			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status %d for method %s, got %d", http.StatusMethodNotAllowed, method, w.Code)
			}
		})
	}
}

// TestAddIncomeController_AddIncome_InvalidJSON tests invalid JSON request
func TestAddIncomeController_AddIncome_InvalidJSON(t *testing.T) {
	// Arrange
	ctrl, _, _ := setupAddIncomeController(t)
	
	// Test malformed JSON
	req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBufferString("{invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	// Act
	ctrl.AddIncome(w, req)
	
	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["success"] != false {
		t.Errorf("Expected success to be false, got %v", response["success"])
	}
	
	if response["error"] != "Invalid JSON" {
		t.Errorf("Expected error message 'Invalid JSON', got %v", response["error"])
	}
}

// TestAddIncomeController_AddIncome_MissingWalletID tests missing wallet_id
func TestAddIncomeController_AddIncome_MissingWalletID(t *testing.T) {
	// Arrange
	ctrl, _, subcategoryID := setupAddIncomeController(t)
	
	requestBody := map[string]interface{}{
		// "wallet_id" intentionally omitted
		"subcategory_id": subcategoryID,
		"amount":         10000,
		"currency":       "USD",
		"description":    "Test income",
		"date":          time.Now().Format(time.RFC3339),
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	// Act
	ctrl.AddIncome(w, req)
	
	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["success"] != false {
		t.Errorf("Expected success to be false, got %v", response["success"])
	}
	
	if response["error"] != "wallet_id is required" {
		t.Errorf("Expected error message 'wallet_id is required', got %v", response["error"])
	}
}

// TestAddIncomeController_AddIncome_EmptyWalletID tests empty wallet_id
func TestAddIncomeController_AddIncome_EmptyWalletID(t *testing.T) {
	// Arrange
	ctrl, _, subcategoryID := setupAddIncomeController(t)
	
	requestBody := map[string]interface{}{
		"wallet_id":      "", // Empty string
		"subcategory_id": subcategoryID,
		"amount":         10000,
		"currency":       "USD",
		"description":    "Test income",
		"date":          time.Now().Format(time.RFC3339),
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	// Act
	ctrl.AddIncome(w, req)
	
	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["success"] != false {
		t.Errorf("Expected success to be false, got %v", response["success"])
	}
	
	if response["error"] != "wallet_id is required" {
		t.Errorf("Expected error message 'wallet_id is required', got %v", response["error"])
	}
}

// TestAddIncomeController_AddIncome_MissingSubcategoryID tests missing subcategory_id
func TestAddIncomeController_AddIncome_MissingSubcategoryID(t *testing.T) {
	// Arrange
	ctrl, walletID, _ := setupAddIncomeController(t)
	
	requestBody := map[string]interface{}{
		"wallet_id": walletID,
		// "subcategory_id" intentionally omitted
		"amount":      10000,
		"currency":    "USD",
		"description": "Test income",
		"date":       time.Now().Format(time.RFC3339),
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	// Act
	ctrl.AddIncome(w, req)
	
	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["error"] != "subcategory_id is required" {
		t.Errorf("Expected error message 'subcategory_id is required', got %v", response["error"])
	}
}

// TestAddIncomeController_AddIncome_EmptySubcategoryID tests empty subcategory_id
func TestAddIncomeController_AddIncome_EmptySubcategoryID(t *testing.T) {
	// Arrange
	ctrl, walletID, _ := setupAddIncomeController(t)
	
	requestBody := map[string]interface{}{
		"wallet_id":      walletID,
		"subcategory_id": "", // Empty string
		"amount":         10000,
		"currency":       "USD",
		"description":    "Test income",
		"date":          time.Now().Format(time.RFC3339),
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	// Act
	ctrl.AddIncome(w, req)
	
	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["error"] != "subcategory_id is required" {
		t.Errorf("Expected error message 'subcategory_id is required', got %v", response["error"])
	}
}

// TestAddIncomeController_AddIncome_InvalidAmount tests invalid amount values
func TestAddIncomeController_AddIncome_InvalidAmount(t *testing.T) {
	ctrl, walletID, subcategoryID := setupAddIncomeController(t)
	
	testCases := []struct {
		name     string
		amount   interface{}
		expected string
	}{
		{"Zero Amount", 0, "amount must be positive"},
		{"Negative Amount", -1000, "amount must be positive"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			requestBody := map[string]interface{}{
				"wallet_id":      walletID,
				"subcategory_id": subcategoryID,
				"amount":         tc.amount,
				"currency":       "USD",
				"description":    "Test income",
				"date":          time.Now().Format(time.RFC3339),
			}
			
			jsonBody, _ := json.Marshal(requestBody)
			req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			// Act
			ctrl.AddIncome(w, req)
			
			// Assert
			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
			}
			
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			
			if response["error"] != tc.expected {
				t.Errorf("Expected error message '%s', got %v", tc.expected, response["error"])
			}
		})
	}
}

// TestAddIncomeController_AddIncome_MissingCurrency tests missing currency
func TestAddIncomeController_AddIncome_MissingCurrency(t *testing.T) {
	// Arrange
	ctrl, walletID, subcategoryID := setupAddIncomeController(t)
	
	requestBody := map[string]interface{}{
		"wallet_id":      walletID,
		"subcategory_id": subcategoryID,
		"amount":         10000,
		// "currency" intentionally omitted
		"description": "Test income",
		"date":       time.Now().Format(time.RFC3339),
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	// Act
	ctrl.AddIncome(w, req)
	
	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["error"] != "currency is required" {
		t.Errorf("Expected error message 'currency is required', got %v", response["error"])
	}
}

// TestAddIncomeController_AddIncome_EmptyCurrency tests empty currency
func TestAddIncomeController_AddIncome_EmptyCurrency(t *testing.T) {
	// Arrange
	ctrl, walletID, subcategoryID := setupAddIncomeController(t)
	
	requestBody := map[string]interface{}{
		"wallet_id":      walletID,
		"subcategory_id": subcategoryID,
		"amount":         10000,
		"currency":       "", // Empty string
		"description":    "Test income",
		"date":          time.Now().Format(time.RFC3339),
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	// Act
	ctrl.AddIncome(w, req)
	
	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["error"] != "currency is required" {
		t.Errorf("Expected error message 'currency is required', got %v", response["error"])
	}
}

// TestAddIncomeController_AddIncome_WalletNotFound tests non-existent wallet
func TestAddIncomeController_AddIncome_WalletNotFound(t *testing.T) {
	// Arrange
	ctrl, _, subcategoryID := setupAddIncomeController(t)
	
	requestBody := map[string]interface{}{
		"wallet_id":      "non-existent-wallet-id",
		"subcategory_id": subcategoryID,
		"amount":         10000,
		"currency":       "USD",
		"description":    "Test income",
		"date":          time.Now().Format(time.RFC3339),
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	// Act
	ctrl.AddIncome(w, req)
	
	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["success"] != false {
		t.Errorf("Expected success to be false, got %v", response["success"])
	}
	
	if response["message"] != "Wallet not found" {
		t.Errorf("Expected message 'Wallet not found', got %v", response["message"])
	}
}

// TestAddIncomeController_AddIncome_SubcategoryNotFound tests non-existent subcategory
func TestAddIncomeController_AddIncome_SubcategoryNotFound(t *testing.T) {
	// Arrange
	ctrl, walletID, _ := setupAddIncomeController(t)
	
	requestBody := map[string]interface{}{
		"wallet_id":      walletID,
		"subcategory_id": "non-existent-subcategory",
		"amount":         10000,
		"currency":       "USD",
		"description":    "Test income",
		"date":          time.Now().Format(time.RFC3339),
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	// Act
	ctrl.AddIncome(w, req)
	
	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["success"] != false {
		t.Errorf("Expected success to be false, got %v", response["success"])
	}
	
	if response["message"] != "Subcategory not found in any category" {
		t.Errorf("Expected message 'Subcategory not found in any category', got %v", response["message"])
	}
}

// TestAddIncomeController_AddIncome_FloatAmount tests that float amounts cause JSON error
func TestAddIncomeController_AddIncome_FloatAmount(t *testing.T) {
	// Arrange
	ctrl, walletID, subcategoryID := setupAddIncomeController(t)
	
	// Create a JSON string manually to include float value
	jsonString := `{
		"wallet_id": "` + walletID + `",
		"subcategory_id": "` + subcategoryID + `",
		"amount": 100.50,
		"currency": "USD",
		"description": "Test income",
		"date": "` + time.Now().Format(time.RFC3339) + `"
	}`
	
	req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBufferString(jsonString))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	// Act
	ctrl.AddIncome(w, req)
	
	// Assert - JSON unmarshaling fails when trying to unmarshal float into int64
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["error"] != "Invalid JSON" {
		t.Errorf("Expected error message 'Invalid JSON', got %v", response["error"])
	}
}

// TestAddIncomeController_AddIncome_EmptyRequestBody tests completely empty request
func TestAddIncomeController_AddIncome_EmptyRequestBody(t *testing.T) {
	// Arrange
	ctrl, _, _ := setupAddIncomeController(t)
	
	req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	// Act
	ctrl.AddIncome(w, req)
	
	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["error"] != "wallet_id is required" {
		t.Errorf("Expected first validation error 'wallet_id is required', got %v", response["error"])
	}
}

// TestAddIncomeController_AddIncome_NoContentType tests request without Content-Type
func TestAddIncomeController_AddIncome_NoContentType(t *testing.T) {
	// Arrange
	ctrl, walletID, subcategoryID := setupAddIncomeController(t)
	
	requestBody := map[string]interface{}{
		"wallet_id":      walletID,
		"subcategory_id": subcategoryID,
		"amount":         10000,
		"currency":       "USD",
		"description":    "Test income",
		"date":          time.Now().Format(time.RFC3339),
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBuffer(jsonBody))
	// Intentionally not setting Content-Type
	w := httptest.NewRecorder()
	
	// Act
	ctrl.AddIncome(w, req)
	
	// Assert - Should still work as Go's JSON decoder doesn't require Content-Type
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
	}
}

// TestAddIncomeController_AddIncome_OptionalFieldsSuccess tests that description and date are optional
func TestAddIncomeController_AddIncome_OptionalFieldsSuccess(t *testing.T) {
	// Arrange
	ctrl, walletID, subcategoryID := setupAddIncomeController(t)
	
	requestBody := map[string]interface{}{
		"wallet_id":      walletID,
		"subcategory_id": subcategoryID,
		"amount":         10000,
		"currency":       "USD",
		// description and date omitted to test they're optional
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/incomes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	// Act
	ctrl.AddIncome(w, req)
	
	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["success"] != true {
		t.Errorf("Expected success to be true, got %v", response["success"])
	}
}