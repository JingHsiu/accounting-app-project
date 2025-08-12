package controller

import (
	"encoding/json"
	"net/http"

	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
)

type CategoryController struct {
	createExpenseCategoryUseCase usecase.CreateExpenseCategoryUseCase
	createIncomeCategoryUseCase  usecase.CreateIncomeCategoryUseCase
}

func NewCategoryController(
	createExpenseCategoryUseCase usecase.CreateExpenseCategoryUseCase,
	createIncomeCategoryUseCase usecase.CreateIncomeCategoryUseCase,
) *CategoryController {
	return &CategoryController{
		createExpenseCategoryUseCase: createExpenseCategoryUseCase,
		createIncomeCategoryUseCase:  createIncomeCategoryUseCase,
	}
}

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
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	input := command.CreateExpenseCategoryInput{
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
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	input := command.CreateIncomeCategoryInput{
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
