import { apiRequest } from './api'
import type { Category, CategoryType, ApiResponse } from '@/types'

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
  getCategories: async (type?: CategoryType): Promise<ApiResponse<Category[]>> => {
    const url = type ? `/categories?type=${type}` : '/categories'
    return apiRequest.get<Category[]>(url)
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