// Enhanced React Query hooks with robust 500 error handling
import { useQuery, useMutation, useQueryClient, UseQueryOptions, UseMutationOptions } from 'react-query'
import { toast } from 'react-hot-toast'
import { 
  ErrorClassifier, 
  ErrorSeverity, 
  ErrorType,
  type EnhancedError 
} from '@/utils/errorHandling'

// Enhanced query options with error handling
export interface EnhancedQueryOptions<TData, TError = EnhancedError> extends Omit<UseQueryOptions<TData, TError>, 'onError'> {
  // Error handling options
  showErrorToast?: boolean
  errorToastDuration?: number
  onError?: (error: EnhancedError) => void
  
  // Retry configuration for server errors
  retryServerErrors?: boolean
  maxServerErrorRetries?: number
  
  // Fallback data when server errors occur
  fallbackData?: TData
  
  // Progressive degradation options
  enableProgressiveDegradation?: boolean
  degradedModeData?: TData
}

// Enhanced mutation options
export interface EnhancedMutationOptions<TData, TVariables, TError = EnhancedError> 
  extends Omit<UseMutationOptions<TData, TError, TVariables>, 'onError'> {
  
  // Error handling options
  showErrorToast?: boolean
  errorToastDuration?: number
  onError?: (error: EnhancedError, variables: TVariables, context?: any) => void
  
  // Success feedback
  showSuccessToast?: boolean
  successMessage?: string | ((data: TData, variables: TVariables) => string)
  
  // Retry for server errors (careful with mutations)
  retryServerErrors?: boolean
  maxServerErrorRetries?: number
}

// Custom error toast function
function showErrorToast(error: EnhancedError, duration?: number) {
  const toastId = `error_${error.timestamp}`
  
  // Don't show duplicate toasts for the same error
  if ((toast as any).isActive?.(toastId)) {
    return
  }

  const toastOptions = {
    id: toastId,
    duration: duration || (error.severity === ErrorSeverity.CRITICAL ? 8000 : 4000),
    style: {
      backgroundColor: error.severity === ErrorSeverity.CRITICAL ? '#fef2f2' : '#fefce8',
      border: `1px solid ${error.severity === ErrorSeverity.CRITICAL ? '#fecaca' : '#fde68a'}`,
      color: error.severity === ErrorSeverity.CRITICAL ? '#991b1b' : '#92400e'
    }
  }

  toast.error(error.userMessage, toastOptions)
}

// Enhanced useQuery hook
export function useEnhancedQuery<TData = unknown>(
  queryKey: any[],
  queryFn: () => Promise<TData>,
  options: EnhancedQueryOptions<TData> = {}
) {
  const {
    showErrorToast = true,
    errorToastDuration,
    onError,
    retryServerErrors = true,
    maxServerErrorRetries = 3,
    fallbackData,
    enableProgressiveDegradation = false,
    degradedModeData,
    ...queryOptions
  } = options

  // Enhanced retry function for server errors
  const retryFunction = (failureCount: number, error: any): boolean => {
    // Always apply custom retry logic first
    if (queryOptions.retry !== undefined) {
      if (typeof queryOptions.retry === 'boolean') {
        return queryOptions.retry
      } else if (typeof queryOptions.retry === 'number') {
        return failureCount < queryOptions.retry
      } else if (typeof queryOptions.retry === 'function') {
        return queryOptions.retry(failureCount, error)
      }
    }

    // Classify error for intelligent retry
    const classifiedError = error?.type ? error : ErrorClassifier.classify(error)
    
    // Retry server errors more aggressively
    if (retryServerErrors && classifiedError.type === ErrorType.SERVER_ERROR) {
      return failureCount < maxServerErrorRetries
    }

    // Retry network and timeout errors
    if (classifiedError.type === ErrorType.NETWORK_ERROR || classifiedError.type === ErrorType.TIMEOUT_ERROR) {
      return failureCount < 2
    }

    // Don't retry client errors (400-499)
    if (classifiedError.type === ErrorType.CLIENT_ERROR || classifiedError.type === ErrorType.VALIDATION_ERROR) {
      return false
    }

    // Default retry once for unknown errors
    return failureCount < 1
  }

  // Enhanced error handler
  const errorHandler = (error: any) => {
    const classifiedError = error?.type ? error : ErrorClassifier.classify(error)
    
    // Show error toast if enabled
    if (showErrorToast) {
      showErrorToast(classifiedError, errorToastDuration)
    }
    
    // Call custom error handler
    onError?.(classifiedError)
    
    // Log error for debugging
    console.error('Enhanced query error:', {
      queryKey,
      error: classifiedError,
      timestamp: new Date().toISOString()
    })
  }

  const result = useQuery(queryKey, queryFn, {
    ...queryOptions,
    retry: retryFunction,
    retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000), // Exponential backoff
    onError: errorHandler,
    
    // Provide fallback data on error if configured
    placeholderData: enableProgressiveDegradation ? (degradedModeData || fallbackData) : queryOptions.placeholderData,
  })

  // Enhanced return with error classification and recovery options
  return {
    ...result,
    // Classify error for better handling in components
    classifiedError: result.error?.type ? result.error as EnhancedError : 
                     result.error ? ErrorClassifier.classify(result.error) : undefined,
    
    // Helper methods for error recovery
    canRetry: () => {
      if (!result.error) return false
      const classified = result.error?.type ? result.error as EnhancedError : ErrorClassifier.classify(result.error)
      return classified.retryable
    },
    
    // Check if in degraded mode (showing fallback data due to errors)
    isDegraded: enableProgressiveDegradation && result.error && !!result.data,
    
    // Manual retry with error classification
    retryQuery: () => result.refetch(),
    
    // Get user-friendly error message
    getUserErrorMessage: (): string | undefined => {
      if (!result.error) return undefined
      const classified = result.error?.type ? result.error as EnhancedError : ErrorClassifier.classify(result.error)
      return classified.userMessage
    }
  }
}

// Enhanced useMutation hook  
export function useEnhancedMutation<TData = unknown, TVariables = unknown>(
  mutationFn: (variables: TVariables) => Promise<TData>,
  options: EnhancedMutationOptions<TData, TVariables> = {}
) {
  const queryClient = useQueryClient()
  
  const {
    showErrorToast = true,
    errorToastDuration,
    showSuccessToast = false,
    successMessage,
    onError,
    retryServerErrors = false, // Conservative default for mutations
    maxServerErrorRetries = 1,
    ...mutationOptions
  } = options

  // Enhanced retry function (conservative for mutations)
  const retryFunction = (failureCount: number, error: any): boolean => {
    if (!retryServerErrors) return false
    
    const classifiedError = error?.type ? error : ErrorClassifier.classify(error)
    
    // Only retry certain server errors and only once for mutations
    if (classifiedError.type === ErrorType.SERVER_ERROR && 
        (classifiedError.statusCode === 500 || classifiedError.statusCode === 503)) {
      return failureCount < maxServerErrorRetries
    }

    return false
  }

  // Enhanced error handler
  const errorHandler = (error: any, variables: TVariables, context?: any) => {
    const classifiedError = error?.type ? error : ErrorClassifier.classify(error)
    
    // Show error toast if enabled
    if (showErrorToast) {
      showErrorToast(classifiedError, errorToastDuration)
    }
    
    // Call custom error handler
    onError?.(classifiedError, variables, context)
    
    // Log error for debugging
    console.error('Enhanced mutation error:', {
      variables,
      error: classifiedError,
      timestamp: new Date().toISOString()
    })
  }

  // Enhanced success handler
  const successHandler = (data: TData, variables: TVariables, context?: any) => {
    // Show success toast if configured
    if (showSuccessToast && successMessage) {
      const message = typeof successMessage === 'function' 
        ? successMessage(data, variables) 
        : successMessage
      
      toast.success(message, {
        duration: 3000,
        style: {
          backgroundColor: '#f0fdf4',
          border: '1px solid #bbf7d0',
          color: '#166534'
        }
      })
    }
    
    // Call original success handler
    mutationOptions.onSuccess?.(data, variables, context)
  }

  const result = useMutation(mutationFn, {
    ...mutationOptions,
    retry: retryFunction,
    retryDelay: (attemptIndex) => Math.min(2000 * 2 ** attemptIndex, 10000), // Conservative backoff
    onError: errorHandler,
    onSuccess: successHandler,
  })

  // Enhanced return with error classification
  return {
    ...result,
    // Classify error for better handling
    classifiedError: result.error?.type ? result.error as EnhancedError : 
                     result.error ? ErrorClassifier.classify(result.error) : undefined,
    
    // Helper methods
    canRetry: () => {
      if (!result.error) return false
      const classified = result.error?.type ? result.error as EnhancedError : ErrorClassifier.classify(result.error)
      return classified.retryable && retryServerErrors
    },
    
    getUserErrorMessage: (): string | undefined => {
      if (!result.error) return undefined
      const classified = result.error?.type ? result.error as EnhancedError : ErrorClassifier.classify(result.error)
      return classified.userMessage
    }
  }
}

// Specialized hooks for common patterns
export function useWalletsQuery(userID: string, options: EnhancedQueryOptions<any[]> = {}) {
  // This would integrate with your existing wallet service
  return useEnhancedQuery(
    ['wallets', userID],
    () => {
      // Your wallet fetching logic here
      throw new Error('Wallet service integration needed')
    },
    {
      fallbackData: [],
      enableProgressiveDegradation: true,
      showErrorToast: true,
      ...options
    }
  )
}

export function useTransactionsQuery(filters: any = {}, options: EnhancedQueryOptions<any[]> = {}) {
  return useEnhancedQuery(
    ['transactions', filters],
    () => {
      // Your transaction fetching logic here  
      throw new Error('Transaction service integration needed')
    },
    {
      fallbackData: [],
      enableProgressiveDegradation: true,
      showErrorToast: true,
      retryServerErrors: true,
      maxServerErrorRetries: 5, // More retries for non-critical data
      ...options
    }
  )
}