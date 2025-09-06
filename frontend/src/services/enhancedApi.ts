// Enhanced API service with robust 500 error handling and retry logic
import axios, { AxiosResponse } from 'axios'
import type { ApiResponse } from '@/types'
import { apiDebugger } from '@/utils/apiDebug'
import { 
  ErrorClassifier, 
  RetryHandler, 
  serverErrorRetryConfig,
  ErrorReportingService,
  type EnhancedError 
} from '@/utils/errorHandling'

// Enhanced axios instance with circuit breaker logic
export const enhancedApi = axios.create({
  baseURL: '/api/v1',
  timeout: 15000, // Increased timeout for 500 error scenarios
  headers: {
    'Content-Type': 'application/json',
  },
})

// Circuit breaker state
class CircuitBreaker {
  private failures = 0
  private lastFailureTime = 0
  private state: 'CLOSED' | 'OPEN' | 'HALF_OPEN' = 'CLOSED'
  
  constructor(
    private readonly threshold = 5,      // Open after 5 failures
    private readonly timeout = 30000,    // 30 seconds timeout
    private readonly resetTimeout = 60000 // 1 minute reset
  ) {}

  async execute<T>(operation: () => Promise<T>): Promise<T> {
    if (this.state === 'OPEN') {
      if (Date.now() - this.lastFailureTime > this.resetTimeout) {
        this.state = 'HALF_OPEN'
      } else {
        throw new Error('Circuit breaker is OPEN - service temporarily unavailable')
      }
    }

    try {
      const result = await operation()
      
      if (this.state === 'HALF_OPEN') {
        this.reset()
      }
      
      return result
    } catch (error) {
      this.recordFailure()
      throw error
    }
  }

  private recordFailure(): void {
    this.failures++
    this.lastFailureTime = Date.now()
    
    if (this.failures >= this.threshold) {
      this.state = 'OPEN'
    }
  }

  private reset(): void {
    this.failures = 0
    this.state = 'CLOSED'
  }

  getState(): string {
    return this.state
  }
}

const circuitBreaker = new CircuitBreaker()

// Request interceptor with enhanced logging
enhancedApi.interceptors.request.use(
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

// Enhanced response interceptor with 500 error handling
enhancedApi.interceptors.response.use(
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
      component: 'enhancedApi',
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
    
    // Classify the error for better handling
    const classifiedError = ErrorClassifier.classify(error)
    
    // Log error response with classification
    apiDebugger.log({
      url: error.config?.url || 'unknown',
      method: (error.config?.method || 'unknown').toUpperCase(),
      fullUrl: error.request?.responseURL || `${window.location.origin}${error.config?.url}`,
      component: 'enhancedApi',
      success: false,
      error: {
        status: error.response?.status,
        statusText: error.response?.statusText,
        data: error.response?.data,
        message: error.message,
        type: classifiedError.type,
        severity: classifiedError.severity,
        userMessage: classifiedError.userMessage,
        retryable: classifiedError.retryable
      },
      responseTime
    })

    // Report error for monitoring
    ErrorReportingService.reportError(classifiedError, `API ${error.config?.method?.toUpperCase()} ${error.config?.url}`)

    // Handle 401 unauthorized
    if (error.response?.status === 401) {
      localStorage.removeItem('authToken')
      window.location.href = '/login'
    }
    
    return Promise.reject(classifiedError)
  }
)

// Enhanced generic API functions with retry logic
export const enhancedApiRequest = {
  get: async <T>(url: string, options?: { retries?: boolean, component?: string }): Promise<ApiResponse<T>> => {
    const operation = async (): Promise<ApiResponse<T>> => {
      try {
        const response: AxiosResponse = await circuitBreaker.execute(
          () => enhancedApi.get(url)
        )
        
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
        // If it's already a classified error, preserve classification
        if (error.type && error.severity) {
          return {
            success: false,
            error: error.userMessage || error.message,
            errorDetails: error
          }
        }
        
        // Classify new errors
        const classifiedError = ErrorClassifier.classify(error)
        return {
          success: false,
          error: classifiedError.userMessage,
          errorDetails: classifiedError
        }
      }
    }

    // Apply retry logic for retryable errors
    if (options?.retries !== false) {
      try {
        return await RetryHandler.withRetry(operation, serverErrorRetryConfig)
      } catch (error: any) {
        // If retry failed, return the error response
        if (error.type && error.severity) {
          return {
            success: false,
            error: error.userMessage || error.message,
            errorDetails: error
          }
        }
        
        return {
          success: false,
          error: 'Request failed after multiple attempts',
          errorDetails: error
        }
      }
    }

    return operation()
  },

  post: async <T>(url: string, data?: any, options?: { retries?: boolean }): Promise<ApiResponse<T>> => {
    const operation = async (): Promise<ApiResponse<T>> => {
      try {
        const response: AxiosResponse = await circuitBreaker.execute(
          () => enhancedApi.post(url, data)
        )
        
        // Handle backend response format: {success: boolean, data: T, error?: string}
        if (response.data && typeof response.data === 'object' && 'success' in response.data) {
          if (response.data.success) {
            return {
              success: true,
              data: response.data.data || response.data,
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
        // If it's already a classified error, preserve classification
        if (error.type && error.severity) {
          return {
            success: false,
            error: error.userMessage || error.message,
            errorDetails: error
          }
        }
        
        // Classify new errors
        const classifiedError = ErrorClassifier.classify(error)
        return {
          success: false,
          error: classifiedError.userMessage,
          errorDetails: classifiedError
        }
      }
    }

    // Apply retry logic for retryable errors (but not for most POST operations)
    if (options?.retries === true) {
      try {
        return await RetryHandler.withRetry(operation, {
          ...serverErrorRetryConfig,
          maxRetries: 2 // Fewer retries for POST to avoid duplicate operations
        })
      } catch (error: any) {
        if (error.type && error.severity) {
          return {
            success: false,
            error: error.userMessage || error.message,
            errorDetails: error
          }
        }
        
        return {
          success: false,
          error: 'Request failed after multiple attempts',
          errorDetails: error
        }
      }
    }

    return operation()
  },

  put: async <T>(url: string, data?: any, options?: { retries?: boolean }): Promise<ApiResponse<T>> => {
    // Similar implementation to POST
    const operation = async (): Promise<ApiResponse<T>> => {
      try {
        const response: AxiosResponse = await circuitBreaker.execute(
          () => enhancedApi.put(url, data)
        )
        
        if (response.data && typeof response.data === 'object' && 'success' in response.data) {
          if (response.data.success) {
            return {
              success: true,
              data: response.data.data || response.data,
            }
          } else {
            return {
              success: false,
              error: response.data.error || response.data.message || 'Backend returned unsuccessful response',
            }
          }
        }
        
        return {
          success: true,
          data: response.data,
        }
      } catch (error: any) {
        if (error.type && error.severity) {
          return {
            success: false,
            error: error.userMessage || error.message,
            errorDetails: error
          }
        }
        
        const classifiedError = ErrorClassifier.classify(error)
        return {
          success: false,
          error: classifiedError.userMessage,
          errorDetails: classifiedError
        }
      }
    }

    if (options?.retries === true) {
      try {
        return await RetryHandler.withRetry(operation, {
          ...serverErrorRetryConfig,
          maxRetries: 2
        })
      } catch (error: any) {
        if (error.type && error.severity) {
          return {
            success: false,
            error: error.userMessage || error.message,
            errorDetails: error
          }
        }
        
        return {
          success: false,
          error: 'Request failed after multiple attempts',
          errorDetails: error
        }
      }
    }

    return operation()
  },

  delete: async <T>(url: string, options?: { retries?: boolean }): Promise<ApiResponse<T>> => {
    // Similar implementation to GET
    const operation = async (): Promise<ApiResponse<T>> => {
      try {
        const response: AxiosResponse = await circuitBreaker.execute(
          () => enhancedApi.delete(url)
        )
        
        if (response.data && typeof response.data === 'object' && 'success' in response.data) {
          if (response.data.success) {
            return {
              success: true,
              data: response.data.data || response.data,
            }
          } else {
            return {
              success: false,
              error: response.data.error || response.data.message || 'Backend returned unsuccessful response',
            }
          }
        }
        
        return {
          success: true,
          data: response.data,
        }
      } catch (error: any) {
        if (error.type && error.severity) {
          return {
            success: false,
            error: error.userMessage || error.message,
            errorDetails: error
          }
        }
        
        const classifiedError = ErrorClassifier.classify(error)
        return {
          success: false,
          error: classifiedError.userMessage,
          errorDetails: classifiedError
        }
      }
    }

    if (options?.retries !== false) {
      try {
        return await RetryHandler.withRetry(operation, serverErrorRetryConfig)
      } catch (error: any) {
        if (error.type && error.severity) {
          return {
            success: false,
            error: error.userMessage || error.message,
            errorDetails: error
          }
        }
        
        return {
          success: false,
          error: 'Request failed after multiple attempts',
          errorDetails: error
        }
      }
    }

    return operation()
  }
}

// Health check endpoint for monitoring circuit breaker state
export const getApiHealth = () => ({
  circuitBreakerState: circuitBreaker.getState(),
  timestamp: new Date().toISOString()
})