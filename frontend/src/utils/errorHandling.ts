// Enhanced error handling utilities for 500 errors and user experience

export enum ErrorType {
  NETWORK_ERROR = 'NETWORK_ERROR',
  SERVER_ERROR = 'SERVER_ERROR',
  CLIENT_ERROR = 'CLIENT_ERROR',
  VALIDATION_ERROR = 'VALIDATION_ERROR',
  TIMEOUT_ERROR = 'TIMEOUT_ERROR',
  UNKNOWN_ERROR = 'UNKNOWN_ERROR'
}

export enum ErrorSeverity {
  LOW = 'LOW',       // Non-critical, app continues working
  MEDIUM = 'MEDIUM', // Feature impaired but workaround available
  HIGH = 'HIGH',     // Core feature broken
  CRITICAL = 'CRITICAL' // App unusable
}

export interface EnhancedError {
  type: ErrorType
  severity: ErrorSeverity
  message: string
  userMessage: string
  retryable: boolean
  statusCode?: number
  originalError?: any
  timestamp: number
  action?: 'retry' | 'refresh' | 'contact_support' | 'ignore'
}

export class ErrorClassifier {
  static classify(error: any): EnhancedError {
    const timestamp = Date.now()
    
    // Network/Connection errors
    if (error.code === 'NETWORK_ERROR' || error.message?.includes('Network Error')) {
      return {
        type: ErrorType.NETWORK_ERROR,
        severity: ErrorSeverity.HIGH,
        message: error.message || 'Network connection failed',
        userMessage: '網路連線不穩定，請檢查您的網路狀態',
        retryable: true,
        originalError: error,
        timestamp,
        action: 'retry'
      }
    }

    // Server errors (500-599)
    if (error.response?.status >= 500 && error.response?.status < 600) {
      const statusCode = error.response.status
      
      // Specific 500 error handling
      if (statusCode === 500) {
        return {
          type: ErrorType.SERVER_ERROR,
          severity: ErrorSeverity.HIGH,
          message: `Internal server error (${statusCode})`,
          userMessage: '伺服器暫時無法處理您的請求，我們正在處理中',
          retryable: true,
          statusCode,
          originalError: error,
          timestamp,
          action: 'retry'
        }
      }

      // Service unavailable (503)
      if (statusCode === 503) {
        return {
          type: ErrorType.SERVER_ERROR,
          severity: ErrorSeverity.CRITICAL,
          message: `Service unavailable (${statusCode})`,
          userMessage: '服務暫時維護中，請稍後再試',
          retryable: true,
          statusCode,
          originalError: error,
          timestamp,
          action: 'retry'
        }
      }

      // Gateway errors (502, 504)
      if (statusCode === 502 || statusCode === 504) {
        return {
          type: ErrorType.SERVER_ERROR,
          severity: ErrorSeverity.HIGH,
          message: `Gateway error (${statusCode})`,
          userMessage: '伺服器連線異常，請稍後重試',
          retryable: true,
          statusCode,
          originalError: error,
          timestamp,
          action: 'retry'
        }
      }

      // Generic server error
      return {
        type: ErrorType.SERVER_ERROR,
        severity: ErrorSeverity.HIGH,
        message: `Server error (${statusCode})`,
        userMessage: '伺服器發生錯誤，請稍後再試',
        retryable: true,
        statusCode,
        originalError: error,
        timestamp,
        action: 'retry'
      }
    }

    // Client errors (400-499)
    if (error.response?.status >= 400 && error.response?.status < 500) {
      const statusCode = error.response.status
      
      if (statusCode === 401) {
        return {
          type: ErrorType.CLIENT_ERROR,
          severity: ErrorSeverity.MEDIUM,
          message: 'Unauthorized access',
          userMessage: '登入已過期，請重新登入',
          retryable: false,
          statusCode,
          originalError: error,
          timestamp,
          action: 'refresh'
        }
      }

      if (statusCode === 403) {
        return {
          type: ErrorType.CLIENT_ERROR,
          severity: ErrorSeverity.MEDIUM,
          message: 'Forbidden access',
          userMessage: '您沒有執行此操作的權限',
          retryable: false,
          statusCode,
          originalError: error,
          timestamp,
          action: 'contact_support'
        }
      }

      if (statusCode === 404) {
        return {
          type: ErrorType.CLIENT_ERROR,
          severity: ErrorSeverity.LOW,
          message: 'Resource not found',
          userMessage: '找不到您要的資料，可能已被刪除',
          retryable: false,
          statusCode,
          originalError: error,
          timestamp,
          action: 'refresh'
        }
      }

      // Validation errors (422)
      if (statusCode === 422) {
        return {
          type: ErrorType.VALIDATION_ERROR,
          severity: ErrorSeverity.LOW,
          message: 'Validation failed',
          userMessage: '輸入的資料格式不正確，請檢查後重試',
          retryable: false,
          statusCode,
          originalError: error,
          timestamp,
          action: 'ignore'
        }
      }

      return {
        type: ErrorType.CLIENT_ERROR,
        severity: ErrorSeverity.MEDIUM,
        message: `Client error (${statusCode})`,
        userMessage: '請求有誤，請檢查輸入內容',
        retryable: false,
        statusCode,
        originalError: error,
        timestamp,
        action: 'ignore'
      }
    }

    // Timeout errors
    if (error.code === 'ECONNABORTED' || error.message?.includes('timeout')) {
      return {
        type: ErrorType.TIMEOUT_ERROR,
        severity: ErrorSeverity.MEDIUM,
        message: 'Request timeout',
        userMessage: '請求超時，可能是網路較慢或伺服器忙碌',
        retryable: true,
        originalError: error,
        timestamp,
        action: 'retry'
      }
    }

    // Unknown errors
    return {
      type: ErrorType.UNKNOWN_ERROR,
      severity: ErrorSeverity.MEDIUM,
      message: error.message || 'Unknown error occurred',
      userMessage: '發生未知錯誤，請稍後再試',
      retryable: true,
      originalError: error,
      timestamp,
      action: 'retry'
    }
  }
}

// Retry configuration for different error types
export interface RetryConfig {
  maxRetries: number
  baseDelay: number
  maxDelay: number
  backoffFactor: number
  retryCondition: (error: EnhancedError) => boolean
}

export const defaultRetryConfig: RetryConfig = {
  maxRetries: 3,
  baseDelay: 1000,    // 1 second
  maxDelay: 30000,    // 30 seconds
  backoffFactor: 2,   // Exponential backoff
  retryCondition: (error: EnhancedError) => error.retryable
}

export const serverErrorRetryConfig: RetryConfig = {
  maxRetries: 5,      // More retries for server errors
  baseDelay: 2000,    // 2 seconds initial delay
  maxDelay: 60000,    // 1 minute max
  backoffFactor: 1.5, // Gentler backoff
  retryCondition: (error: EnhancedError) => 
    error.type === ErrorType.SERVER_ERROR || error.type === ErrorType.TIMEOUT_ERROR
}

export class RetryHandler {
  static async withRetry<T>(
    operation: () => Promise<T>,
    config: RetryConfig = defaultRetryConfig
  ): Promise<T> {
    let lastError: EnhancedError | null = null
    
    for (let attempt = 0; attempt <= config.maxRetries; attempt++) {
      try {
        return await operation()
      } catch (error) {
        const classifiedError = ErrorClassifier.classify(error)
        lastError = classifiedError
        
        // Don't retry if not retryable or max attempts reached
        if (!config.retryCondition(classifiedError) || attempt === config.maxRetries) {
          throw classifiedError
        }
        
        // Calculate delay with exponential backoff and jitter
        const delay = Math.min(
          config.baseDelay * Math.pow(config.backoffFactor, attempt),
          config.maxDelay
        )
        const jitter = delay * 0.1 * Math.random() // 10% jitter
        
        console.warn(`Retrying operation (attempt ${attempt + 1}/${config.maxRetries}) after ${delay + jitter}ms`, classifiedError)
        
        await new Promise(resolve => setTimeout(resolve, delay + jitter))
      }
    }
    
    throw lastError
  }
}

// Error boundary helpers
export interface ErrorBoundaryState {
  hasError: boolean
  error?: EnhancedError
  errorId: string
}

export class ErrorReportingService {
  static reportError(error: EnhancedError, context?: string): void {
    // In production, this would send to error monitoring service
    console.error('Error Report:', {
      errorId: `error_${error.timestamp}_${Math.random().toString(36).substr(2, 9)}`,
      type: error.type,
      severity: error.severity,
      message: error.message,
      statusCode: error.statusCode,
      context,
      timestamp: new Date(error.timestamp).toISOString(),
      userAgent: navigator.userAgent,
      url: window.location.href
    })
  }
}