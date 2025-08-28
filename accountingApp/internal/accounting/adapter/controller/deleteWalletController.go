package controller

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
)

// DeleteWalletController represents the controller responsible for wallet deletion
type DeleteWalletController struct {
	deleteWalletUseCase usecase.DeleteWalletUseCase
}

// NewDeleteWalletController creates a new instance of DeleteWalletController
func NewDeleteWalletController(deleteWalletUseCase usecase.DeleteWalletUseCase) *DeleteWalletController {
	return &DeleteWalletController{
		deleteWalletUseCase: deleteWalletUseCase,
	}
}

// DeleteWallet handles DELETE /api/v1/wallets/{walletID}
func (c *DeleteWalletController) DeleteWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		c.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	walletID := c.extractWalletID(r.URL.Path)
	if walletID == "" {
		c.sendError(w, "Invalid wallet ID", http.StatusBadRequest)
		return
	}

	result := c.deleteWalletUseCase.Execute(usecase.DeleteWalletInput{
		WalletID: walletID,
	})

	if result.GetExitCode() != common.Success {
		message := result.GetMessage()
		if message == "Wallet not found" {
			c.sendError(w, message, http.StatusNotFound)
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
func (c *DeleteWalletController) extractWalletID(path string) string {
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

func (c *DeleteWalletController) sendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func (c *DeleteWalletController) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}