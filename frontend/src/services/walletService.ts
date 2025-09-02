import { apiRequest } from './api'
import type { Wallet, WalletType, Money, ApiResponse } from '@/types'

export interface CreateWalletRequest {
  name: string
  type: WalletType
  currency: string  // Required by backend (ISO 4217 format)
  user_id: string   // Match backend format
  initialBalance?: number // Optional initial balance
}

export interface UpdateWalletRequest {
  name?: string
  type?: WalletType
}

export const walletService = {
  // Get all wallets for a user
  getWallets: async (userID: string, component = 'WalletService'): Promise<Wallet[]> => {
    console.group(`🔍 [${component}] Getting wallets for userID: ${userID}`)
    
    try {
      // Step 1: Make API call
      console.log('📡 Step 1: Making API request...')
      const response = await apiRequest.get<{data: Wallet[], count: number}>(`/wallets?userID=${userID}`)
      
      // Step 2: Check response success
      if (!response.success) {
        console.error('❌ API Error:', response.error)
        console.groupEnd()
        throw new Error(response.error || 'Failed to load wallets')
      }

      if (!response.data) {
        console.error('❌ No data in response')
        console.groupEnd()
        throw new Error('No data received from server')
      }

      // Step 3: Extract wallet array (apiRequest.get now handles unwrapping)
      let walletArray: Wallet[] = []
      
      if (Array.isArray(response.data)) {
        // Direct array format
        console.log('📊 Step 3: Direct array format detected')
        walletArray = response.data
      } else if (response.data && typeof response.data === 'object') {
        // Expected format: {data: [...], count: number}
        console.log('📊 Step 3: Nested data structure:', {
          hasDataArray: Array.isArray((response.data as any).data),
          dataLength: (response.data as any).data?.length || 0,
          count: (response.data as any).count
        })
        
        if (Array.isArray((response.data as any).data)) {
          walletArray = (response.data as any).data
        } else {
          console.error('❌ Unrecognized data structure:', response.data)
          console.groupEnd()
          throw new Error('Invalid data structure: unable to extract wallet array')
        }
      } else {
        console.error('❌ Invalid response data format:', typeof response.data)
        console.groupEnd()
        throw new Error('Invalid data structure: expected object or array')
      }

      // Step 4: Validate and return
      if (!Array.isArray(walletArray)) {
        console.error('❌ Extracted data is not an array:', typeof walletArray)
        console.groupEnd()
        throw new Error('Invalid data structure: expected array of wallets')
      }

      console.log('✅ Returning wallet array:', walletArray.length, 'wallets')
      console.groupEnd()
      return walletArray
      
    } catch (error) {
      console.error('💥 Exception in getWallets:', error)
      console.groupEnd()
      
      if (error instanceof Error) {
        throw error
      } else {
        throw new Error('Unknown error occurred while fetching wallets')
      }
    }
  },

  // Get a specific wallet with transactions
  getWallet: async (walletID: string): Promise<ApiResponse<Wallet>> => {
    const response = await apiRequest.get<{data: Wallet}>(`/wallets/${walletID}`)
    
    if (!response.success) {
      return {
        success: false,
        error: response.error || 'Failed to get wallet'
      }
    }

    if (!response.data) {
      return {
        success: false,
        error: 'No data received from server'
      }
    }

    // apiRequest.get now handles unwrapping, check if we have wallet data directly
    let wallet: Wallet | null = null
    
    if (response.data && typeof response.data === 'object') {
      // Check for nested format: {data: wallet}
      if ((response.data as any).data && typeof (response.data as any).data === 'object') {
        wallet = (response.data as any).data
      } 
      // Check for direct wallet format: wallet object
      else if ((response.data as any).id && (response.data as any).name) {
        wallet = response.data as Wallet
      }
    }

    if (!wallet) {
      return {
        success: false,
        error: 'Invalid wallet data structure received'
      }
    }

    return {
      success: true,
      data: wallet
    }
  },

  // Create a new wallet
  createWallet: async (wallet: CreateWalletRequest): Promise<{id: string}> => {
    console.log('✨ Creating wallet:', wallet)
    const response = await apiRequest.post<{id: string, message: string}>('/wallets', wallet)
    console.log('💾 Create wallet API response:', response)
    
    // apiRequest.post now handles unwrapping and success/error checking
    if (response.success && response.data) {
      return { id: response.data.id }
    }
    
    // Handle request-level errors
    throw new Error(response.error || 'Failed to create wallet')
  },

  // Update wallet details  
  updateWallet: async (walletID: string, updates: UpdateWalletRequest): Promise<ApiResponse<Wallet>> => {
    const response = await apiRequest.put<{data: Wallet}>(`/wallets/${walletID}`, updates)
    
    if (!response.success) {
      return {
        success: false,
        error: response.error || 'Failed to update wallet'
      }
    }

    if (!response.data) {
      return {
        success: false,
        error: 'No data received from server'
      }
    }

    // Handle different response formats defensively
    let wallet: Wallet | null = null
    
    if (response.data && typeof response.data === 'object') {
      // Check for double-wrapped format: {success: true, data: {data: wallet}}
      if ((response.data as any).data && typeof (response.data as any).data === 'object') {
        wallet = (response.data as any).data
      } 
      // Check for direct wallet format: {success: true, data: wallet}
      else if ((response.data as any).id && (response.data as any).name) {
        wallet = response.data as Wallet
      }
    }

    if (!wallet) {
      return {
        success: false,
        error: 'Invalid wallet data structure received'
      }
    }

    return {
      success: true,
      data: wallet
    }
  },

  // Delete a wallet
  deleteWallet: async (walletID: string): Promise<ApiResponse<void>> => {
    const response = await apiRequest.delete<{data: {message: string}}>(`/wallets/${walletID}`)
    
    if (!response.success) {
      return {
        success: false,
        error: response.error || 'Failed to delete wallet'
      }
    }

    if (!response.data) {
      return {
        success: false,
        error: 'No data received from server'
      }
    }

    // Handle different response formats defensively
    let message = 'Wallet deleted successfully'
    
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

  // Get wallet balance
  getWalletBalance: async (walletID: string): Promise<ApiResponse<Money>> => {
    return apiRequest.get<Money>(`/wallets/${walletID}/balance`)
  },
}