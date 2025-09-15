package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
)

// QueryExpenseController handles expense query operations
type QueryExpenseController struct {
	getExpensesUseCase usecase.GetExpensesUseCase
}

// NewQueryExpenseController creates a new QueryExpenseController
func NewQueryExpenseController(getExpensesUseCase usecase.GetExpensesUseCase) *QueryExpenseController {
	return &QueryExpenseController{
		getExpensesUseCase: getExpensesUseCase,
	}
}

// GetExpenses handles GET /api/v1/expenses
func (c *QueryExpenseController) GetExpenses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract query parameters
	query := r.URL.Query()
	
	// For now, use a demo user ID (in production this would come from auth)
	userID := "demo-user-123"
	if queryUserID := query.Get("userID"); queryUserID != "" {
		userID = queryUserID
	}

	input := usecase.GetExpensesInput{
		UserID: userID,
	}

	// Process optional filters
	if walletID := query.Get("walletID"); walletID != "" {
		input.WalletID = &walletID
	}

	if categoryID := query.Get("categoryID"); categoryID != "" {
		input.CategoryID = &categoryID
	}

	if startDateStr := query.Get("startDate"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			input.StartDate = &startDate
		}
	}

	if endDateStr := query.Get("endDate"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			input.EndDate = &endDate
		}
	}

	if minAmountStr := query.Get("minAmount"); minAmountStr != "" {
		if minAmount, err := strconv.ParseInt(minAmountStr, 10, 64); err == nil {
			input.MinAmount = &minAmount
		}
	}

	if maxAmountStr := query.Get("maxAmount"); maxAmountStr != "" {
		if maxAmount, err := strconv.ParseInt(maxAmountStr, 10, 64); err == nil {
			input.MaxAmount = &maxAmount
		}
	}

	if description := query.Get("description"); description != "" {
		input.Description = &description
	}

	// Execute use case
	output := c.getExpensesUseCase.Execute(input)

	w.Header().Set("Content-Type", "application/json")
	
	if output.GetExitCode() != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   output.GetMessage(),
		})
		return
	}

	// Cast to specific output type to access data
	expensesOutput, ok := output.(usecase.GetExpensesOutput)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid output type",
		})
		return
	}

	// Return successful response in format expected by frontend
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    expensesOutput.Data,
		"count":   expensesOutput.Count,
		"message": expensesOutput.Message,
	})
}

// Helper methods
func (c *QueryExpenseController) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}