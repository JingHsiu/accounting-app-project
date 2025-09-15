package mapper

import (
	"time"
	"github.com/JingHsiu/accountingApp/internal/accounting/adapter/store"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

// WalletData Wallet的持久化資料結構
type WalletData struct {
	ID              string    `db:"id"`
	UserID          string    `db:"user_id"`
	Name            string    `db:"name"`
	Type            string    `db:"type"`
	Currency        string    `db:"currency"`
	BalanceAmount   int64     `db:"balance_amount"`
	BalanceCurrency string    `db:"balance_currency"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	
	// 子實體資料 (不映射到資料庫欄位，透過關聯表處理)
	IncomeRecords  []IncomeRecordData  `db:"-"`
	ExpenseRecords []ExpenseRecordData `db:"-"`
	Transfers      []TransferData      `db:"-"`
	IsFullyLoaded  bool                `db:"-"`
}

// IncomeRecordData Income Record的持久化資料結構  
type IncomeRecordData struct {
	ID            string    `db:"id"`
	WalletID      string    `db:"wallet_id"`
	SubcategoryID string    `db:"category_id"`
	Amount        int64     `db:"amount"`
	Currency      string    `db:"currency"`
	Description   string    `db:"description"`
	Date          time.Time `db:"date"`
	CreatedAt     time.Time `db:"created_at"`
}

// ExpenseRecordData Expense Record的持久化資料結構
type ExpenseRecordData struct {
	ID            string    `db:"id"`
	WalletID      string    `db:"wallet_id"`
	SubcategoryID string    `db:"category_id"`
	Amount        int64     `db:"amount"`
	Currency      string    `db:"currency"`
	Description   string    `db:"description"`
	Date          time.Time `db:"date"`
	CreatedAt     time.Time `db:"created_at"`
}

// TransferData Transfer的持久化資料結構
type TransferData struct {
	ID              string    `db:"id"`
	FromWalletID    string    `db:"from_wallet_id"`
	ToWalletID      string    `db:"to_wallet_id"`
	Amount          int64     `db:"amount"`
	Currency        string    `db:"currency"`
	Fee             int64     `db:"fee"`
	Description     string    `db:"description"`
	Date            time.Time `db:"date"`
	CreatedAt       time.Time `db:"created_at"`
}

func (wd WalletData) GetID() string {
	return wd.ID
}

// GetID implementations for child entities to satisfy AggregateData interface
func (ird IncomeRecordData) GetID() string {
	return ird.ID
}

func (erd ExpenseRecordData) GetID() string {
	return erd.ID
}

func (td TransferData) GetID() string {
	return td.ID
}

// WalletMapper Wallet聚合的資料轉換器
type WalletMapper struct{}

func NewWalletMapper() *WalletMapper {
	return &WalletMapper{}
}

// ToData 將Wallet Domain Model轉換為WalletData (包含子實體)
func (m *WalletMapper) ToData(wallet *model.Wallet) WalletData {
	walletData := WalletData{
		ID:              wallet.ID,
		UserID:          wallet.UserID,
		Name:            wallet.Name,
		Type:            string(wallet.Type),
		Currency:        wallet.Currency(),
		BalanceAmount:   wallet.Balance.Amount,
		BalanceCurrency: wallet.Balance.Currency,
		CreatedAt:       wallet.CreatedAt,
		UpdatedAt:       wallet.UpdatedAt,
		IsFullyLoaded:   wallet.IsFullyLoaded(),
	}

	// 映射 IncomeRecords
	incomeRecords := wallet.GetIncomeRecords()
	walletData.IncomeRecords = make([]IncomeRecordData, len(incomeRecords))
	for i, income := range incomeRecords {
		walletData.IncomeRecords[i] = IncomeRecordData{
			ID:            income.ID,
			WalletID:      income.WalletID,
			SubcategoryID: income.SubcategoryID,
			Amount:        income.Amount.Amount,
			Currency:      income.Amount.Currency,
			Description:   income.Description,
			Date:          income.Date,
			CreatedAt:     income.CreatedAt,
		}
	}

	// 映射 ExpenseRecords
	expenseRecords := wallet.GetExpenseRecords()
	walletData.ExpenseRecords = make([]ExpenseRecordData, len(expenseRecords))
	for i, expense := range expenseRecords {
		walletData.ExpenseRecords[i] = ExpenseRecordData{
			ID:            expense.ID,
			WalletID:      expense.WalletID,
			SubcategoryID: expense.SubcategoryID,
			Amount:        expense.Amount.Amount,
			Currency:      expense.Amount.Currency,
			Description:   expense.Description,
			Date:          expense.Date,
			CreatedAt:     expense.CreatedAt,
		}
	}

	// 映射 Transfers
	transfers := wallet.GetTransfers()
	walletData.Transfers = make([]TransferData, len(transfers))
	for i, transfer := range transfers {
		walletData.Transfers[i] = TransferData{
			ID:           transfer.ID,
			FromWalletID: transfer.FromWalletID,
			ToWalletID:   transfer.ToWalletID,
			Amount:       transfer.Amount.Amount,
			Currency:     transfer.Amount.Currency,
			Fee:          transfer.Fee.Amount,
			Description:  transfer.Description,
			Date:         transfer.Date,
			CreatedAt:    transfer.CreatedAt,
		}
	}

	return walletData
}

// ToDomain 將WalletData轉換為Wallet Domain Model (包含子實體)
func (m *WalletMapper) ToDomain(data WalletData) (*model.Wallet, error) {
	walletType, err := model.ParseWalletType(data.Type)
	if err != nil {
		return nil, err
	}
	
	balance, err := model.NewMoney(data.BalanceAmount, data.BalanceCurrency)
	if err != nil {
		return nil, err
	}
	
	// 創建基本錢包
	wallet := &model.Wallet{
		ID:        data.ID,
		UserID:    data.UserID,
		Name:      data.Name,
		Type:      walletType,
		Balance:   *balance,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	// 如果有子實體資料，重建完整聚合
	if data.IsFullyLoaded {
		// 重建 IncomeRecords
		for _, incomeData := range data.IncomeRecords {
			amount, err := model.NewMoney(incomeData.Amount, incomeData.Currency)
			if err != nil {
				return nil, err
			}
			
			incomeRecord := model.IncomeRecord{
				ID:            incomeData.ID,
				WalletID:      incomeData.WalletID,
				SubcategoryID: incomeData.SubcategoryID,
				Amount:        *amount,
				Description:   incomeData.Description,
				Date:          incomeData.Date,
				CreatedAt:     incomeData.CreatedAt,
			}
			
			// 透過聚合方法添加到錢包 (這會驗證業務規則)
			err = wallet.LoadIncomeRecord(incomeRecord)
			if err != nil {
				return nil, err
			}
		}
		
		// 重建 ExpenseRecords (類似處理)
		for _, expenseData := range data.ExpenseRecords {
			amount, err := model.NewMoney(expenseData.Amount, expenseData.Currency)
			if err != nil {
				return nil, err
			}
			
			expenseRecord := model.ExpenseRecord{
				ID:            expenseData.ID,
				WalletID:      expenseData.WalletID,
				SubcategoryID: expenseData.SubcategoryID,
				Amount:        *amount,
				Description:   expenseData.Description,
				Date:          expenseData.Date,
				CreatedAt:     expenseData.CreatedAt,
			}
			
			err = wallet.LoadExpenseRecord(expenseRecord)
			if err != nil {
				return nil, err
			}
		}
		
		// 重建 Transfers (類似處理)
		for _, transferData := range data.Transfers {
			amount, err := model.NewMoney(transferData.Amount, transferData.Currency)
			if err != nil {
				return nil, err
			}
			fee, err := model.NewMoney(transferData.Fee, transferData.Currency)
			if err != nil {
				return nil, err
			}
			
			transfer := model.Transfer{
				ID:           transferData.ID,
				FromWalletID: transferData.FromWalletID,
				ToWalletID:   transferData.ToWalletID,
				Amount:       *amount,
				Fee:          *fee,
				Description:  transferData.Description,
				Date:         transferData.Date,
				CreatedAt:    transferData.CreatedAt,
			}
			
			err = wallet.LoadTransfer(transfer)
			if err != nil {
				return nil, err
			}
		}
	}
	
	return wallet, nil
}

// 確保WalletData實現AggregateData介面
var _ store.AggregateData = (*WalletData)(nil)

// 確保子實體資料結構實現AggregateData介面
var _ store.AggregateData = (*IncomeRecordData)(nil)
var _ store.AggregateData = (*ExpenseRecordData)(nil)
var _ store.AggregateData = (*TransferData)(nil)

// 確保WalletMapper實現Mapper介面和AggregateMapper介面
var _ Mapper[*model.Wallet, WalletData] = (*WalletMapper)(nil)
var _ store.AggregateMapper[*model.Wallet, WalletData] = (*WalletMapper)(nil)