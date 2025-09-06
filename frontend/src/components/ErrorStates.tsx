// Reusable error state components for consistent UX during 500 errors
import React from 'react'
import { 
  AlertTriangle, 
  RefreshCw, 
  WifiOff, 
  Server, 
  Clock,
  ArrowLeft,
  Settings,
  MessageCircle
} from 'lucide-react'
import { Card, CardContent, Button } from '@/components/ui'
import { ErrorType, ErrorSeverity, type EnhancedError } from '@/utils/errorHandling'

interface BaseErrorStateProps {
  error: EnhancedError
  onRetry?: () => void
  onRefresh?: () => void
  onGoBack?: () => void
  showTechnicalDetails?: boolean
  className?: string
}

// Generic error state component
export const ErrorState: React.FC<BaseErrorStateProps> = ({
  error,
  onRetry,
  onRefresh,
  onGoBack,
  showTechnicalDetails = false,
  className = ''
}) => {
  const getErrorIcon = () => {
    switch (error.type) {
      case ErrorType.NETWORK_ERROR:
        return <WifiOff className="w-12 h-12" />
      case ErrorType.SERVER_ERROR:
        return <Server className="w-12 h-12" />
      case ErrorType.TIMEOUT_ERROR:
        return <Clock className="w-12 h-12" />
      default:
        return <AlertTriangle className="w-12 h-12" />
    }
  }

  const getErrorColor = () => {
    switch (error.severity) {
      case ErrorSeverity.LOW:
        return 'text-yellow-600'
      case ErrorSeverity.MEDIUM:
        return 'text-orange-600'
      case ErrorSeverity.HIGH:
        return 'text-red-600'
      case ErrorSeverity.CRITICAL:
        return 'text-red-700'
      default:
        return 'text-gray-600'
    }
  }

  const getBgColor = () => {
    switch (error.severity) {
      case ErrorSeverity.LOW:
        return 'bg-yellow-50 border-yellow-200'
      case ErrorSeverity.MEDIUM:
        return 'bg-orange-50 border-orange-200'
      case ErrorSeverity.HIGH:
        return 'bg-red-50 border-red-200'
      case ErrorSeverity.CRITICAL:
        return 'bg-red-100 border-red-300'
      default:
        return 'bg-gray-50 border-gray-200'
    }
  }

  return (
    <Card className={`${getBgColor()} border-2 ${className}`}>
      <CardContent className="text-center py-8">
        <div className={`${getErrorColor()} mb-4 flex justify-center`}>
          {getErrorIcon()}
        </div>
        
        <h3 className="text-lg font-semibold text-gray-800 mb-2">
          {getErrorTitle(error)}
        </h3>
        
        <p className="text-gray-700 mb-6 max-w-md mx-auto">
          {error.userMessage}
        </p>

        <div className="flex flex-wrap justify-center gap-3">
          {error.retryable && onRetry && (
            <Button
              variant="primary"
              size="sm"
              onClick={onRetry}
              className="flex items-center gap-2"
            >
              <RefreshCw className="w-4 h-4" />
              重試
            </Button>
          )}
          
          {onRefresh && (
            <Button
              variant="secondary"
              size="sm"
              onClick={onRefresh}
              className="flex items-center gap-2"
            >
              <RefreshCw className="w-4 h-4" />
              重新整理
            </Button>
          )}
          
          {onGoBack && (
            <Button
              variant="outline"
              size="sm"
              onClick={onGoBack}
              className="flex items-center gap-2"
            >
              <ArrowLeft className="w-4 h-4" />
              返回
            </Button>
          )}
        </div>

        {/* Technical details for development or support */}
        {showTechnicalDetails && (
          <details className="mt-6 text-left">
            <summary className="cursor-pointer text-sm text-gray-500 hover:text-gray-700 text-center">
              技術資訊
            </summary>
            <div className="mt-3 p-3 bg-gray-100 rounded text-xs font-mono text-left overflow-auto">
              <div><strong>錯誤類型:</strong> {error.type}</div>
              <div><strong>嚴重程度:</strong> {error.severity}</div>
              <div><strong>狀態碼:</strong> {error.statusCode || 'N/A'}</div>
              <div><strong>可重試:</strong> {error.retryable ? '是' : '否'}</div>
              <div><strong>時間:</strong> {new Date(error.timestamp).toLocaleString()}</div>
              <div><strong>詳細訊息:</strong> {error.message}</div>
            </div>
          </details>
        )}
      </CardContent>
    </Card>
  )
}

// Specialized error states for different scenarios
export const ServerErrorState: React.FC<Omit<BaseErrorStateProps, 'error'> & { 
  statusCode?: number,
  customMessage?: string 
}> = ({ 
  statusCode = 500, 
  customMessage,
  onRetry, 
  onRefresh,
  ...props 
}) => {
  const error: EnhancedError = {
    type: ErrorType.SERVER_ERROR,
    severity: statusCode >= 500 ? ErrorSeverity.HIGH : ErrorSeverity.MEDIUM,
    message: `Server error (${statusCode})`,
    userMessage: customMessage || getServerErrorMessage(statusCode),
    retryable: true,
    statusCode,
    timestamp: Date.now(),
    action: 'retry'
  }

  return (
    <ErrorState 
      error={error}
      onRetry={onRetry}
      onRefresh={onRefresh}
      {...props}
    />
  )
}

export const NetworkErrorState: React.FC<Omit<BaseErrorStateProps, 'error'>> = ({ 
  onRetry, 
  onRefresh,
  ...props 
}) => {
  const error: EnhancedError = {
    type: ErrorType.NETWORK_ERROR,
    severity: ErrorSeverity.HIGH,
    message: 'Network connection failed',
    userMessage: '網路連線出現問題，請檢查您的網路狀態後重試',
    retryable: true,
    timestamp: Date.now(),
    action: 'retry'
  }

  return (
    <ErrorState 
      error={error}
      onRetry={onRetry}
      onRefresh={onRefresh}
      {...props}
    />
  )
}

export const TimeoutErrorState: React.FC<Omit<BaseErrorStateProps, 'error'>> = ({ 
  onRetry, 
  onRefresh,
  ...props 
}) => {
  const error: EnhancedError = {
    type: ErrorType.TIMEOUT_ERROR,
    severity: ErrorSeverity.MEDIUM,
    message: 'Request timeout',
    userMessage: '請求超時，可能是網路較慢或伺服器忙碌，請稍後重試',
    retryable: true,
    timestamp: Date.now(),
    action: 'retry'
  }

  return (
    <ErrorState 
      error={error}
      onRetry={onRetry}
      onRefresh={onRefresh}
      {...props}
    />
  )
}

// Inline error message for smaller spaces
export const InlineError: React.FC<{
  error: EnhancedError
  onRetry?: () => void
  className?: string
}> = ({ error, onRetry, className = '' }) => {
  return (
    <div className={`flex items-center justify-between p-3 bg-red-50 border border-red-200 rounded-lg ${className}`}>
      <div className="flex items-center gap-2">
        <AlertTriangle className="w-4 h-4 text-red-600" />
        <span className="text-sm text-red-700">{error.userMessage}</span>
      </div>
      {error.retryable && onRetry && (
        <Button
          variant="ghost"
          size="sm"
          onClick={onRetry}
          className="text-red-600 hover:text-red-700 hover:bg-red-100"
        >
          重試
        </Button>
      )}
    </div>
  )
}

// Loading state with error fallback
export const LoadingWithError: React.FC<{
  isLoading: boolean
  error?: EnhancedError
  onRetry?: () => void
  children: React.ReactNode
  loadingComponent?: React.ReactNode
  className?: string
}> = ({ isLoading, error, onRetry, children, loadingComponent, className = '' }) => {
  if (error) {
    return (
      <div className={className}>
        <ErrorState error={error} onRetry={onRetry} />
      </div>
    )
  }

  if (isLoading) {
    return (
      <div className={className}>
        {loadingComponent || (
          <div className="flex items-center justify-center py-8">
            <div className="w-6 h-6 border-2 border-primary-300 border-t-primary-600 rounded-full animate-spin mr-3" />
            <span className="text-gray-600">載入中...</span>
          </div>
        )}
      </div>
    )
  }

  return <>{children}</>
}

// Helper functions
function getErrorTitle(error: EnhancedError): string {
  switch (error.type) {
    case ErrorType.NETWORK_ERROR:
      return '網路連線問題'
    case ErrorType.SERVER_ERROR:
      if (error.statusCode === 500) return '伺服器內部錯誤'
      if (error.statusCode === 502) return '網關錯誤'
      if (error.statusCode === 503) return '服務暫時無法使用'
      if (error.statusCode === 504) return '網關超時'
      return '伺服器錯誤'
    case ErrorType.TIMEOUT_ERROR:
      return '請求超時'
    case ErrorType.CLIENT_ERROR:
      return '請求錯誤'
    default:
      return '發生錯誤'
  }
}

function getServerErrorMessage(statusCode: number): string {
  switch (statusCode) {
    case 500:
      return '伺服器暫時無法處理您的請求，我們正在處理中，請稍後再試'
    case 502:
      return '伺服器連線異常，請稍後重試'
    case 503:
      return '服務暫時維護中，請稍後再試'
    case 504:
      return '伺服器回應超時，請稍後重試'
    default:
      return '伺服器發生錯誤，請稍後再試'
  }
}