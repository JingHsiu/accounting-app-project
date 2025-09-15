import { apiRequest } from './api'
import { expenseTransformer, toBackendFormatWithMapping, handleApiResponse } from '@/utils/apiTransform'
import type { ExpenseRecord, CreateExpenseRequest, UpdateExpenseRequest, IncomeExpenseFilters, ApiResponse } from '@/types'

export const expenseService = {
  // Get all expense records with optional filtering
  getExpenses: async (filters?: IncomeExpenseFilters, component = 'ExpenseService'): Promise<ExpenseRecord[]> => {
    console.group(`🔍 [${component}] Getting expense records`)
    
    try {
      // Build query parameters (transform to backend format)
      const queryParams = new URLSearchParams()
      if (filters?.walletID) queryParams.append('wallet_id', filters.walletID)
      if (filters?.categoryID) queryParams.append('subcategory_id', filters.categoryID)
      if (filters?.startDate) queryParams.append('start_date', filters.startDate)
      if (filters?.endDate) queryParams.append('end_date', filters.endDate)
      if (filters?.minAmount) queryParams.append('min_amount', filters.minAmount.toString())
      if (filters?.maxAmount) queryParams.append('max_amount', filters.maxAmount.toString())
      if (filters?.description) queryParams.append('description', filters.description)
      
      const queryString = queryParams.toString()
      const url = queryString ? `/expenses?${queryString}` : '/expenses'
      
      console.log('📡 Making API request to:', url)
      const response = await apiRequest.get<{data: any[], count: number}>(url)
      
      // Use shared transformer to handle response
      const transformedExpenses = expenseTransformer.handleResponse(response) as ExpenseRecord[]
      
      console.log('✅ Transformed', transformedExpenses.length, 'expense records')
      console.groupEnd()
      
      return transformedExpenses
      
    } catch (error) {
      console.error('💥 Exception in getExpenses:', error)
      console.groupEnd()
      
      if (error instanceof Error) {
        throw error
      } else {
        throw new Error('Unknown error occurred while fetching expense records')
      }
    }
  },

  // Get a specific expense record
  getExpense: async (expenseID: string): Promise<ApiResponse<ExpenseRecord>> => {
    try {
      const response = await apiRequest.get<{data: any}>(`/expenses/${expenseID}`)
      const transformedExpense = expenseTransformer.handleResponse(response) as ExpenseRecord
      
      return {
        success: true,
        data: transformedExpense
      }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to get expense record'
      }
    }
  },

  // Create a new expense record  
  createExpense: async (expense: CreateExpenseRequest): Promise<{id: string}> => {
    console.log('✨ Creating expense record:', expense)
    
    // expense is already in backend format (wallet_id, subcategory_id)
    const response = await apiRequest.post<{id: string, message: string}>('/expenses', expense)
    console.log('💾 Create expense API response:', response)
    
    if (response.success && response.data) {
      return { id: response.data.id }
    }
    
    throw new Error(response.error || 'Failed to create expense record')
  },

  // Update expense record details  
  updateExpense: async (expenseID: string, updates: UpdateExpenseRequest): Promise<ApiResponse<ExpenseRecord>> => {
    try {
      // Transform frontend camelCase to backend snake_case
      const backendUpdates = toBackendFormatWithMapping(updates)
      const response = await apiRequest.put<{data: any}>(`/expenses/${expenseID}`, backendUpdates)
      
      // Transform response back to frontend format
      const transformedExpense = expenseTransformer.handleResponse(response) as ExpenseRecord
      
      return {
        success: true,
        data: transformedExpense
      }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to update expense record'
      }
    }
  },

  // Delete an expense record
  deleteExpense: async (expenseID: string): Promise<ApiResponse<void>> => {
    try {
      const response = await apiRequest.delete<{data: {message: string}}>(`/expenses/${expenseID}`)
      
      // Extract message from response
      const message = response.data?.data?.message || 
                     response.data?.message || 
                     'Expense record deleted successfully'
      
      return {
        success: true,
        message
      }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to delete expense record'
      }
    }
  },
}