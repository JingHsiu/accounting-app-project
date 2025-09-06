import { apiRequest } from './api'
import type { ExpenseRecord, CreateExpenseRequest, UpdateExpenseRequest, IncomeExpenseFilters, ApiResponse } from '@/types'

export const expenseService = {
  // Get all expense records with optional filtering
  getExpenses: async (filters?: IncomeExpenseFilters, component = 'ExpenseService'): Promise<ExpenseRecord[]> => {
    console.group(`🔍 [${component}] Getting expense records`)
    
    try {
      // Build query parameters
      const queryParams = new URLSearchParams()
      if (filters?.walletID) queryParams.append('walletID', filters.walletID)
      if (filters?.categoryID) queryParams.append('categoryID', filters.categoryID)
      if (filters?.startDate) queryParams.append('startDate', filters.startDate)
      if (filters?.endDate) queryParams.append('endDate', filters.endDate)
      if (filters?.minAmount) queryParams.append('minAmount', filters.minAmount.toString())
      if (filters?.maxAmount) queryParams.append('maxAmount', filters.maxAmount.toString())
      if (filters?.description) queryParams.append('description', filters.description)
      
      const queryString = queryParams.toString()
      const url = queryString ? `/expenses?${queryString}` : '/expenses'
      
      console.log('📡 Making API request to:', url)
      const response = await apiRequest.get<{data: ExpenseRecord[], count: number}>(url)
      
      if (!response.success) {
        console.error('❌ API Error:', response.error)
        console.groupEnd()
        throw new Error(response.error || 'Failed to load expense records')
      }

      if (!response.data) {
        console.error('❌ No data in response')
        console.groupEnd()
        throw new Error('No data received from server')
      }

      // Extract expense array (handle different response formats)
      let expenseArray: ExpenseRecord[] = []
      
      if (Array.isArray(response.data)) {
        // Direct array format
        console.log('📊 Direct array format detected')
        expenseArray = response.data
      } else if (response.data && typeof response.data === 'object') {
        // Expected format: {data: [...], count: number}
        console.log('📊 Nested data structure:', {
          hasDataArray: Array.isArray((response.data as any).data),
          dataLength: (response.data as any).data?.length || 0,
          count: (response.data as any).count
        })
        
        if (Array.isArray((response.data as any).data)) {
          expenseArray = (response.data as any).data
        } else {
          console.error('❌ Unrecognized data structure:', response.data)
          console.groupEnd()
          throw new Error('Invalid data structure: unable to extract expense array')
        }
      } else {
        console.error('❌ Invalid response data format:', typeof response.data)
        console.groupEnd()
        throw new Error('Invalid data structure: expected object or array')
      }

      if (!Array.isArray(expenseArray)) {
        console.error('❌ Extracted data is not an array:', typeof expenseArray)
        console.groupEnd()
        throw new Error('Invalid data structure: expected array of expense records')
      }

      console.log('✅ Returning expense array:', expenseArray.length, 'records')
      console.groupEnd()
      return expenseArray
      
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
    const response = await apiRequest.get<{data: ExpenseRecord}>(`/expenses/${expenseID}`)
    
    if (!response.success) {
      return {
        success: false,
        error: response.error || 'Failed to get expense record'
      }
    }

    if (!response.data) {
      return {
        success: false,
        error: 'No data received from server'
      }
    }

    // Handle different response formats defensively
    let expense: ExpenseRecord | null = null
    
    if (response.data && typeof response.data === 'object') {
      // Check for nested format: {data: expense}
      if ((response.data as any).data && typeof (response.data as any).data === 'object') {
        expense = (response.data as any).data
      } 
      // Check for direct expense format: expense object
      else if ((response.data as any).id && (response.data as any).amount) {
        expense = response.data as ExpenseRecord
      }
    }

    if (!expense) {
      return {
        success: false,
        error: 'Invalid expense data structure received'
      }
    }

    return {
      success: true,
      data: expense
    }
  },

  // Create a new expense record
  createExpense: async (expense: CreateExpenseRequest): Promise<{id: string}> => {
    console.log('✨ Creating expense record:', expense)
    const response = await apiRequest.post<{id: string, message: string}>('/expenses', expense)
    console.log('💾 Create expense API response:', response)
    
    if (response.success && response.data) {
      return { id: response.data.id }
    }
    
    throw new Error(response.error || 'Failed to create expense record')
  },

  // Update expense record details  
  updateExpense: async (expenseID: string, updates: UpdateExpenseRequest): Promise<ApiResponse<ExpenseRecord>> => {
    const response = await apiRequest.put<{data: ExpenseRecord}>(`/expenses/${expenseID}`, updates)
    
    if (!response.success) {
      return {
        success: false,
        error: response.error || 'Failed to update expense record'
      }
    }

    if (!response.data) {
      return {
        success: false,
        error: 'No data received from server'
      }
    }

    // Handle different response formats defensively
    let expense: ExpenseRecord | null = null
    
    if (response.data && typeof response.data === 'object') {
      // Check for double-wrapped format: {success: true, data: {data: expense}}
      if ((response.data as any).data && typeof (response.data as any).data === 'object') {
        expense = (response.data as any).data
      } 
      // Check for direct expense format: {success: true, data: expense}
      else if ((response.data as any).id && (response.data as any).amount) {
        expense = response.data as ExpenseRecord
      }
    }

    if (!expense) {
      return {
        success: false,
        error: 'Invalid expense data structure received'
      }
    }

    return {
      success: true,
      data: expense
    }
  },

  // Delete an expense record
  deleteExpense: async (expenseID: string): Promise<ApiResponse<void>> => {
    const response = await apiRequest.delete<{data: {message: string}}>(`/expenses/${expenseID}`)
    
    if (!response.success) {
      return {
        success: false,
        error: response.error || 'Failed to delete expense record'
      }
    }

    if (!response.data) {
      return {
        success: false,
        error: 'No data received from server'
      }
    }

    // Handle different response formats defensively
    let message = 'Expense record deleted successfully'
    
    if (response.data && typeof response.data === 'object') {
      // Check for double-wrapped format: {success: true, data: {data: {message: string}}}
      if ((response.data as any).data && (response.data as any).data.message) {
        message = (response.data as any).data.message
      }
      // Check for direct message format: {success: true, data: {message: string}}
      else if ((response.data as any).message) {
        message = (response.data as any).message
      }
    }

    return {
      success: true,
      message
    }
  },
}