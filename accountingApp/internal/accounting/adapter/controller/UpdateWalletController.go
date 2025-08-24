package controller

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
)

// UpdateWalletController represents the controller responsible for wallet updates
type UpdateWalletController struct {
	updateWalletUseCase usecase.UpdateWalletUseCase
}

// NewUpdateWalletController creates a new instance of UpdateWalletController
func NewUpdateWalletController(updateWalletUseCase usecase.UpdateWalletUseCase) *UpdateWalletController {
	return &UpdateWalletController{
		updateWalletUseCase: updateWalletUseCase,
	}
}

// UpdateWallet handles PUT /api/v1/wallets/{walletID}
func (c *UpdateWalletController) UpdateWallet(w http.ResponseWriter, r *http.Request) {
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

	// Prepare input for use case
	var name, walletType, currency *string
	if req.Name != "" {
		name = &req.Name
	}
	if req.Type != "" {
		walletType = &req.Type
	}
	if req.Currency != "" {
		currency = &req.Currency
	}

	result := c.updateWalletUseCase.Execute(command.UpdateWalletInput{
		WalletID: walletID,
		Name:     name,
		Type:     walletType,
		Currency: currency,
	})

	if result.GetExitCode() != common.Success {
		message := result.GetMessage()
		if message == "Wallet not found" {
			c.sendError(w, message, http.StatusNotFound)
		} else if strings.Contains(message, "Invalid") {
			c.sendError(w, message, http.StatusBadRequest)
		} else {
			c.sendError(w, message, http.StatusInternalServerError)
		}
		return
	}

	c.sendSuccess(w, map[string]interface{}{
		"message": result.GetMessage(),
	})
}

// Helper methods
func (c *UpdateWalletController) extractWalletID(path string) string {
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


func (c *UpdateWalletController) sendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func (c *UpdateWalletController) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}