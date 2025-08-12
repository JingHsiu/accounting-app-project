import { apiRequest } from './api'
import type { Wallet, WalletType, Money, ApiResponse } from '@/types'

export interface CreateWalletRequest {
  name: string
  type: WalletType
  userID: string
}

export interface UpdateWalletRequest {
  name?: string
  type?: WalletType
}

export const walletService = {
  // Get all wallets for a user
  getWallets: async (userID: string): Promise<ApiResponse<Wallet[]>> => {
    return apiRequest.get<Wallet[]>(`/wallets?userID=${userID}`)
  },

  // Get a specific wallet with transactions
  getWallet: async (walletID: string): Promise<ApiResponse<Wallet>> => {
    return apiRequest.get<Wallet>(`/wallets/${walletID}`)
  },

  // Create a new wallet
  createWallet: async (wallet: CreateWalletRequest): Promise<ApiResponse<Wallet>> => {
    return apiRequest.post<Wallet>('/wallets', wallet)
  },

  // Update wallet details
  updateWallet: async (walletID: string, updates: UpdateWalletRequest): Promise<ApiResponse<Wallet>> => {
    return apiRequest.put<Wallet>(`/wallets/${walletID}`, updates)
  },

  // Delete a wallet
  deleteWallet: async (walletID: string): Promise<ApiResponse<void>> => {
    return apiRequest.delete<void>(`/wallets/${walletID}`)
  },

  // Get wallet balance
  getWalletBalance: async (walletID: string): Promise<ApiResponse<Money>> => {
    return apiRequest.get<Money>(`/wallets/${walletID}/balance`)
  },
}