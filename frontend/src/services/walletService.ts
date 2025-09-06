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
      const response = await apiRequest.get<Wallet[]>(`/wallets?userID=${userID}`)
      
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

      // Step 3: Validate wallet array (apiRequest.get handles unwrapping)
      if (!Array.isArray(response.data)) {
        console.error('❌ Response data is not an array:', typeof response.data)
        console.groupEnd()
        throw new Error('Invalid data structure: expected array of wallets')
      }

      console.log('✅ Returning wallet array:', response.data.length, 'wallets')
      console.groupEnd()
      return response.data
      
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
    const response = await apiRequest.get<Wallet>(`/wallets/${walletID}`)
    
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

    // Validate wallet object structure
    if (!response.data.id || !response.data.name) {
      return {
        success: false,
        error: 'Invalid wallet data structure received'
      }
    }

    return {
      success: true,
      data: response.data
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
    const response = await apiRequest.put<Wallet>(`/wallets/${walletID}`, updates)
    
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

    // Validate wallet object structure
    if (!response.data.id || !response.data.name) {
      return {
        success: false,
        error: 'Invalid wallet data structure received'
      }
    }

    return {
      success: true,
      data: response.data
    }
  },

  // Delete a wallet
  deleteWallet: async (walletID: string): Promise<ApiResponse<void>> => {
    const response = await apiRequest.delete<{message: string}>(`/wallets/${walletID}`)
    
    if (!response.success) {
      return {
        success: false,
        error: response.error || 'Failed to delete wallet'
      }
    }

    let message = 'Wallet deleted successfully'
    if (response.data && (response.data as any).message) {
      message = (response.data as any).message
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