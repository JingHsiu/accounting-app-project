package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
)

// AddIncomeController handles add income operations
type AddIncomeController struct {
	addIncomeUseCase usecase.AddIncomeUseCase
}

// NewAddIncomeController creates a new AddIncomeController
func NewAddIncomeController(addIncomeUseCase usecase.AddIncomeUseCase) *AddIncomeController {
	return &AddIncomeController{
		addIncomeUseCase: addIncomeUseCase,
	}
}

// AddIncome handles POST /api/v1/incomes
func (c *AddIncomeController) AddIncome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		WalletID      string    `json:"wallet_id"`
		SubcategoryID string    `json:"subcategory_id"`
		Amount        int64     `json:"amount"`
		Currency      string    `json:"currency"`
		Description   string    `json:"description"`
		Date          time.Time `json:"date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.WalletID == "" {
		c.sendError(w, "wallet_id is required", http.StatusBadRequest)
		return
	}
	if req.SubcategoryID == "" {
		c.sendError(w, "subcategory_id is required", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		c.sendError(w, "amount must be positive", http.StatusBadRequest)
		return
	}
	if req.Currency == "" {
		c.sendError(w, "currency is required", http.StatusBadRequest)
		return
	}

	input := usecase.AddIncomeInput{
		WalletID:      req.WalletID,
		SubcategoryID: req.SubcategoryID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   req.Description,
		Date:          req.Date,
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

// Helper methods
func (c *AddIncomeController) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}