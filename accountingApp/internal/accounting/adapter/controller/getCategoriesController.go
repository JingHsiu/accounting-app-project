package controller

import (
	"encoding/json"
	"net/http"
)

// Category represents a category structure for API responses
type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// GetCategoriesController handles categories endpoint
type GetCategoriesController struct{}

// NewGetCategoriesController creates a new categories controller
func NewGetCategoriesController() *GetCategoriesController {
	return &GetCategoriesController{}
}

// GetCategories handles GET /api/v1/categories
func (c *GetCategoriesController) GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Return hardcoded default categories for now
	// This will be replaced with database queries later
	categories := []Category{
		// Expense categories
		{ID: "default-expense-1", Name: "餐飲", Type: "expense"},
		{ID: "default-expense-2", Name: "交通", Type: "expense"},
		{ID: "default-expense-3", Name: "購物", Type: "expense"},
		{ID: "default-expense-4", Name: "娛樂", Type: "expense"},
		{ID: "default-expense-5", Name: "醫療", Type: "expense"},
		{ID: "default-expense-6", Name: "教育", Type: "expense"},
		{ID: "default-expense-7", Name: "居住", Type: "expense"},
		{ID: "default-expense-8", Name: "其他", Type: "expense"},
		
		// Income categories
		{ID: "default-income-1", Name: "薪資", Type: "income"},
		{ID: "default-income-2", Name: "投資", Type: "income"},
		{ID: "default-income-3", Name: "副業", Type: "income"},
		{ID: "default-income-4", Name: "其他收入", Type: "income"},
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"data":    categories,
	}
	json.NewEncoder(w).Encode(response)
}

// GetExpenseCategories handles GET /api/v1/categories/expense
func (c *GetCategoriesController) GetExpenseCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	expenseCategories := []Category{
		{ID: "default-expense-1", Name: "餐飲", Type: "expense"},
		{ID: "default-expense-2", Name: "交通", Type: "expense"},
		{ID: "default-expense-3", Name: "購物", Type: "expense"},
		{ID: "default-expense-4", Name: "娛樂", Type: "expense"},
		{ID: "default-expense-5", Name: "醫療", Type: "expense"},
		{ID: "default-expense-6", Name: "教育", Type: "expense"},
		{ID: "default-expense-7", Name: "居住", Type: "expense"},
		{ID: "default-expense-8", Name: "其他", Type: "expense"},
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"data":    expenseCategories,
	}
	json.NewEncoder(w).Encode(response)
}

// GetIncomeCategories handles GET /api/v1/categories/income
func (c *GetCategoriesController) GetIncomeCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	incomeCategories := []Category{
		{ID: "default-income-1", Name: "薪資", Type: "income"},
		{ID: "default-income-2", Name: "投資", Type: "income"},
		{ID: "default-income-3", Name: "副業", Type: "income"},
		{ID: "default-income-4", Name: "其他收入", Type: "income"},
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"data":    incomeCategories,
	}
	json.NewEncoder(w).Encode(response)
}