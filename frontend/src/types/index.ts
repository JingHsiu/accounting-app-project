// Domain Types matching the backend Go models
export interface Money {
  currency: string;
  amount: number;
}

export interface Wallet {
  id: string;
  user_id: string;    // Match backend format
  name: string;
  type: WalletType;
  balance: Money;
  created_at: string; // Match backend format  
  updated_at: string; // Match backend format
  currency: string;   // Additional field from backend
  is_fully_loaded?: boolean; // Additional field from backend
}

export enum WalletType {
  CASH = "CASH",
  BANK = "BANK", 
  CREDIT = "CREDIT",
  INVESTMENT = "INVESTMENT"
}

export interface Category {
  id: string;
  name: string;
  type: CategoryType;
  parentID?: string;
  color: string;
  icon: string;
  createdAt: string;
}

export enum CategoryType {
  INCOME = "INCOME",
  EXPENSE = "EXPENSE"
}

export interface ExpenseRecord {
  id: string;
  walletID: string;
  categoryID: string;
  amount: Money;
  description: string;
  date: string;
  createdAt: string;
  updatedAt: string;
}

export interface IncomeRecord {
  id: string;
  walletID: string;
  categoryID: string;
  amount: Money;
  description: string;
  date: string;
  createdAt: string;
  updatedAt: string;
}

export interface Transfer {
  id: string;
  fromWalletID: string;
  toWalletID: string;
  amount: Money;
  description: string;
  date: string;
  createdAt: string;
}

// API Response Types
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}

// UI State Types
export interface DashboardData {
  totalBalance: Money;
  monthlyIncome: Money;
  monthlyExpense: Money;
  wallets: Wallet[];
  recentTransactions: (ExpenseRecord | IncomeRecord | Transfer)[];
}

export interface TransactionFilters {
  walletID?: string;
  categoryID?: string;
  startDate?: string;
  endDate?: string;
  type?: 'income' | 'expense' | 'transfer';
}

// Note: WalletType and CategoryType are already exported above as enums