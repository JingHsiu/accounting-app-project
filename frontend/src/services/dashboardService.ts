import { apiRequest } from './api'
import type { DashboardData, Money, ApiResponse } from '@/types'

export interface DashboardParams {
  userID: string
  startDate?: string
  endDate?: string
}

export interface MonthlyStats {
  month: string
  income: Money
  expense: Money
  balance: Money
}

export interface CategoryStats {
  categoryID: string
  categoryName: string
  amount: Money
  percentage: number
  transactionCount: number
}

export const dashboardService = {
  // Get dashboard overview data
  getDashboardData: async (params: DashboardParams): Promise<ApiResponse<DashboardData>> => {
    const queryString = new URLSearchParams(
      Object.entries(params)
        .filter(([_, value]) => value !== undefined)
        .map(([key, value]) => [key, String(value)])
    ).toString()
    return apiRequest.get<DashboardData>(`/dashboard?${queryString}`)
  },

  // Get monthly statistics for charts
  getMonthlyStats: async (userID: string, months: number = 12): Promise<ApiResponse<MonthlyStats[]>> => {
    return apiRequest.get<MonthlyStats[]>(`/dashboard/monthly?userID=${userID}&months=${months}`)
  },

  // Get expense breakdown by category
  getExpenseByCategory: async (userID: string, startDate?: string, endDate?: string): Promise<ApiResponse<CategoryStats[]>> => {
    const params = new URLSearchParams({ userID })
    if (startDate) params.append('startDate', startDate)
    if (endDate) params.append('endDate', endDate)
    return apiRequest.get<CategoryStats[]>(`/dashboard/expense-by-category?${params}`)
  },

  // Get income breakdown by category
  getIncomeByCategory: async (userID: string, startDate?: string, endDate?: string): Promise<ApiResponse<CategoryStats[]>> => {
    const params = new URLSearchParams({ userID })
    if (startDate) params.append('startDate', startDate)
    if (endDate) params.append('endDate', endDate)
    return apiRequest.get<CategoryStats[]>(`/dashboard/income-by-category?${params}`)
  },

  // Get net worth trend
  getNetWorthTrend: async (userID: string, months: number = 12): Promise<ApiResponse<{ date: string; netWorth: Money }[]>> => {
    return apiRequest.get(`/dashboard/net-worth?userID=${userID}&months=${months}`)
  },
}