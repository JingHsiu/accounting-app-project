package controller

import (
	"encoding/json"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
	"net/http"

	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
)

// CreateWalletController represents the controller responsible for wallet creation
type CreateWalletController struct {
	createWalletUseCase usecase.CreateWalletUseCase
}

// NewCreateWalletController creates a new instance of CreateWalletController
func NewCreateWalletController(createWalletUseCase usecase.CreateWalletUseCase) *CreateWalletController {
	return &CreateWalletController{
		createWalletUseCase: createWalletUseCase,
	}
}

// CreateWallet handles POST /api/v1/wallets
func (c *CreateWalletController) CreateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID         string `json:"user_id"`
		Name           string `json:"name"`
		Type           string `json:"type"`
		Currency       string `json:"currency"`
		InitialBalance *int64 `json:"initialBalance,omitempty"`
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
	if req.Type == "" {
		c.sendError(w, "type is required", http.StatusBadRequest)
		return
	}
	if req.Currency == "" {
		c.sendError(w, "currency is required", http.StatusBadRequest)
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
		"id":      output.GetID(),
		"success": output.GetExitCode() == 0,
		"message": output.GetMessage(),
	})
}

func (c *CreateWalletController) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
