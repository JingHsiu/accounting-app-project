package web

import (
	"net/http"

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

	// Wallet endpoints
	mux.HandleFunc("/api/v1/wallets", r.walletController.CreateWallet)
	mux.HandleFunc("/api/v1/wallets/", r.walletController.GetWalletBalance)

	// Category endpoints
	mux.HandleFunc("/api/v1/categories/expense", r.categoryController.CreateExpenseCategory)
	mux.HandleFunc("/api/v1/categories/income", r.categoryController.CreateIncomeCategory)

	// Transaction endpoints
	mux.HandleFunc("/api/v1/expenses", r.walletController.AddExpense)
	mux.HandleFunc("/api/v1/incomes", r.walletController.AddIncome)

	return mux
}