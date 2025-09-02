import React, { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from 'react-query'
import { 
  Wallet as WalletIcon, 
  Plus, 
  Edit, 
  Trash2, 
  CreditCard,
  Banknote,
  TrendingUp
} from 'lucide-react'
import { Card, CardContent, Button, Modal, Input, Select } from '@/components/ui'
import { walletService } from '@/services'
import { formatMoney, getWalletTypeDisplayName } from '@/utils/format'
import type { CreateWalletRequest } from '@/services/walletService'
import { WalletType } from '@/types'
import WalletDebugPanel from '@/components/WalletDebugPanel'

const DEMO_USER_ID = "demo-user-123"

const Wallets: React.FC = () => {
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [editingWallet, setEditingWallet] = useState<string | null>(null)
  const [formData, setFormData] = useState({
    name: '',
    type: WalletType.CASH,
    currency: 'TWD',
    initialBalance: 0
  })

  const queryClient = useQueryClient()

  // Queries with comprehensive debugging
  const { data: walletsData, isLoading, error, dataUpdatedAt, isStale, isFetching } = useQuery(
    ['wallets', DEMO_USER_ID],
    () => walletService.getWallets(DEMO_USER_ID, 'WalletsPage'),
    {
      onSuccess: (data) => {
        console.group('✅ [WalletsPage] React Query SUCCESS')
        console.log('Success data received:', {
          data,
          dataType: typeof data,
          isArray: Array.isArray(data),
          length: data?.length,
          firstItem: data?.[0]
        })
        console.groupEnd()
      },
      onError: (error) => {
        console.group('❌ [WalletsPage] React Query ERROR')
        console.error('Query error:', {
          error,
          errorType: typeof error,
          errorMessage: error instanceof Error ? error.message : 'Unknown error'
        })
        console.groupEnd()
      },
      onSettled: (data, error) => {
        console.group('🏁 [WalletsPage] React Query SETTLED')
        console.log('Query settled:', {
          hasData: !!data,
          hasError: !!error,
          dataLength: data?.length,
          errorMessage: error instanceof Error ? error.message : error
        })
        console.groupEnd()
      },
      staleTime: 0, // Always refetch
      cacheTime: 0  // Don't cache for debugging
    }
  )

  // Mutations
  const createWalletMutation = useMutation(
    (wallet: CreateWalletRequest) => walletService.createWallet(wallet),
    {
      onSuccess: (data) => {
        console.group('✅ [WalletsPage] Create Wallet SUCCESS')
        console.log('Mutation success data:', data)
        console.log('Invalidating wallet cache and refreshing...')
        queryClient.invalidateQueries(['wallets', DEMO_USER_ID])
        setShowCreateModal(false)
        resetForm()
        console.groupEnd()
      },
      onError: (error) => {
        console.group('❌ [WalletsPage] Create Wallet ERROR')
        console.error('Mutation error:', {
          error,
          errorType: typeof error,
          errorMessage: error instanceof Error ? error.message : 'Unknown error'
        })
        console.groupEnd()
      },
      onMutate: (wallet) => {
        console.group('🚀 [WalletsPage] Create Wallet MUTATE')
        console.log('Starting wallet creation:', wallet)
        console.groupEnd()
      },
      onSettled: (data, error) => {
        console.group('🏁 [WalletsPage] Create Wallet SETTLED')
        console.log('Mutation settled:', { 
          hasData: !!data, 
          hasError: !!error,
          currentWalletsCount: wallets.length 
        })
        console.groupEnd()
      }
    }
  )

  const deleteWalletMutation = useMutation(
    (walletID: string) => walletService.deleteWallet(walletID),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['wallets', DEMO_USER_ID])
      }
    }
  )

  const updateWalletMutation = useMutation(
    ({ walletID, updates }: { walletID: string; updates: { name?: string; type?: WalletType } }) => 
      walletService.updateWallet(walletID, updates),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['wallets', DEMO_USER_ID])
        setShowCreateModal(false)
        resetForm()
      },
      onError: (error) => {
        console.error('錢包更新失敗:', error)
        // 可以在這裡加入錯誤提示
      }
    }
  )

  // wallets data is now directly returned from the service  
  const wallets = walletsData || []
  
  // Enhanced component render debugging
  console.group('🔄 [WalletsPage] Component Render Debug')
  console.log('Render state analysis:', {
    queryState: {
      isLoading,
      isFetching,
      isStale,
      error: error?.toString(),
      dataUpdatedAt: new Date(dataUpdatedAt || 0).toLocaleTimeString(),
      hasError: !!error
    },
    dataAnalysis: {
      walletsData,
      walletsDataType: typeof walletsData,
      walletsDataIsArray: Array.isArray(walletsData),
      walletsDataLength: walletsData?.length,
      processedWallets: wallets,
      processedWalletsLength: wallets.length,
      processedWalletsType: typeof wallets,
      firstWallet: wallets[0]
    },
    renderDecision: {
      willShowLoading: isLoading,
      willShowWallets: wallets.length > 0,
      willShowEmptyState: !isLoading && wallets.length === 0,
      willShowError: !!error
    }
  })
  console.groupEnd()

  const resetForm = () => {
    setFormData({ name: '', type: WalletType.CASH, currency: 'TWD', initialBalance: 0 })
    setEditingWallet(null)
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    console.log('🚀 HandleSubmit called with formData:', formData)
    
    // Validate required fields
    if (!formData.name.trim()) {
      console.error('❌ Wallet name is required')
      return
    }
    
    if (editingWallet) {
      console.log('📝 Updating wallet:', editingWallet)
      updateWalletMutation.mutate({
        walletID: editingWallet,
        updates: {
          name: formData.name.trim(),
          type: formData.type
        }
      })
    } else {
      console.log('➕ Creating new wallet with payload:', {
        name: formData.name.trim(),
        type: formData.type,
        currency: formData.currency,
        user_id: DEMO_USER_ID,
        initialBalance: formData.initialBalance
      })
      
      createWalletMutation.mutate({
        name: formData.name.trim(),
        type: formData.type,
        currency: formData.currency,
        user_id: DEMO_USER_ID,
        initialBalance: formData.initialBalance
      })
    }
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

  const walletTypeOptions = [
    { value: 'CASH', label: '現金' },
    { value: 'BANK', label: '銀行帳戶' },
    { value: 'CREDIT', label: '信用卡' },
    { value: 'INVESTMENT', label: '投資帳戶' }
  ]

  const currencyOptions = [
    { value: 'TWD', label: '台幣 (TWD)' },
    { value: 'USD', label: '美元 (USD)' },
    { value: 'JPY', label: '日圓 (JPY)' },
    { value: 'EUR', label: '歐元 (EUR)' },
    { value: 'CNY', label: '人民幣 (CNY)' }
  ]

  if (isLoading || isFetching) {
    return (
      <div className="space-y-6 animate-fade-in">
        {/* Header with loading state */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold text-gradient-primary">錢包管理</h1>
            <p className="text-neutral-600 mt-1 flex items-center gap-2">
              <div className="w-4 h-4 border-2 border-primary-300 border-t-primary-600 rounded-full animate-spin" />
              正在載入錢包資料...
            </p>
          </div>
        </div>
        
        {/* Loading skeleton */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {[...Array(6)].map((_, i) => (
            <Card key={i} glass className="animate-pulse">
              <div className="h-32 bg-primary-200/20 rounded" />
            </Card>
          ))}
        </div>
      </div>
    )
  }

  // Show error state
  if (error) {
    return (
      <div className="space-y-6 animate-fade-in">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold text-gradient-primary">錢包管理</h1>
            <p className="text-neutral-600 mt-1">管理您的錢包和帳戶</p>
          </div>
        </div>
        
        <Card glass className="text-center py-12">
          <CardContent>
            <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <span className="text-2xl">❌</span>
            </div>
            <h3 className="text-lg font-medium text-red-700 mb-2">載入錢包失敗</h3>
            <p className="text-red-600 mb-4">
              {error instanceof Error ? error.message : '發生未知錯誤'}
            </p>
            <Button 
              variant="secondary"
              onClick={() => {
                console.log('🔄 Retry button clicked')
                // Force refetch
                queryClient.invalidateQueries(['wallets', DEMO_USER_ID])
              }}
            >
              重新載入
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold text-gradient-primary">錢包管理</h1>
          <p className="text-neutral-600 mt-1">管理您的錢包和帳戶</p>
        </div>
        <Button 
          variant="primary"
          onClick={() => {
            console.log('🎯 Create wallet button clicked (header)')
            setShowCreateModal(true)
          }}
        >
          <Plus className="w-4 h-4" />
          新增錢包
        </Button>
      </div>

      {/* Wallets Grid */}
      {wallets.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {wallets.map((wallet) => (
            <Card key={wallet.id} glass hover className="card-hover">
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
                      disabled={updateWalletMutation.isLoading && editingWallet === wallet.id}
                      onClick={() => {
                        setEditingWallet(wallet.id)
                        setFormData({ 
                          name: wallet.name, 
                          type: wallet.type, 
                          currency: wallet.currency || 'TWD',
                          initialBalance: wallet.balance.amount || 0
                        })
                        setShowCreateModal(true)
                      }}
                    >
                      <Edit className="w-4 h-4" />
                    </Button>
                    <Button 
                      variant="ghost" 
                      size="sm"
                      onClick={() => deleteWalletMutation.mutate(wallet.id)}
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
          ))}
        </div>
      ) : (
        <Card glass className="text-center py-12">
          <CardContent>
            <WalletIcon className="w-16 h-16 text-neutral-300 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-neutral-700 mb-2">尚無錢包</h3>
            <p className="text-neutral-500 mb-6">建立您的第一個錢包開始記帳</p>
            
            {/* Debug info for empty state */}
            <div className="mb-4 p-3 bg-yellow-50 border border-yellow-200 rounded text-sm text-left">
              <h4 className="font-semibold text-yellow-800 mb-2">🔍 調試資訊 (空狀態)</h4>
              <div className="space-y-1 text-yellow-700">
                <div>walletsData 類型: {typeof walletsData}</div>
                <div>walletsData 是否為陣列: {String(Array.isArray(walletsData))}</div>
                <div>walletsData 長度: {walletsData?.length || 'undefined'}</div>
                <div>處理後 wallets 長度: {wallets.length}</div>
                <div>是否載入中: {String(isLoading)}</div>
                <div>是否有錯誤: {String(!!error)}</div>
                <div>walletsData 內容: {JSON.stringify(walletsData, null, 2)}</div>
              </div>
            </div>
            
            <Button 
              variant="primary"
              onClick={() => {
                console.log('🎯 Create wallet button clicked (empty state)')
                console.log('Current debug state:', { walletsData, wallets, isLoading, error })
                setShowCreateModal(true)
              }}
            >
              <Plus className="w-4 h-4" />
              建立錢包
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Create/Edit Modal */}
      <Modal
        isOpen={showCreateModal}
        onClose={() => {
          setShowCreateModal(false)
          resetForm()
        }}
        title={
          <div className="flex items-center gap-3">
            <div className={`p-2 rounded-lg ${editingWallet ? 'bg-amber-100 text-amber-600' : 'bg-primary-100 text-primary-600'}`}>
              {editingWallet ? <Edit className="w-5 h-5" /> : <Plus className="w-5 h-5" />}
            </div>
            <div>
              <h2 className="text-lg font-semibold">
                {editingWallet ? '編輯錢包' : '新增錢包'}
              </h2>
              {editingWallet && (
                <p className="text-sm text-neutral-500">
                  修改錢包基本資訊
                </p>
              )}
            </div>
          </div>
        }
      >
        <form onSubmit={handleSubmit} className="space-y-4">
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
            options={walletTypeOptions}
          />

          {!editingWallet && (
            <>
              <Select
                label="貨幣"
                value={formData.currency}
                onChange={(e) => setFormData({ ...formData, currency: e.target.value })}
                options={currencyOptions}
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
            </>
          )}
          
          {editingWallet && (
            <div className="bg-neutral-50 border border-neutral-200 rounded-lg p-4 space-y-2">
              <h4 className="text-sm font-medium text-neutral-700">不可修改的資訊</h4>
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <span className="text-neutral-500">貨幣：</span>
                  <span className="font-medium">{formData.currency}</span>
                </div>
                <div>
                  <span className="text-neutral-500">當前餘額：</span>
                  <span className="font-medium">{
                    walletsData?.find(w => w.id === editingWallet)?.balance 
                      ? formatMoney(walletsData.find(w => w.id === editingWallet)!.balance)
                      : '載入中...'
                  }</span>
                </div>
              </div>
              <p className="text-xs text-neutral-500">
                貨幣和餘額無法直接修改，餘額會隨著交易記錄變動
              </p>
            </div>
          )}
          
          <div className="flex gap-2 pt-4">
            <Button
              type="button"
              variant="secondary"
              onClick={() => {
                setShowCreateModal(false)
                resetForm()
              }}
              className="flex-1"
            >
              取消
            </Button>
            <Button
              type="submit"
              variant="primary"
              loading={editingWallet ? updateWalletMutation.isLoading : createWalletMutation.isLoading}
              className="flex-1"
            >
              {editingWallet ? '更新' : '建立'}
            </Button>
          </div>
        </form>
      </Modal>
      
      {/* Debug Panel */}
      <WalletDebugPanel userID={DEMO_USER_ID} show={false} />
    </div>
  )
}

export default Wallets