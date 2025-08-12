package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/query"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
)

type WalletController struct {
	createWalletUseCase     usecase.CreateWalletUseCase
	addExpenseUseCase       usecase.AddExpenseUseCase
	addIncomeUseCase        usecase.AddIncomeUseCase
	getWalletBalanceUseCase usecase.GetWalletBalanceUseCase
}

func NewWalletController(
	createWalletUseCase usecase.CreateWalletUseCase,
	addExpenseUseCase usecase.AddExpenseUseCase,
	addIncomeUseCase usecase.AddIncomeUseCase,
	getWalletBalanceUseCase usecase.GetWalletBalanceUseCase,
) *WalletController {
	return &WalletController{
		createWalletUseCase:     createWalletUseCase,
		addExpenseUseCase:       addExpenseUseCase,
		addIncomeUseCase:        addIncomeUseCase,
		getWalletBalanceUseCase: getWalletBalanceUseCase,
	}
}

func (c *WalletController) CreateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID   string `json:"user_id"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		Currency string `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	input := command.CreateWalletInput{
		UserID:   req.UserID,
		Name:     req.Name,
		Type:     req.Type,
		Currency: req.Currency,
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