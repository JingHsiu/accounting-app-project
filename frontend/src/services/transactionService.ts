import { apiRequest } from './api'
import type { ExpenseRecord, IncomeRecord, Transfer, Money, ApiResponse } from '@/types'

export interface AddExpenseRequest {
  walletID: string
  categoryID: string
  amount: Money
  description: string
  date: string
}

export interface AddIncomeRequest {
  walletID: string
  categoryID: string
  amount: Money
  description: string
  date: string
}

export interface ProcessTransferRequest {
  fromWalletID: string
  toWalletID: string
  amount: Money
  description: string
  date: string
}

export interface TransactionQueryParams {
  walletID?: string
  categoryID?: string
  startDate?: string
  endDate?: string
  limit?: number
  offset?: number
}

export const transactionService = {
  // Expense operations
  addExpense: async (expense: AddExpenseRequest): Promise<ApiResponse<ExpenseRecord>> => {
    return apiRequest.post<ExpenseRecord>('/expenses', expense)
  },

  getExpenses: async (params?: TransactionQueryParams): Promise<ApiResponse<ExpenseRecord[]>> => {
    const queryString = new URLSearchParams(params as Record<string, string>).toString()
    return apiRequest.get<ExpenseRecord[]>(`/expenses?${queryString}`)
  },

  updateExpense: async (expenseID: string, updates: Partial<AddExpenseRequest>): Promise<ApiResponse<ExpenseRecord>> => {
    return apiRequest.put<ExpenseRecord>(`/expenses/${expenseID}`, updates)
  },

  deleteExpense: async (expenseID: string): Promise<ApiResponse<void>> => {
    return apiRequest.delete<void>(`/expenses/${expenseID}`)
  },

  // Income operations
  addIncome: async (income: AddIncomeRequest): Promise<ApiResponse<IncomeRecord>> => {
    return apiRequest.post<IncomeRecord>('/incomes', income)
  },

  getIncomes: async (params?: TransactionQueryParams): Promise<ApiResponse<IncomeRecord[]>> => {
    const queryString = new URLSearchParams(params as Record<string, string>).toString()
    return apiRequest.get<IncomeRecord[]>(`/incomes?${queryString}`)
  },

  updateIncome: async (incomeID: string, updates: Partial<AddIncomeRequest>): Promise<ApiResponse<IncomeRecord>> => {
    return apiRequest.put<IncomeRecord>(`/incomes/${incomeID}`, updates)
  },

  deleteIncome: async (incomeID: string): Promise<ApiResponse<void>> => {
    return apiRequest.delete<void>(`/incomes/${incomeID}`)
  },

  // Transfer operations
  processTransfer: async (transfer: ProcessTransferRequest): Promise<ApiResponse<Transfer>> => {
    return apiRequest.post<Transfer>('/transfers', transfer)
  },

  getTransfers: async (params?: TransactionQueryParams): Promise<ApiResponse<Transfer[]>> => {
    const queryString = new URLSearchParams(params as Record<string, string>).toString()
    return apiRequest.get<Transfer[]>(`/transfers?${queryString}`)
  },

  // Combined transaction history
  getTransactionHistory: async (params?: TransactionQueryParams): Promise<ApiResponse<{
    expenses: ExpenseRecord[]
    incomes: IncomeRecord[]
    transfers: Transfer[]
  }>> => {
    const queryString = new URLSearchParams(params as Record<string, string>).toString()
    return apiRequest.get(`/transactions/history?${queryString}`)
  },
}