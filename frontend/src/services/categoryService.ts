import { apiRequest } from './api'
import { CategoryType } from '@/types'
import type { Category, ApiResponse } from '@/types'

export interface CreateCategoryRequest {
  name: string
  type: CategoryType
  parentID?: string
  color: string
  icon: string
}

export interface UpdateCategoryRequest {
  name?: string
  color?: string
  icon?: string
  parentID?: string
}

export const categoryService = {
  // Get all categories
  getCategories: async (type?: CategoryType, userID?: string): Promise<Category[]> => {
    // Use demo user ID if not provided
    const userId = userID || "demo-user-123"
    
    if (type === CategoryType.EXPENSE) {
      const response = await apiRequest.get<{data: Category[]}>(`/categories/expense?userID=${userId}`)
      if (response.success && response.data) {
        return Array.isArray(response.data) ? response.data : response.data.data || []
      }
      throw new Error(response.error || 'Failed to load expense categories')
    } else if (type === CategoryType.INCOME) {
      const response = await apiRequest.get<{data: Category[]}>(`/categories/income?userID=${userId}`)
      if (response.success && response.data) {
        return Array.isArray(response.data) ? response.data : response.data.data || []
      }
      throw new Error(response.error || 'Failed to load income categories')
    } else {
      // Get both types and combine
      const [expenseResponse, incomeResponse] = await Promise.all([
        apiRequest.get<{data: Category[]}>(`/categories/expense?userID=${userId}`),
        apiRequest.get<{data: Category[]}>(`/categories/income?userID=${userId}`)
      ])
      
      const expenseCategories = expenseResponse.success && expenseResponse.data ? 
        (Array.isArray(expenseResponse.data) ? expenseResponse.data : expenseResponse.data.data || []) : []
      const incomeCategories = incomeResponse.success && incomeResponse.data ? 
        (Array.isArray(incomeResponse.data) ? incomeResponse.data : incomeResponse.data.data || []) : []
      
      return [...expenseCategories, ...incomeCategories]
    }
  },

  // Get category by ID
  getCategory: async (categoryID: string): Promise<ApiResponse<Category>> => {
    return apiRequest.get<Category>(`/categories/${categoryID}`)
  },

  // Create a new category
  createCategory: async (category: CreateCategoryRequest): Promise<ApiResponse<Category>> => {
    return apiRequest.post<Category>('/categories', category)
  },

  // Update category
  updateCategory: async (categoryID: string, updates: UpdateCategoryRequest): Promise<ApiResponse<Category>> => {
    return apiRequest.put<Category>(`/categories/${categoryID}`, updates)
  },

  // Delete category
  deleteCategory: async (categoryID: string): Promise<ApiResponse<void>> => {
    return apiRequest.delete<void>(`/categories/${categoryID}`)
  },

  // Get category tree (hierarchical structure)
  getCategoryTree: async (type?: CategoryType): Promise<ApiResponse<Category[]>> => {
    const url = type ? `/categories/tree?type=${type}` : '/categories/tree'
    return apiRequest.get<Category[]>(url)
  },
}