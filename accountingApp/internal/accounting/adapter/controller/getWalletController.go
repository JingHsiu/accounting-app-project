package controller

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

// QueryWalletController represents the controller responsible for wallet queries
type QueryWalletController struct {
	getWalletsUseCase usecase.GetWalletsUseCase
	getWalletUseCase  usecase.GetWalletUseCase
}

// NewQueryWalletController creates a new instance of QueryWalletController
func NewQueryWalletController(getWalletsUseCase usecase.GetWalletsUseCase, getWalletUseCase usecase.GetWalletUseCase) *QueryWalletController {
	return &QueryWalletController{
		getWalletsUseCase: getWalletsUseCase,
		getWalletUseCase:  getWalletUseCase,
	}
}

// GetWallets handles GET /api/v1/wallets?userID={userID}
func (c *QueryWalletController) GetWallets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		c.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		c.sendError(w, "userID parameter is required", http.StatusBadRequest)
		return
	}

	result := c.getWalletsUseCase.Execute(usecase.GetWalletsInput{
		UserID: userID,
	})

	if result.GetExitCode() != common.Success {
		c.sendError(w, result.GetMessage(), http.StatusInternalServerError)
		return
	}

	// Convert domain models to API response format
	output, ok := result.(usecase.GetWalletsOutput)
	if !ok {
		c.sendError(w, "Internal error: invalid output type", http.StatusInternalServerError)
		return
	}

	response := make([]map[string]interface{}, len(output.Wallets))
	for i, wallet := range output.Wallets {
		response[i] = c.walletToResponse(wallet)
	}

	c.sendSuccess(w, map[string]interface{}{
		"data":  response,
		"count": len(response),
	})
}

// GetWallet handles GET /api/v1/wallets/{walletID}
func (c *QueryWalletController) GetWallet(w http.ResponseWriter, r *http.Request) {
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
	includeTransactions := r.URL.Query().Get("includeTransactions") == "true"

	result := c.getWalletUseCase.Execute(usecase.GetWalletInput{
		WalletID:            walletID,
		IncludeTransactions: includeTransactions,
	})

	if result.GetExitCode() != common.Success {
		if result.GetMessage() == "Wallet not found" {
			c.sendError(w, result.GetMessage(), http.StatusNotFound)
		} else {
			c.sendError(w, result.GetMessage(), http.StatusInternalServerError)
		}
		return
	}

	// Convert domain model to API response format
	output, ok := result.(usecase.GetWalletOutput)
	if !ok {
		c.sendError(w, "Internal error: invalid output type", http.StatusInternalServerError)
		return
	}

	response := c.walletToResponse(output.Wallet)
	c.sendSuccess(w, map[string]interface{}{
		"data": response,
	})
}

// Helper methods
func (c *QueryWalletController) extractWalletID(path string) string {
	// Extract wallet ID from paths like /api/v1/wallets/{walletID}
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

func (c *QueryWalletController) walletToResponse(wallet *model.Wallet) map[string]interface{} {
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
		response["transactions"] = transactions
		response["is_fully_loaded"] = true
	} else {
		response["is_fully_loaded"] = false
	}

	return response
}

func (c *QueryWalletController) sendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func (c *QueryWalletController) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
