package web

import (
	"net/http"
	"strings"

	"github.com/JingHsiu/accountingApp/internal/accounting/adapter/controller"
)

type Router struct {
	// Specialized wallet controllers
	createWalletController     *controller.CreateWalletController
	queryWalletController      *controller.QueryWalletController
	updateWalletController     *controller.UpdateWalletController
	deleteWalletController     *controller.DeleteWalletController
	getWalletBalanceController *controller.GetWalletBalanceController

	// Specialized transaction controllers
	addExpenseController   *controller.AddExpenseController
	addIncomeController    *controller.AddIncomeController
	queryIncomeController  *controller.QueryIncomeController
	queryExpenseController *controller.QueryExpenseController

	// Category controllers
	categoryController    *controller.CategoryController
	getCategoriesController *controller.GetCategoriesController
}

func NewRouter(
	createWalletController *controller.CreateWalletController,
	queryWalletController *controller.QueryWalletController,
	updateWalletController *controller.UpdateWalletController,
	deleteWalletController *controller.DeleteWalletController,
	getWalletBalanceController *controller.GetWalletBalanceController,
	addExpenseController *controller.AddExpenseController,
	addIncomeController *controller.AddIncomeController,
	queryIncomeController *controller.QueryIncomeController,
	queryExpenseController *controller.QueryExpenseController,
	categoryController *controller.CategoryController,
	getCategoriesController *controller.GetCategoriesController,
) *Router {
	return &Router{
		createWalletController:     createWalletController,
		queryWalletController:      queryWalletController,
		updateWalletController:     updateWalletController,
		deleteWalletController:     deleteWalletController,
		getWalletBalanceController: getWalletBalanceController,
		addExpenseController:       addExpenseController,
		addIncomeController:        addIncomeController,
		queryIncomeController:      queryIncomeController,
		queryExpenseController:     queryExpenseController,
		categoryController:         categoryController,
		getCategoriesController:    getCategoriesController,
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
	mux.HandleFunc("/api/v1/wallets", r.handleWalletCollection)                     // GET (with userID param), POST
	mux.HandleFunc("/api/v1/wallets/", r.handleWalletResource)                      // GET, PUT, DELETE by ID
	mux.HandleFunc("/api/v1/wallets/balance/", r.getWalletBalanceController.GetWalletBalance) // Specialized balance endpoint

	// Category endpoints
	mux.HandleFunc("/api/v1/categories", r.getCategoriesController.GetCategories)              // GET all categories
	mux.HandleFunc("/api/v1/categories/expense", r.getCategoriesController.GetExpenseCategories) // GET expense categories
	mux.HandleFunc("/api/v1/categories/income", r.getCategoriesController.GetIncomeCategories)   // GET income categories

	// Transaction endpoints
	mux.HandleFunc("/api/v1/expenses", r.handleExpenses)
	mux.HandleFunc("/api/v1/incomes", r.handleIncomes)

	return mux
}

// handleWalletCollection routes requests to /api/v1/wallets
func (r *Router) handleWalletCollection(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.queryWalletController.GetWallets(w, req)
	case http.MethodPost:
		r.createWalletController.CreateWallet(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleWalletResource routes requests to /api/v1/wallets/{walletID}
func (r *Router) handleWalletResource(w http.ResponseWriter, req *http.Request) {
	// Check if this is a balance request (ends with /balance)
	if strings.HasSuffix(req.URL.Path, "/balance") {
		r.getWalletBalanceController.GetWalletBalance(w, req)
		return
	}

	// Route to appropriate specialized wallet controller
	switch req.Method {
	case http.MethodGet:
		r.queryWalletController.GetWallet(w, req)
	case http.MethodPut:
		r.updateWalletController.UpdateWallet(w, req)
	case http.MethodDelete:
		r.deleteWalletController.DeleteWallet(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleIncomes routes requests to /api/v1/incomes
func (r *Router) handleIncomes(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.queryIncomeController.GetIncomes(w, req)
	case http.MethodPost:
		r.addIncomeController.AddIncome(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleExpenses routes requests to /api/v1/expenses
func (r *Router) handleExpenses(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.queryExpenseController.GetExpenses(w, req)
	case http.MethodPost:
		r.addExpenseController.AddExpense(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

