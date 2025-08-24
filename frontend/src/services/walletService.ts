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
  getWallets: async (userID: string): Promise<ApiResponse<Wallet[]>> => {
    try {
      const response = await fetch(`/api/v1/wallets?userID=${userID}`)
      const data = await response.json()
      
      if (data.success && data.data?.data) {
        return {
          success: true,
          data: data.data.data // Extract wallet array from nested response
        }
      }
      
      return {
        success: false,
        error: data.error || 'Failed to load wallets'
      }
    } catch (error: any) {
      return {
        success: false,
        error: error.message
      }
    }
  },

  // Get a specific wallet with transactions
  getWallet: async (walletID: string): Promise<ApiResponse<Wallet>> => {
    const response = await apiRequest.get<{data: Wallet}>(`/wallets/${walletID}`)
    if (response.success && response.data) {
      return {
        success: true,
        data: response.data.data // Extract the nested data
      }
    }
    return response as unknown as ApiResponse<Wallet>
  },

  // Create a new wallet
  createWallet: async (wallet: CreateWalletRequest): Promise<ApiResponse<{id: string}>> => {
    // Backend returns {id, success, message} format for create
    return apiRequest.post<{id: string}>('/wallets', wallet)
  },

  // Update wallet details  
  updateWallet: async (walletID: string, updates: UpdateWalletRequest): Promise<ApiResponse<Wallet>> => {
    const response = await apiRequest.put<{data: Wallet}>(`/wallets/${walletID}`, updates)
    if (response.success && response.data) {
      return {
        success: true,
        data: response.data.data // Extract the nested data
      }
    }
    return response as unknown as ApiResponse<Wallet>
  },

  // Delete a wallet
  deleteWallet: async (walletID: string): Promise<ApiResponse<void>> => {
    // Backend returns {success: true, data: {message: string}} format for delete
    const response = await apiRequest.delete<{data: {message: string}}>(`/wallets/${walletID}`)
    if (response.success && response.data) {
      return {
        success: true,
        message: response.data.data.message // Extract message from nested structure
      }
    }
    return response as unknown as ApiResponse<void>
  },

  // Get wallet balance
  getWalletBalance: async (walletID: string): Promise<ApiResponse<Money>> => {
    return apiRequest.get<Money>(`/wallets/${walletID}/balance`)
  },
}