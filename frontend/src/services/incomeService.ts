import { apiRequest } from './api'
import { incomeTransformer, toBackendFormatWithMapping } from '@/utils/apiTransform'
import type { IncomeRecord, CreateIncomeRequest, UpdateIncomeRequest, IncomeExpenseFilters, ApiResponse } from '@/types'

export const incomeService = {
  // Get all income records with optional filtering  
  getIncomes: async (filters?: IncomeExpenseFilters, component = 'IncomeService'): Promise<IncomeRecord[]> => {
    console.group(`🔍 [${component}] Getting income records`)
    
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
      const url = queryString ? `/incomes?${queryString}` : '/incomes'
      
      console.log('📡 Making API request to:', url)
      const response = await apiRequest.get<{data: any[], count: number}>(url)
      
      // Use shared transformer to handle response
      const transformedIncomes = incomeTransformer.handleResponse(response) as IncomeRecord[]
      
      console.log('✅ Transformed', transformedIncomes.length, 'income records')
      console.groupEnd()
      
      return transformedIncomes
      
    } catch (error) {
      console.error('💥 Exception in getIncomes:', error)
      console.groupEnd()
      
      if (error instanceof Error) {
        throw error
      } else {
        throw new Error('Unknown error occurred while fetching income records')
      }
    }
  },

  // Get a specific income record
  getIncome: async (incomeID: string): Promise<ApiResponse<IncomeRecord>> => {
    const response = await apiRequest.get<{data: IncomeRecord}>(`/incomes/${incomeID}`)
    
    if (!response.success) {
      return {
        success: false,
        error: response.error || 'Failed to get income record'
      }
    }

    if (!response.data) {
      return {
        success: false,
        error: 'No data received from server'
      }
    }

    // Handle different response formats defensively
    let income: IncomeRecord | null = null
    
    if (response.data && typeof response.data === 'object') {
      // Check for nested format: {data: income}
      if ((response.data as any).data && typeof (response.data as any).data === 'object') {
        income = (response.data as any).data
      } 
      // Check for direct income format: income object
      else if ((response.data as any).id && (response.data as any).amount) {
        income = response.data as IncomeRecord
      }
    }

    if (!income) {
      return {
        success: false,
        error: 'Invalid income data structure received'
      }
    }

    return {
      success: true,
      data: income
    }
  },

  // Create a new income record
  createIncome: async (income: CreateIncomeRequest): Promise<{id: string}> => {
    console.log('✨ Creating income record:', income)
    const response = await apiRequest.post<{id: string, message: string}>('/incomes', income)
    console.log('💾 Create income API response:', response)
    
    if (response.success && response.data) {
      return { id: response.data.id }
    }
    
    throw new Error(response.error || 'Failed to create income record')
  },

  // Update income record details  
  updateIncome: async (incomeID: string, updates: UpdateIncomeRequest): Promise<ApiResponse<IncomeRecord>> => {
    const response = await apiRequest.put<{data: IncomeRecord}>(`/incomes/${incomeID}`, updates)
    
    if (!response.success) {
      return {
        success: false,
        error: response.error || 'Failed to update income record'
      }
    }

    if (!response.data) {
      return {
        success: false,
        error: 'No data received from server'
      }
    }

    // Handle different response formats defensively with transformation
    let incomeData: any = null
    
    if (response.data && typeof response.data === 'object') {
      // Check for double-wrapped format: {success: true, data: {data: income}}
      if ((response.data as any).data && typeof (response.data as any).data === 'object') {
        incomeData = (response.data as any).data
      } 
      // Check for direct income format: {success: true, data: income}
      else if ((response.data as any).id && (response.data as any).amount) {
        incomeData = response.data
      }
    }

    if (!incomeData) {
      return {
        success: false,
        error: 'Invalid income data structure received'
      }
    }

    // Transform API snake_case to frontend camelCase
    const transformedIncome: IncomeRecord = {
      id: incomeData.id,
      walletID: incomeData.wallet_id || incomeData.walletID,
      categoryID: incomeData.subcategory_id || incomeData.categoryID,
      amount: incomeData.amount,
      description: incomeData.description,
      date: incomeData.date,
      createdAt: incomeData.created_at || incomeData.createdAt,
      updatedAt: incomeData.updated_at || incomeData.updatedAt
    }

    return {
      success: true,
      data: transformedIncome
    }
  },

  // Delete an income record
  deleteIncome: async (incomeID: string): Promise<ApiResponse<void>> => {
    const response = await apiRequest.delete<{data: {message: string}}>(`/incomes/${incomeID}`)
    
    if (!response.success) {
      return {
        success: false,
        error: response.error || 'Failed to delete income record'
      }
    }

    if (!response.data) {
      return {
        success: false,
        error: 'No data received from server'
      }
    }

    // Handle different response formats defensively
    let message = 'Income record deleted successfully'
    
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