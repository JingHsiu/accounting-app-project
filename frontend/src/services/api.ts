import axios from 'axios'
import type { ApiResponse } from '@/types'
import { apiDebugger } from '@/utils/apiDebug'

// Create axios instance with base configuration
export const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Add auth token if available
    const token = localStorage.getItem('authToken')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    // Store request start time for debugging
    config.metadata = { startTime: Date.now() }
    
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor
api.interceptors.response.use(
  (response) => {
    // Calculate response time for debugging
    const responseTime = response.config.metadata 
      ? Date.now() - response.config.metadata.startTime 
      : undefined
    
    // Log successful response
    apiDebugger.log({
      url: response.config.url || 'unknown',
      method: (response.config.method || 'unknown').toUpperCase(),
      fullUrl: response.request?.responseURL || `${window.location.origin}${response.config.url}`,
      component: 'unknown', // Will be overridden by specific services
      success: true,
      data: response.data,
      responseTime
    })
    
    return response
  },
  (error) => {
    // Calculate response time for debugging
    const responseTime = error.config?.metadata 
      ? Date.now() - error.config.metadata.startTime 
      : undefined
    
    // Log error response
    apiDebugger.log({
      url: error.config?.url || 'unknown',
      method: (error.config?.method || 'unknown').toUpperCase(),
      fullUrl: error.request?.responseURL || `${window.location.origin}${error.config?.url}`,
      component: 'unknown', // Will be overridden by specific services
      success: false,
      error: {
        status: error.response?.status,
        statusText: error.response?.statusText,
        data: error.response?.data,
        message: error.message
      },
      responseTime
    })
    
    if (error.response?.status === 401) {
      // Handle unauthorized access
      localStorage.removeItem('authToken')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Generic API functions
export const apiRequest = {
  get: async <T>(url: string): Promise<ApiResponse<T>> => {
    try {
      const response = await api.get(url)
      
      // Handle backend response format: {success: boolean, data: T, error?: string}
      if (response.data && typeof response.data === 'object' && 'success' in response.data) {
        if (response.data.success) {
          return {
            success: true,
            data: response.data.data,
          }
        } else {
          return {
            success: false,
            error: response.data.error || 'Backend returned unsuccessful response',
          }
        }
      }
      
      // Fallback for unexpected response format
      return {
        success: true,
        data: response.data,
      }
    } catch (error: any) {
      // Enhanced error context for better debugging
      const errorContext = {
        url,
        method: 'GET',
        statusCode: error.response?.status,
        isNetworkError: !error.response,
        isProxyError: error.code === 'ECONNREFUSED' || error.message?.includes('ECONNREFUSED'),
        fullUrl: `${window.location.origin}${url}`
      }
      
      console.error('🚨 API GET Error:', errorContext, error)
      
      return {
        success: false,
        error: error.response?.data?.error || error.response?.data?.message || error.message,
      }
    }
  },

  post: async <T>(url: string, data?: any): Promise<ApiResponse<T>> => {
    try {
      const response = await api.post(url, data)
      
      // Handle backend response format: {success: boolean, data: T, error?: string}
      if (response.data && typeof response.data === 'object' && 'success' in response.data) {
        if (response.data.success) {
          return {
            success: true,
            data: response.data.data || response.data, // Some endpoints return data directly in root
          }
        } else {
          return {
            success: false,
            error: response.data.error || response.data.message || 'Backend returned unsuccessful response',
          }
        }
      }
      
      // Fallback for unexpected response format
      return {
        success: true,
        data: response.data,
      }
    } catch (error: any) {
      return {
        success: false,
        error: error.response?.data?.error || error.response?.data?.message || error.message,
      }
    }
  },

  put: async <T>(url: string, data?: any): Promise<ApiResponse<T>> => {
    try {
      const response = await api.put(url, data)
      
      // Handle backend response format: {success: boolean, data: T, error?: string}
      if (response.data && typeof response.data === 'object' && 'success' in response.data) {
        if (response.data.success) {
          return {
            success: true,
            data: response.data.data || response.data, // Some endpoints return data directly in root
          }
        } else {
          return {
            success: false,
            error: response.data.error || response.data.message || 'Backend returned unsuccessful response',
          }
        }
      }
      
      // Fallback for unexpected response format
      return {
        success: true,
        data: response.data,
      }
    } catch (error: any) {
      return {
        success: false,
        error: error.response?.data?.error || error.response?.data?.message || error.message,
      }
    }
  },

  delete: async <T>(url: string): Promise<ApiResponse<T>> => {
    try {
      const response = await api.delete(url)
      
      // Handle backend response format: {success: boolean, data: T, error?: string}
      if (response.data && typeof response.data === 'object' && 'success' in response.data) {
        if (response.data.success) {
          return {
            success: true,
            data: response.data.data || response.data, // Some endpoints return data directly in root
          }
        } else {
          return {
            success: false,
            error: response.data.error || response.data.message || 'Backend returned unsuccessful response',
          }
        }
      }
      
      // Fallback for unexpected response format
      return {
        success: true,
        data: response.data,
      }
    } catch (error: any) {
      return {
        success: false,
        error: error.response?.data?.error || error.response?.data?.message || error.message,
      }
    }
  },
}