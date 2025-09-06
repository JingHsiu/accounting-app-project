import { apiRequest } from './api'
import type { IncomeRecord, CreateIncomeRequest, UpdateIncomeRequest, IncomeExpenseFilters, ApiResponse } from '@/types'

export const incomeService = {
  // Get all income records with optional filtering
  getIncomes: async (filters?: IncomeExpenseFilters, component = 'IncomeService'): Promise<IncomeRecord[]> => {
    console.group(`🔍 [${component}] Getting income records`)
    
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
      const url = queryString ? `/incomes?${queryString}` : '/incomes'
      
      console.log('📡 Making API request to:', url)
      const response = await apiRequest.get<{data: IncomeRecord[], count: number}>(url)
      
      if (!response.success) {
        console.error('❌ API Error:', response.error)
        console.groupEnd()
        throw new Error(response.error || 'Failed to load income records')
      }

      if (!response.data) {
        console.error('❌ No data in response')
        console.groupEnd()
        throw new Error('No data received from server')
      }

      // Extract income array (handle different response formats)
      let incomeArray: IncomeRecord[] = []
      
      if (Array.isArray(response.data)) {
        // Direct array format
        console.log('📊 Direct array format detected')
        incomeArray = response.data
      } else if (response.data && typeof response.data === 'object') {
        // Expected format: {data: [...], count: number}
        console.log('📊 Nested data structure:', {
          hasDataArray: Array.isArray((response.data as any).data),
          dataLength: (response.data as any).data?.length || 0,
          count: (response.data as any).count
        })
        
        if (Array.isArray((response.data as any).data)) {
          incomeArray = (response.data as any).data
        } else {
          console.error('❌ Unrecognized data structure:', response.data)
          console.groupEnd()
          throw new Error('Invalid data structure: unable to extract income array')
        }
      } else {
        console.error('❌ Invalid response data format:', typeof response.data)
        console.groupEnd()
        throw new Error('Invalid data structure: expected object or array')
      }

      if (!Array.isArray(incomeArray)) {
        console.error('❌ Extracted data is not an array:', typeof incomeArray)
        console.groupEnd()
        throw new Error('Invalid data structure: expected array of income records')
      }

      console.log('✅ Returning income array:', incomeArray.length, 'records')
      console.groupEnd()
      return incomeArray
      
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

    // Handle different response formats defensively
    let income: IncomeRecord | null = null
    
    if (response.data && typeof response.data === 'object') {
      // Check for double-wrapped format: {success: true, data: {data: income}}
      if ((response.data as any).data && typeof (response.data as any).data === 'object') {
        income = (response.data as any).data
      } 
      // Check for direct income format: {success: true, data: income}
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