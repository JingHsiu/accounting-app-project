package controller

import (
	"encoding/json"
	"net/http"

	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
)

// CategoryController handles category operations
type CategoryController struct {
	createExpenseCategoryUseCase usecase.CreateExpenseCategoryUseCase
	createIncomeCategoryUseCase  usecase.CreateIncomeCategoryUseCase
	getExpenseCategoriesUseCase  usecase.GetExpenseCategoriesUseCase
	getIncomeCategoriesUseCase   usecase.GetIncomeCategoriesUseCase
}

// NewCategoryController creates a new CategoryController
func NewCategoryController(
	createExpenseCategoryUseCase usecase.CreateExpenseCategoryUseCase,
	createIncomeCategoryUseCase usecase.CreateIncomeCategoryUseCase,
	getExpenseCategoriesUseCase usecase.GetExpenseCategoriesUseCase,
	getIncomeCategoriesUseCase usecase.GetIncomeCategoriesUseCase,
) *CategoryController {
	return &CategoryController{
		createExpenseCategoryUseCase: createExpenseCategoryUseCase,
		createIncomeCategoryUseCase:  createIncomeCategoryUseCase,
		getExpenseCategoriesUseCase:  getExpenseCategoriesUseCase,
		getIncomeCategoriesUseCase:   getIncomeCategoriesUseCase,
	}
}

// CreateExpenseCategory handles POST /api/v1/categories/expense
func (c *CategoryController) CreateExpenseCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID string `json:"user_id"`
		Name   string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.UserID == "" {
		c.sendError(w, "user_id is required", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		c.sendError(w, "name is required", http.StatusBadRequest)
		return
	}

	input := usecase.CreateExpenseCategoryInput{
		UserID: req.UserID,
		Name:   req.Name,
	}

	output := c.createExpenseCategoryUseCase.Execute(input)

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

// CreateIncomeCategory handles POST /api/v1/categories/income
func (c *CategoryController) CreateIncomeCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID string `json:"user_id"`
		Name   string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.UserID == "" {
		c.sendError(w, "user_id is required", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		c.sendError(w, "name is required", http.StatusBadRequest)
		return
	}

	input := usecase.CreateIncomeCategoryInput{
		UserID: req.UserID,
		Name:   req.Name,
	}

	output := c.createIncomeCategoryUseCase.Execute(input)

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

// GetExpenseCategories handles GET /api/v1/categories/expense
func (c *CategoryController) GetExpenseCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		c.sendError(w, "userID is required", http.StatusBadRequest)
		return
	}

	input := usecase.GetExpenseCategoriesInput{
		UserID: userID,
	}

	output := c.getExpenseCategoriesUseCase.Execute(input)

	w.Header().Set("Content-Type", "application/json")
	if output.GetExitCode() != 0 {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Convert to expected frontend format
	if categoryOutput, ok := output.(usecase.GetExpenseCategoriesOutput); ok {
		response := map[string]interface{}{
			"success": output.GetExitCode() == 0,
			"message": output.GetMessage(),
		}
		if output.GetExitCode() == 0 {
			response["data"] = categoryOutput.Categories
		}
		json.NewEncoder(w).Encode(response)
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": output.GetExitCode() == 0,
			"message": output.GetMessage(),
		})
	}
}

// GetIncomeCategories handles GET /api/v1/categories/income
func (c *CategoryController) GetIncomeCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		c.sendError(w, "userID is required", http.StatusBadRequest)
		return
	}

	input := usecase.GetIncomeCategoriesInput{
		UserID: userID,
	}

	output := c.getIncomeCategoriesUseCase.Execute(input)

	w.Header().Set("Content-Type", "application/json")
	if output.GetExitCode() != 0 {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Convert to expected frontend format
	if categoryOutput, ok := output.(usecase.GetIncomeCategoriesOutput); ok {
		response := map[string]interface{}{
			"success": output.GetExitCode() == 0,
			"message": output.GetMessage(),
		}
		if output.GetExitCode() == 0 {
			response["data"] = categoryOutput.Categories
		}
		json.NewEncoder(w).Encode(response)
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": output.GetExitCode() == 0,
			"message": output.GetMessage(),
		})
	}
}

// Helper methods
func (c *CategoryController) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}