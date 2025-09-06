// Enhanced Wallets page with comprehensive 500 error handling
import React, { useState } from 'react'
import { useQueryClient } from 'react-query'
import { 
  Wallet as WalletIcon, 
  Plus, 
  Edit, 
  Trash2, 
  CreditCard,
  Banknote,
  TrendingUp,
  AlertCircle
} from 'lucide-react'
import { toast } from 'react-hot-toast'

// Enhanced imports
import { Card, CardContent, Button, Modal, Input, Select } from '@/components/ui'
import { ErrorBoundary, SectionErrorBoundary } from '@/components/ErrorBoundary'
import { ServerErrorState, LoadingWithError, InlineError } from '@/components/ErrorStates'
import { useEnhancedQuery, useEnhancedMutation } from '@/hooks/useEnhancedQuery'

// Enhanced service import (would replace existing walletService)
import { enhancedApiRequest } from '@/services/enhancedApi'
import { formatMoney, getWalletTypeDisplayName } from '@/utils/format'
import type { CreateWalletRequest } from '@/services/walletService'
import { WalletType, type Wallet } from '@/types'

const DEMO_USER_ID = "demo-user-123"

// Enhanced wallet service functions
const enhancedWalletService = {
  getWallets: async (userID: string): Promise<Wallet[]> => {
    const response = await enhancedApiRequest.get<Wallet[]>(`/wallets?userID=${userID}`, {
      retries: true,
      component: 'WalletsEnhanced'
    })
    
    if (!response.success) {
      throw new Error(response.error || 'Failed to load wallets')
    }
    
    return response.data || []
  },

  createWallet: async (wallet: CreateWalletRequest): Promise<{id: string}> => {
    const response = await enhancedApiRequest.post<{id: string}>('/wallets', wallet, {
      retries: true // Allow retries for wallet creation
    })
    
    if (!response.success) {
      throw new Error(response.error || 'Failed to create wallet')
    }
    
    return response.data!
  },

  deleteWallet: async (walletID: string): Promise<void> => {
    const response = await enhancedApiRequest.delete<void>(`/wallets/${walletID}`)
    
    if (!response.success) {
      throw new Error(response.error || 'Failed to delete wallet')
    }
  }
}

const WalletsEnhanced: React.FC = () => {
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [formData, setFormData] = useState({
    name: '',
    type: WalletType.CASH,
    currency: 'TWD',
    initialBalance: 0
  })

  const queryClient = useQueryClient()

  // Enhanced wallet query with comprehensive error handling
  const {
    data: wallets = [],
    isLoading,
    classifiedError,
    isDegraded,
    canRetry,
    retryQuery,
    getUserErrorMessage
  } = useEnhancedQuery(
    ['wallets', DEMO_USER_ID],
    () => enhancedWalletService.getWallets(DEMO_USER_ID),
    {
      showErrorToast: true,
      retryServerErrors: true,
      maxServerErrorRetries: 5,
      fallbackData: [],
      enableProgressiveDegradation: true,
      staleTime: 5 * 60 * 1000, // 5 minutes
      onError: (error) => {
        console.error('Wallet loading failed:', error)
        
        // Custom handling for critical errors
        if (error.severity === 'CRITICAL') {
          // Could trigger additional error reporting
          console.error('Critical wallet error - may need immediate attention')
        }
      }
    }
  )

  // Enhanced create wallet mutation
  const createWalletMutation = useEnhancedMutation(
    (wallet: CreateWalletRequest) => enhancedWalletService.createWallet(wallet),
    {
      showErrorToast: true,
      showSuccessToast: true,
      successMessage: '錢包建立成功！',
      retryServerErrors: true,
      maxServerErrorRetries: 2,
      onSuccess: () => {
        queryClient.invalidateQueries(['wallets', DEMO_USER_ID])
        setShowCreateModal(false)
        resetForm()
      },
      onError: (error, variables) => {
        console.error('Wallet creation failed:', { error, variables })
        
        // Keep modal open on error for user to retry
        // Don't reset form data so user doesn't lose input
      }
    }
  )

  // Enhanced delete wallet mutation
  const deleteWalletMutation = useEnhancedMutation(
    (walletID: string) => enhancedWalletService.deleteWallet(walletID),
    {
      showErrorToast: true,
      showSuccessToast: true,
      successMessage: '錢包已刪除',
      retryServerErrors: false, // Don't retry deletes automatically
      onSuccess: () => {
        queryClient.invalidateQueries(['wallets', DEMO_USER_ID])
      }
    }
  )

  const resetForm = () => {
    setFormData({ name: '', type: WalletType.CASH, currency: 'TWD', initialBalance: 0 })
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!formData.name.trim()) {
      toast.error('請輸入錢包名稱')
      return
    }
    
    createWalletMutation.mutate({
      name: formData.name.trim(),
      type: formData.type,
      currency: formData.currency,
      user_id: DEMO_USER_ID,
      initialBalance: formData.initialBalance
    })
  }

  // Render different states based on error classification
  if (classifiedError && !isDegraded) {
    return (
      <ErrorBoundary level="page">
        <div className="space-y-6 animate-fade-in">
          {/* Header */}
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div>
              <h1 className="text-3xl font-bold text-gradient-primary">錢包管理</h1>
              <p className="text-neutral-600 mt-1">管理您的錢包和帳戶</p>
            </div>
          </div>
          
          {/* Enhanced error display */}
          <ServerErrorState 
            statusCode={classifiedError.statusCode}
            customMessage={getUserErrorMessage()}
            onRetry={canRetry() ? retryQuery : undefined}
            onRefresh={() => window.location.reload()}
            showTechnicalDetails={process.env.NODE_ENV === 'development'}
          />
        </div>
      </ErrorBoundary>
    )
  }

  return (
    <ErrorBoundary level="page">
      <div className="space-y-6 animate-fade-in">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold text-gradient-primary">錢包管理</h1>
            <p className="text-neutral-600 mt-1">管理您的錢包和帳戶</p>
            
            {/* Degraded mode indicator */}
            {isDegraded && (
              <div className="mt-2">
                <InlineError 
                  error={classifiedError!}
                  onRetry={canRetry() ? retryQuery : undefined}
                  className="text-sm"
                />
              </div>
            )}
          </div>
          
          <Button 
            variant="primary"
            onClick={() => setShowCreateModal(true)}
            disabled={createWalletMutation.isLoading}
          >
            <Plus className="w-4 h-4" />
            新增錢包
          </Button>
        </div>

        {/* Wallets Content with Error Boundaries */}
        <SectionErrorBoundary>
          <LoadingWithError
            isLoading={isLoading}
            error={classifiedError && !isDegraded ? classifiedError : undefined}
            onRetry={canRetry() ? retryQuery : undefined}
          >
            {wallets.length > 0 ? (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {wallets.map((wallet) => (
                  <SectionErrorBoundary key={wallet.id}>
                    <Card glass hover className="card-hover">
                      <CardContent>
                        <div className="flex items-start justify-between mb-4">
                          <div className="flex items-center gap-3">
                            <div className="p-3 bg-primary-100 rounded-xl text-primary-600">
                              {getWalletIcon(wallet.type)}
                            </div>
                            <div>
                              <h3 className="font-semibold text-neutral-800">{wallet.name}</h3>
                              <p className="text-sm text-neutral-500">
                                {getWalletTypeDisplayName(wallet.type)}
                              </p>
                            </div>
                          </div>
                          
                          <div className="flex gap-1">
                            <Button 
                              variant="ghost" 
                              size="sm"
                              onClick={() => {
                                // Edit functionality would go here
                                console.log('Edit wallet:', wallet.id)
                              }}
                            >
                              <Edit className="w-4 h-4" />
                            </Button>
                            <Button 
                              variant="ghost" 
                              size="sm"
                              disabled={deleteWalletMutation.isLoading}
                              onClick={() => {
                                if (window.confirm(`確定要刪除錢包「${wallet.name}」嗎？`)) {
                                  deleteWalletMutation.mutate(wallet.id)
                                }
                              }}
                            >
                              <Trash2 className="w-4 h-4 text-accent-600" />
                            </Button>
                          </div>
                        </div>
                        
                        <div className="text-center py-4 border-t border-primary-100">
                          <p className="text-sm text-neutral-600 mb-1">餘額</p>
                          <p className="text-2xl font-bold text-gradient-primary">
                            {formatMoney(wallet.balance)}
                          </p>
                        </div>
                      </CardContent>
                    </Card>
                  </SectionErrorBoundary>
                ))}
              </div>
            ) : (
              <Card glass className="text-center py-12">
                <CardContent>
                  <WalletIcon className="w-16 h-16 text-neutral-300 mx-auto mb-4" />
                  <h3 className="text-lg font-medium text-neutral-700 mb-2">尚無錢包</h3>
                  <p className="text-neutral-500 mb-6">建立您的第一個錢包開始記帳</p>
                  
                  <Button 
                    variant="primary"
                    onClick={() => setShowCreateModal(true)}
                    disabled={createWalletMutation.isLoading}
                  >
                    <Plus className="w-4 h-4" />
                    建立錢包
                  </Button>
                </CardContent>
              </Card>
            )}
          </LoadingWithError>
        </SectionErrorBoundary>

        {/* Create Modal with Enhanced Error Handling */}
        <Modal
          isOpen={showCreateModal}
          onClose={() => {
            setShowCreateModal(false)
            resetForm()
          }}
          title={
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-lg bg-primary-100 text-primary-600">
                <Plus className="w-5 h-5" />
              </div>
              <h2 className="text-lg font-semibold">新增錢包</h2>
            </div>
          }
        >
          <form onSubmit={handleSubmit} className="space-y-4">
            {/* Show mutation error inline */}
            {createWalletMutation.classifiedError && (
              <InlineError 
                error={createWalletMutation.classifiedError}
                onRetry={createWalletMutation.canRetry() ? () => handleSubmit({preventDefault: () => {}} as any) : undefined}
              />
            )}
            
            <Input
              label="錢包名稱"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              placeholder="輸入錢包名稱"
              required
            />
            
            <Select
              label="錢包類型"
              value={formData.type}
              onChange={(e) => setFormData({ ...formData, type: e.target.value as WalletType })}
              options={[
                { value: 'CASH', label: '現金' },
                { value: 'BANK', label: '銀行帳戶' },
                { value: 'CREDIT', label: '信用卡' },
                { value: 'INVESTMENT', label: '投資帳戶' }
              ]}
            />

            <Select
              label="貨幣"
              value={formData.currency}
              onChange={(e) => setFormData({ ...formData, currency: e.target.value })}
              options={[
                { value: 'TWD', label: '台幣 (TWD)' },
                { value: 'USD', label: '美元 (USD)' },
                { value: 'JPY', label: '日圓 (JPY)' },
                { value: 'EUR', label: '歐元 (EUR)' },
                { value: 'CNY', label: '人民幣 (CNY)' }
              ]}
            />

            <Input
              label="初始餘額"
              type="number"
              value={formData.initialBalance}
              onChange={(e) => setFormData({ ...formData, initialBalance: Number(e.target.value) })}
              placeholder="0"
              step="0.01"
              min="0"
              required
            />
            
            <div className="flex gap-2 pt-4">
              <Button
                type="button"
                variant="secondary"
                onClick={() => {
                  setShowCreateModal(false)
                  resetForm()
                }}
                className="flex-1"
                disabled={createWalletMutation.isLoading}
              >
                取消
              </Button>
              <Button
                type="submit"
                variant="primary"
                loading={createWalletMutation.isLoading}
                className="flex-1"
              >
                建立
              </Button>
            </div>
          </form>
        </Modal>
      </div>
    </ErrorBoundary>
  )
}

const getWalletIcon = (type: WalletType) => {
  switch (type) {
    case WalletType.CASH:
      return <Banknote className="w-5 h-5" />
    case WalletType.BANK:
      return <WalletIcon className="w-5 h-5" />
    case WalletType.CREDIT:
      return <CreditCard className="w-5 h-5" />
    case WalletType.INVESTMENT:
      return <TrendingUp className="w-5 h-5" />
    default:
      return <WalletIcon className="w-5 h-5" />
  }
}

export default WalletsEnhanced