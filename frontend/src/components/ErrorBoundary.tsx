// React Error Boundary with comprehensive error handling and recovery
import React, { Component, ReactNode } from 'react'
import { AlertTriangle, RefreshCw, ArrowLeft, Bug } from 'lucide-react'
import { Card, CardHeader, CardTitle, CardContent, Button } from '@/components/ui'
import { 
  ErrorClassifier, 
  ErrorReportingService, 
  type EnhancedError, 
  ErrorSeverity 
} from '@/utils/errorHandling'

interface ErrorBoundaryProps {
  children: ReactNode
  fallback?: ReactNode
  level?: 'page' | 'section' | 'component'
  onError?: (error: EnhancedError, errorInfo: React.ErrorInfo) => void
}

interface ErrorBoundaryState {
  hasError: boolean
  error?: EnhancedError
  errorId: string
  retryCount: number
}

export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  private maxRetries = 3

  constructor(props: ErrorBoundaryProps) {
    super(props)
    this.state = {
      hasError: false,
      errorId: '',
      retryCount: 0
    }
  }

  static getDerivedStateFromError(error: Error): Partial<ErrorBoundaryState> {
    // Classify the error
    const classifiedError = ErrorClassifier.classify(error)
    
    return {
      hasError: true,
      error: classifiedError,
      errorId: `boundary_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    }
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    const classifiedError = this.state.error || ErrorClassifier.classify(error)
    
    // Report error
    ErrorReportingService.reportError(
      classifiedError, 
      `ErrorBoundary (${this.props.level || 'unknown'})`
    )
    
    // Call custom error handler
    this.props.onError?.(classifiedError, errorInfo)
    
    console.error('ErrorBoundary caught error:', {
      error: classifiedError,
      errorInfo,
      componentStack: errorInfo.componentStack,
      errorBoundary: this.props.level || 'unknown'
    })
  }

  handleRetry = () => {
    if (this.state.retryCount < this.maxRetries) {
      this.setState(prevState => ({
        hasError: false,
        error: undefined,
        errorId: '',
        retryCount: prevState.retryCount + 1
      }))
    }
  }

  handleRefresh = () => {
    window.location.reload()
  }

  handleGoBack = () => {
    window.history.back()
  }

  render() {
    if (!this.state.hasError) {
      return this.props.children
    }

    // Use custom fallback if provided
    if (this.props.fallback) {
      return this.props.fallback
    }

    const { error } = this.state
    const level = this.props.level || 'component'
    
    if (!error) {
      return (
        <div className="flex items-center justify-center p-8">
          <p className="text-red-600">發生未知錯誤</p>
        </div>
      )
    }

    // Different UI based on error severity and boundary level
    const showFullError = error.severity === ErrorSeverity.CRITICAL || level === 'page'
    const canRetry = this.state.retryCount < this.maxRetries && error.retryable
    
    return (
      <Card className={`${showFullError ? 'min-h-96' : 'min-h-32'} ${this.getErrorSeverityClass(error.severity)}`}>
        <CardHeader className="pb-3">
          <CardTitle className="flex items-center gap-3">
            <div className={`p-2 rounded-lg ${this.getIconClass(error.severity)}`}>
              <AlertTriangle className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-lg font-semibold">
                {this.getErrorTitle(error.severity, level)}
              </h3>
              {showFullError && (
                <p className="text-sm text-neutral-500 mt-1 font-normal">
                  錯誤 ID: {this.state.errorId}
                </p>
              )}
            </div>
          </CardTitle>
        </CardHeader>
        
        <CardContent>
          {/* User-friendly error message */}
          <div className="mb-4">
            <p className="text-neutral-700 mb-2">
              {error.userMessage}
            </p>
            
            {error.severity === ErrorSeverity.CRITICAL && (
              <p className="text-sm text-neutral-600 bg-neutral-50 p-3 rounded-lg border border-neutral-200">
                如果問題持續發生，請聯繫客服並提供錯誤 ID: <code className="font-mono text-xs bg-neutral-200 px-1 py-0.5 rounded">{this.state.errorId}</code>
              </p>
            )}
          </div>

          {/* Action buttons based on error type and severity */}
          <div className="flex flex-wrap gap-2">
            {canRetry && (
              <Button
                variant="primary"
                size="sm"
                onClick={this.handleRetry}
                className="flex items-center gap-2"
              >
                <RefreshCw className="w-4 h-4" />
                重試 ({this.maxRetries - this.state.retryCount} 次機會)
              </Button>
            )}
            
            {!canRetry && error.retryable && (
              <Button
                variant="secondary"
                size="sm"
                onClick={this.handleRefresh}
                className="flex items-center gap-2"
              >
                <RefreshCw className="w-4 h-4" />
                重新整理頁面
              </Button>
            )}
            
            {level === 'page' && (
              <Button
                variant="outline"
                size="sm"
                onClick={this.handleGoBack}
                className="flex items-center gap-2"
              >
                <ArrowLeft className="w-4 h-4" />
                返回上一頁
              </Button>
            )}
          </div>

          {/* Technical details for development */}
          {process.env.NODE_ENV === 'development' && (
            <details className="mt-4">
              <summary className="cursor-pointer text-sm text-neutral-500 hover:text-neutral-700 flex items-center gap-2">
                <Bug className="w-4 h-4" />
                開發者資訊
              </summary>
              <div className="mt-2 p-3 bg-neutral-100 rounded text-xs font-mono overflow-auto">
                <div><strong>Type:</strong> {error.type}</div>
                <div><strong>Severity:</strong> {error.severity}</div>
                <div><strong>Status:</strong> {error.statusCode || 'N/A'}</div>
                <div><strong>Retryable:</strong> {error.retryable ? 'Yes' : 'No'}</div>
                <div><strong>Message:</strong> {error.message}</div>
                <div><strong>Timestamp:</strong> {new Date(error.timestamp).toLocaleString()}</div>
                {error.originalError && (
                  <div><strong>Original:</strong> {JSON.stringify(error.originalError, null, 2)}</div>
                )}
              </div>
            </details>
          )}
        </CardContent>
      </Card>
    )
  }

  private getErrorTitle(severity: ErrorSeverity, level: string): string {
    switch (severity) {
      case ErrorSeverity.LOW:
        return level === 'page' ? '頁面載入問題' : '功能暫時無法使用'
      case ErrorSeverity.MEDIUM:
        return level === 'page' ? '頁面發生錯誤' : '功能執行錯誤'
      case ErrorSeverity.HIGH:
        return level === 'page' ? '頁面無法載入' : '核心功能異常'
      case ErrorSeverity.CRITICAL:
        return '系統發生嚴重錯誤'
      default:
        return '發生未知錯誤'
    }
  }

  private getErrorSeverityClass(severity: ErrorSeverity): string {
    switch (severity) {
      case ErrorSeverity.LOW:
        return 'border-l-4 border-l-yellow-400'
      case ErrorSeverity.MEDIUM:
        return 'border-l-4 border-l-orange-400'
      case ErrorSeverity.HIGH:
        return 'border-l-4 border-l-red-400'
      case ErrorSeverity.CRITICAL:
        return 'border-l-4 border-l-red-600'
      default:
        return 'border-l-4 border-l-gray-400'
    }
  }

  private getIconClass(severity: ErrorSeverity): string {
    switch (severity) {
      case ErrorSeverity.LOW:
        return 'bg-yellow-100 text-yellow-600'
      case ErrorSeverity.MEDIUM:
        return 'bg-orange-100 text-orange-600'
      case ErrorSeverity.HIGH:
        return 'bg-red-100 text-red-600'
      case ErrorSeverity.CRITICAL:
        return 'bg-red-200 text-red-700'
      default:
        return 'bg-gray-100 text-gray-600'
    }
  }
}

// Specialized error boundaries for different contexts
export const PageErrorBoundary: React.FC<{ children: ReactNode }> = ({ children }) => (
  <ErrorBoundary level="page">{children}</ErrorBoundary>
)

export const SectionErrorBoundary: React.FC<{ children: ReactNode }> = ({ children }) => (
  <ErrorBoundary level="section">{children}</ErrorBoundary>
)

export const ComponentErrorBoundary: React.FC<{ children: ReactNode; fallback?: ReactNode }> = ({ 
  children, 
  fallback 
}) => (
  <ErrorBoundary level="component" fallback={fallback}>{children}</ErrorBoundary>
)