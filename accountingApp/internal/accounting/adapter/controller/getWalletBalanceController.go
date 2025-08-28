package controller

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
)

// GetWalletBalanceController handles get wallet balance operations
type GetWalletBalanceController struct {
	getWalletBalanceUseCase usecase.GetWalletBalanceUseCase
}

// NewGetWalletBalanceController creates a new GetWalletBalanceController
func NewGetWalletBalanceController(getWalletBalanceUseCase usecase.GetWalletBalanceUseCase) *GetWalletBalanceController {
	return &GetWalletBalanceController{
		getWalletBalanceUseCase: getWalletBalanceUseCase,
	}
}

// GetWalletBalance handles GET /api/v1/wallets/{id}/balance
func (c *GetWalletBalanceController) GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract wallet ID from URL path like /api/v1/wallets/{id}/balance
	walletID := c.extractWalletIDFromBalancePath(r.URL.Path)
	if walletID == "" {
		c.sendError(w, "Invalid wallet ID", http.StatusBadRequest)
		return
	}

	input := usecase.GetWalletBalanceInput{
		WalletID: walletID,
	}

	output := c.getWalletBalanceUseCase.Execute(input)

	w.Header().Set("Content-Type", "application/json")
	if output.GetExitCode() != 0 {
		if output.GetMessage() == "Wallet not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	// Cast to specific output type to access fields
	balanceOutput, ok := output.(usecase.GetWalletBalanceOutput)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Internal error: invalid output type",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"walletId": balanceOutput.ID,
		"balance":  balanceOutput.Balance,
		"currency": balanceOutput.Currency,
		"success":  output.GetExitCode() == 0,
		"message":  output.GetMessage(),
	})
}

// Helper methods
func (c *GetWalletBalanceController) extractWalletIDFromBalancePath(path string) string {
	// Extract from paths like /api/v1/wallets/{walletID}/balance
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "wallets" && i+1 < len(parts) {
			walletID := parts[i+1]
			// URL decode in case there are special characters
			decoded, err := url.QueryUnescape(walletID)
			if err != nil {
				return walletID // fallback to original if decode fails
			}
			return decoded
		}
	}
	return ""
}

func (c *GetWalletBalanceController) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}