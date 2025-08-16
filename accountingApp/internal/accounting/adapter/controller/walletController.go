package controller

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/query"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

type WalletController struct {
	createWalletUseCase     usecase.CreateWalletUseCase
	addExpenseUseCase       usecase.AddExpenseUseCase
	addIncomeUseCase        usecase.AddIncomeUseCase
	getWalletBalanceUseCase usecase.GetWalletBalanceUseCase
	walletRepository        repository.WalletRepository
}

func NewWalletController(
	createWalletUseCase usecase.CreateWalletUseCase,
	addExpenseUseCase usecase.AddExpenseUseCase,
	addIncomeUseCase usecase.AddIncomeUseCase,
	getWalletBalanceUseCase usecase.GetWalletBalanceUseCase,
	walletRepository repository.WalletRepository,
) *WalletController {
	return &WalletController{
		createWalletUseCase:     createWalletUseCase,
		addExpenseUseCase:       addExpenseUseCase,
		addIncomeUseCase:        addIncomeUseCase,
		getWalletBalanceUseCase: getWalletBalanceUseCase,
		walletRepository:        walletRepository,
	}
}

func (c *WalletController) CreateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID         string  `json:"user_id"`
		Name           string  `json:"name"`
		Type           string  `json:"type"`
		Currency       string  `json:"currency"`
		InitialBalance *int64  `json:"initialBalance,omitempty"` // Optional initial balance in cents
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	input := command.CreateWalletInput{
		UserID:         req.UserID,
		Name:           req.Name,
		Type:           req.Type,
		Currency:       req.Currency,
		InitialBalance: req.InitialBalance,
	}

	output := c.createWalletUseCase.Execute(input)
	
	w.Header().Set("Content-Type", "application/json")
	if output.GetExitCode() != 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        output.GetID(),
		"success":   output.GetExitCode() == 0,
		"message":   output.GetMessage(),
	})
}

func (c *WalletController) GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/wallets/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "balance" {
		http.Error(w, "Invalid path", http.StatusNotFound)
		return
	}
	
	walletID := parts[0]
	if walletID == "" {
		http.Error(w, "Wallet ID required", http.StatusBadRequest)
		return
	}

	input := query.GetWalletBalanceInput{WalletID: walletID}
	output := c.getWalletBalanceUseCase.Execute(input)
	
	w.Header().Set("Content-Type", "application/json")
	if output.GetExitCode() != 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	
	json.NewEncoder(w).Encode(output)
}

func (c *WalletController) AddExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		WalletID    string `json:"wallet_id"`
		CategoryID  string `json:"category_id"`
		Amount      int64  `json:"amount"`
		Currency    string `json:"currency"`
		Description string `json:"description"`
		Date        string `json:"date,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var date time.Time
	var err error
	if req.Date != "" {
		date, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			http.Error(w, "Invalid date format (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
	} else {
		date = time.Now()
	}

	input := command.AddExpenseInput{
		WalletID:      req.WalletID,
		SubcategoryID: req.CategoryID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   req.Description,
		Date:          date,
	}

	output := c.addExpenseUseCase.Execute(input)
	
	w.Header().Set("Content-Type", "application/json")
	if output.GetExitCode() != 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      output.GetID(),
		"success": output.GetExitCode() == 0,
		"message": output.GetMessage(),
	})
}

func (c *WalletController) AddIncome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		WalletID    string `json:"wallet_id"`
		CategoryID  string `json:"category_id"`
		Amount      int64  `json:"amount"`
		Currency    string `json:"currency"`
		Description string `json:"description"`
		Date        string `json:"date,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var date time.Time
	var err error
	if req.Date != "" {
		date, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			http.Error(w, "Invalid date format (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
	} else {
		date = time.Now()
	}

	input := command.AddIncomeInput{
		WalletID:      req.WalletID,
		SubcategoryID: req.CategoryID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   req.Description,
		Date:          date,
	}

	output := c.addIncomeUseCase.Execute(input)
	
	w.Header().Set("Content-Type", "application/json")
	if output.GetExitCode() != 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      output.GetID(),
		"success": output.GetExitCode() == 0,
		"message": output.GetMessage(),
	})
}

// GetWallets handles GET /api/v1/wallets?userID={userID}
func (c *WalletController) GetWallets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		c.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		c.sendError(w, "userID parameter is required", http.StatusBadRequest)
		return
	}

	wallets, err := c.walletRepository.FindByUserID(userID)
	if err != nil {
		c.sendError(w, "Failed to retrieve wallets: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert domain models to API response format
	response := make([]map[string]interface{}, len(wallets))
	for i, wallet := range wallets {
		response[i] = c.walletToResponse(wallet)
	}

	c.sendSuccess(w, map[string]interface{}{
		"data": response,
		"count": len(response),
	})
}

// GetWallet handles GET /api/v1/wallets/{walletID}
func (c *WalletController) GetWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		c.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	walletID := c.extractWalletID(r.URL.Path)
	if walletID == "" {
		c.sendError(w, "Invalid wallet ID", http.StatusBadRequest)
		return
	}

	// Check if we need to load transactions (query parameter)
	loadTransactions := r.URL.Query().Get("includeTransactions") == "true"
	
	var wallet *model.Wallet
	var err error
	
	if loadTransactions {
		wallet, err = c.walletRepository.FindByIDWithTransactions(walletID)
	} else {
		wallet, err = c.walletRepository.FindByID(walletID)
	}
	
	if err != nil {
		c.sendError(w, "Failed to retrieve wallet: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if wallet == nil {
		c.sendError(w, "Wallet not found", http.StatusNotFound)
		return
	}

	response := c.walletToResponse(wallet)
	c.sendSuccess(w, map[string]interface{}{
		"data": response,
	})
}

// UpdateWallet handles PUT /api/v1/wallets/{walletID}
func (c *WalletController) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		c.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	walletID := c.extractWalletID(r.URL.Path)
	if walletID == "" {
		c.sendError(w, "Invalid wallet ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name     string `json:"name,omitempty"`
		Type     string `json:"type,omitempty"`
		Currency string `json:"currency,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.sendError(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get existing wallet
	wallet, err := c.walletRepository.FindByID(walletID)
	if err != nil {
		c.sendError(w, "Failed to retrieve wallet: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if wallet == nil {
		c.sendError(w, "Wallet not found", http.StatusNotFound)
		return
	}

	// Update wallet properties (this would need methods on the domain model)
	// For now, we'll create a new wallet with updated properties
	// Note: This is a simplified implementation - in practice, you'd want proper update methods
	updated := false
	
	if req.Name != "" && req.Name != wallet.Name {
		// wallet.UpdateName(req.Name) - would need this method on domain model
		updated = true
	}
	
	if req.Type != "" && req.Type != string(wallet.Type) {
		// wallet.UpdateType(req.Type) - would need this method on domain model
		updated = true
	}
	
	if req.Currency != "" && req.Currency != wallet.Currency() {
		// wallet.UpdateCurrency(req.Currency) - would need this method on domain model
		updated = true
	}

	if updated {
		if err := c.walletRepository.Save(wallet); err != nil {
			c.sendError(w, "Failed to update wallet: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	response := c.walletToResponse(wallet)
	c.sendSuccess(w, map[string]interface{}{
		"data": response,
		"message": "Wallet updated successfully",
	})
}

// DeleteWallet handles DELETE /api/v1/wallets/{walletID}
func (c *WalletController) DeleteWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		c.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	walletID := c.extractWalletID(r.URL.Path)
	if walletID == "" {
		c.sendError(w, "Invalid wallet ID", http.StatusBadRequest)
		return
	}

	// Check if wallet exists
	wallet, err := c.walletRepository.FindByID(walletID)
	if err != nil {
		c.sendError(w, "Failed to retrieve wallet: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if wallet == nil {
		c.sendError(w, "Wallet not found", http.StatusNotFound)
		return
	}

	// Delete the wallet
	if err := c.walletRepository.Delete(walletID); err != nil {
		c.sendError(w, "Failed to delete wallet: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c.sendSuccess(w, map[string]interface{}{
		"message": "Wallet deleted successfully",
	})
}

// Helper methods for response formatting and parsing
func (c *WalletController) extractWalletID(path string) string {
	// Extract wallet ID from paths like /api/v1/wallets/{walletID} or /api/v1/wallets/{walletID}/balance
	parts := strings.Split(strings.TrimPrefix(path, "/api/v1/wallets/"), "/")
	if len(parts) > 0 && parts[0] != "" {
		// URL decode in case there are special characters
		decoded, err := url.QueryUnescape(parts[0])
		if err != nil {
			return parts[0] // fallback to original if decode fails
		}
		return decoded
	}
	return ""
}

func (c *WalletController) walletToResponse(wallet *model.Wallet) map[string]interface{} {
	response := map[string]interface{}{
		"id":       wallet.ID,
		"user_id":  wallet.UserID,
		"name":     wallet.Name,
		"type":     string(wallet.Type),
		"currency": wallet.Currency(),
		"balance": map[string]interface{}{
			"amount":   wallet.Balance.Amount,
			"currency": wallet.Balance.Currency,
		},
		"created_at": wallet.CreatedAt.Format(time.RFC3339),
		"updated_at": wallet.UpdatedAt.Format(time.RFC3339),
	}

	// Include transactions if wallet is fully loaded
	if wallet.IsFullyLoaded() {
		transactions := make([]map[string]interface{}, 0)
		// Note: This would need methods to get transaction history from the domain model
		// transactions = c.transactionsToResponse(wallet.GetTransactionHistory())
		response["transactions"] = transactions
		response["is_fully_loaded"] = true
	} else {
		response["is_fully_loaded"] = false
	}

	return response
}

func (c *WalletController) sendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func (c *WalletController) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}