package web

import (
	"net/http"
	"strings"

	"github.com/JingHsiu/accountingApp/internal/accounting/adapter/controller"
)

type Router struct {
	walletController   *controller.WalletController
	categoryController *controller.CategoryController
}

func NewRouter(
	walletController *controller.WalletController,
	categoryController *controller.CategoryController,
) *Router {
	return &Router{
		walletController:   walletController,
		categoryController: categoryController,
	}
}

func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	
	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Wallet endpoints - REST API design
	mux.HandleFunc("/api/v1/wallets", r.handleWalletCollection)          // GET (with userID param), POST
	mux.HandleFunc("/api/v1/wallets/", r.handleWalletResource)           // GET, PUT, DELETE by ID
	mux.HandleFunc("/api/v1/wallets/balance/", r.walletController.GetWalletBalance) // Keep existing balance endpoint

	// Category endpoints
	mux.HandleFunc("/api/v1/categories/expense", r.categoryController.CreateExpenseCategory)
	mux.HandleFunc("/api/v1/categories/income", r.categoryController.CreateIncomeCategory)

	// Transaction endpoints
	mux.HandleFunc("/api/v1/expenses", r.walletController.AddExpense)
	mux.HandleFunc("/api/v1/incomes", r.walletController.AddIncome)

	return mux
}

// handleWalletCollection routes requests to /api/v1/wallets
func (r *Router) handleWalletCollection(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.walletController.GetWallets(w, req)
	case http.MethodPost:
		r.walletController.CreateWallet(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleWalletResource routes requests to /api/v1/wallets/{walletID}
func (r *Router) handleWalletResource(w http.ResponseWriter, req *http.Request) {
	// Check if this is a balance request (ends with /balance)
	if strings.HasSuffix(req.URL.Path, "/balance") {
		r.walletController.GetWalletBalance(w, req)
		return
	}

	// Route to appropriate wallet resource method
	switch req.Method {
	case http.MethodGet:
		r.walletController.GetWallet(w, req)
	case http.MethodPut:
		r.walletController.UpdateWallet(w, req)
	case http.MethodDelete:
		r.walletController.DeleteWallet(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}