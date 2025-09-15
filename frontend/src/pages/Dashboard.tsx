import React, { useState } from 'react'
import { useQuery } from 'react-query'
import { 
  TrendingUp, 
  Wallet, 
  DollarSign,
  Plus,
  ArrowUpRight,
  ArrowDownRight
} from 'lucide-react'
import { Card, CardHeader, CardTitle, CardContent, Button } from '@/components/ui'
import { walletService, categoryService } from '@/services'
import { formatMoney } from '@/utils/format'
import { useCreateExpense, useCreateIncome, useExpenses, useIncomes } from '@/hooks'
import WalletDebugPanel from '@/components/WalletDebugPanel'
import { EnhancedTransactionModal } from '@/components/EnhancedTransactionModal'
import { EnhancedTransactionList } from '@/components/EnhancedTransactionItem'
import { CategoryType } from '@/types'

// Mock user ID for demo purposes
const DEMO_USER_ID = "demo-user-123"

// Modal types for better type safety
type ModalType = 'expense' | 'income' | null

const Dashboard: React.FC = () => {
  // Modal state management
  const [activeModal, setActiveModal] = useState<ModalType>(null)

  // Queries - Dashboard API not implemented yet, so disable for now
  // const { data: dashboardData, isLoading: dashboardLoading } = useQuery(
  //   ['dashboard', DEMO_USER_ID],
  //   () => dashboardService.getDashboardData({ userID: DEMO_USER_ID }),
  //   { refetchInterval: 30000 }
  // )
  const dashboardData = null
  const dashboardLoading = false

  const { data: walletsData, isLoading: walletsLoading, error: walletsError } = useQuery(
    ['wallets', DEMO_USER_ID],
    () => walletService.getWallets(DEMO_USER_ID, 'Dashboard'),
    {
      onSuccess: (data) => {
        console.group('✅ [Dashboard] React Query SUCCESS')
        console.log('Dashboard wallets data:', {
          data,
          length: data?.length,
          isArray: Array.isArray(data)
        })
        console.groupEnd()
      },
      onError: (error) => {
        console.group('❌ [Dashboard] React Query ERROR')
        console.error('Dashboard query error:', error)
        console.groupEnd()
      }
    }
  )

  // Fetch recent transactions for dashboard display
  const { data: recentExpenses = [] } = useExpenses(
    { walletID: undefined }, // Get all expenses for recent view
    { enabled: true }
  )

  const { data: recentIncomes = [] } = useIncomes(
    { walletID: undefined }, // Get all incomes for recent view  
    { enabled: true }
  )

  // Fetch categories for both display and form dropdowns
  const { data: expenseCategories = [] } = useQuery(
    ['categories', 'expense'],
    () => categoryService.getCategories(CategoryType.EXPENSE),
    { enabled: true } // Always load categories for transaction display
  )

  const { data: incomeCategories = [] } = useQuery(
    ['categories', 'income'],
    () => categoryService.getCategories(CategoryType.INCOME),
    { enabled: true } // Always load categories for transaction display
  )

  // Mutation hooks for creating transactions
  const createExpenseMutation = useCreateExpense({
    onSuccess: () => {
      console.log('✅ Expense created successfully')
      setActiveModal(null)
    },
    onError: (error) => {
      console.error('❌ Failed to create expense:', error)
    }
  })

  const createIncomeMutation = useCreateIncome({
    onSuccess: () => {
      console.log('✅ Income created successfully')
      setActiveModal(null)
    },
    onError: (error) => {
      console.error('❌ Failed to create income:', error)
    }
  })

  // TODO: Implement monthly stats visualization
  // const { data: monthlyStats } = useQuery(
  //   ['monthlyStats', DEMO_USER_ID],
  //   () => dashboardService.getMonthlyStats(DEMO_USER_ID, 6)
  // )

  const dashboard = dashboardData
  const wallets = Array.isArray(walletsData) ? walletsData : []

  // Modal handling functions
  const handleOpenModal = (type: ModalType) => {
    if (wallets.length === 0) {
      console.warn('No wallets available for transaction creation')
      return
    }
    setActiveModal(type)
  }

  const handleCloseModal = () => {
    setActiveModal(null)
  }

  // Enhanced modal submit handler
  const handleTransactionSubmit = (data: {
    wallet_id: string;
    subcategory_id: string;
    amount: number;
    currency: string;
    description?: string;
    date: string;
  }) => {
    if (activeModal === 'expense') {
      createExpenseMutation.mutate(data)
    } else if (activeModal === 'income') {
      createIncomeMutation.mutate(data)
    }
  }

  // Combine and sort recent transactions for display
  const recentTransactions = [
    ...recentExpenses.slice(0, 3).map(expense => ({
      ...expense,
      type: 'expense' as const,
      displayAmount: expense.amount,
      icon: ArrowDownRight,
      colorClass: 'text-red-600 bg-red-100'
    })),
    ...recentIncomes.slice(0, 3).map(income => ({
      ...income,
      type: 'income' as const,
      displayAmount: income.amount,
      icon: ArrowUpRight,
      colorClass: 'text-green-600 bg-green-100'
    }))
  ].sort((a, b) => new Date(b.date || b.createdAt).getTime() - new Date(a.date || a.createdAt).getTime()).slice(0, 5)
  
  // Enhanced debugging for Dashboard
  console.group('🏠 [Dashboard] Component Render Debug')
  console.log('Dashboard render analysis:', {
    dashboardState: {
      dashboardData,
      dashboardLoading
    },
    walletsState: {
      walletsData,
      walletsDataType: typeof walletsData,
      walletsIsArray: Array.isArray(walletsData),
      walletsLength: walletsData?.length || 0,
      walletsLoading,
      walletsError: walletsError?.toString(),
      processedWallets: wallets,
      processedWalletsLength: wallets.length
    },
    renderDecision: {
      willShowDashboardLoading: dashboardLoading,
      willShowWalletsLoading: walletsLoading && !walletsError,
      willShowWalletsError: !!walletsError,
      willShowWallets: wallets.length > 0,
      willShowEmptyWallets: !walletsLoading && !walletsError && wallets.length === 0
    }
  })
  console.groupEnd()

  if (dashboardLoading) {
    return (
      <div className="space-y-6">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {[...Array(4)].map((_, i) => (
            <Card key={i} glass className="animate-pulse">
              <div className="h-24 bg-primary-200/20 rounded" />
            </Card>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold text-gradient-primary">儀表板</h1>
          <p className="text-neutral-600 mt-1">歡迎回來！這是您的財務概況</p>
        </div>
        <div className="flex gap-2">
          <Button variant="secondary" size="sm">
            <Plus className="w-4 h-4" />
            快速記錄
          </Button>
          <Button variant="primary" size="sm">
            <Plus className="w-4 h-4" />
            新增交易
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card glass hover className="card-hover">
          <CardContent>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-neutral-600">總資產</p>
                <p className="text-2xl font-bold text-gradient-primary">
                  {dashboard ? formatMoney(dashboard.totalBalance) : 'NT$ 0'}
                </p>
              </div>
              <div className="p-3 bg-gradient-primary rounded-xl">
                <DollarSign className="w-6 h-6 text-white" />
              </div>
            </div>
          </CardContent>
        </Card>

        <Card glass hover className="card-hover">
          <CardContent>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-neutral-600">本月收入</p>
                <p className="text-2xl font-bold text-gradient-secondary">
                  {dashboard ? formatMoney(dashboard.monthlyIncome) : 'NT$ 0'}
                </p>
              </div>
              <div className="p-3 bg-gradient-secondary rounded-xl">
                <ArrowUpRight className="w-6 h-6 text-white" />
              </div>
            </div>
          </CardContent>
        </Card>

        <Card glass hover className="card-hover">
          <CardContent>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-neutral-600">本月支出</p>
                <p className="text-2xl font-bold text-gradient-accent">
                  {dashboard ? formatMoney(dashboard.monthlyExpense) : 'NT$ 0'}
                </p>
              </div>
              <div className="p-3 bg-gradient-accent rounded-xl">
                <ArrowDownRight className="w-6 h-6 text-white" />
              </div>
            </div>
          </CardContent>
        </Card>

        <Card glass hover className="card-hover">
          <CardContent>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-neutral-600">錢包數量</p>
                <p className="text-2xl font-bold text-primary-600">
                  {wallets.length}
                </p>
              </div>
              <div className="p-3 bg-primary-500 rounded-xl">
                <Wallet className="w-6 h-6 text-white" />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Content */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Wallets Overview */}
        <div className="lg:col-span-2">
          <Card glass>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Wallet className="w-5 h-5 text-primary-600" />
                錢包概況
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {walletsLoading && (
                  <div className="text-center py-4 text-neutral-500">載入中...</div>
                )}
                {walletsError && (
                  <div className="text-center py-4 text-red-500">載入錢包失敗</div>
                )}
                {!walletsLoading && !walletsError && wallets.map((wallet) => (
                  <div key={wallet.id} className="flex items-center justify-between p-4 bg-white/50 rounded-lg border border-primary-100">
                    <div className="flex items-center gap-3">
                      <div className="p-2 bg-primary-100 rounded-lg">
                        <Wallet className="w-4 h-4 text-primary-600" />
                      </div>
                      <div>
                        <p className="font-medium text-neutral-800">{wallet.name}</p>
                        <p className="text-sm text-neutral-500">{wallet.type}</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="font-semibold text-neutral-800">
                        {formatMoney(wallet.balance)}
                      </p>
                    </div>
                  </div>
                ))}
                {!walletsLoading && !walletsError && wallets.length === 0 && (
                  <div className="text-center py-8">
                    <Wallet className="w-12 h-12 text-neutral-300 mx-auto mb-4" />
                    <p className="text-neutral-500">尚未建立錢包</p>
                    <Button variant="primary" size="sm" className="mt-2">
                      建立第一個錢包
                    </Button>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Quick Actions */}
        <div>
          <Card glass>
            <CardHeader>
              <CardTitle>快速操作</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <Button 
                  variant="primary" 
                  className="w-full justify-start"
                  onClick={() => handleOpenModal('expense')}
                  disabled={wallets.length === 0}
                >
                  <Plus className="w-4 h-4" />
                  新增支出
                </Button>
                <Button 
                  variant="secondary" 
                  className="w-full justify-start"
                  onClick={() => handleOpenModal('income')}
                  disabled={wallets.length === 0}
                >
                  <Plus className="w-4 h-4" />
                  新增收入
                </Button>
                <Button variant="outline" className="w-full justify-start" disabled>
                  <TrendingUp className="w-4 h-4" />
                  轉帳 (即將推出)
                </Button>
                <Button variant="ghost" className="w-full justify-start" disabled>
                  <Wallet className="w-4 h-4" />
                  新增錢包 (即將推出)
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Recent Transactions */}
      <Card glass>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            <span>近期交易</span>
            <Button variant="ghost" size="sm">查看全部</Button>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <EnhancedTransactionList
            transactions={recentTransactions}
            wallets={wallets}
            categories={[...expenseCategories, ...incomeCategories]}
            className="space-y-3"
          />
        </CardContent>
      </Card>
      
      {/* Debug Panel */}
      <WalletDebugPanel userID={DEMO_USER_ID} show={false} />
      
      {/* Enhanced Transaction Creation Modal */}
      <EnhancedTransactionModal
        open={!!activeModal}
        onOpenChange={handleCloseModal}
        onSubmit={handleTransactionSubmit}
        wallets={wallets}
        categories={activeModal === 'expense' ? expenseCategories : incomeCategories}
        type={activeModal || 'expense'}
      />
    </div>
  )
}

export default Dashboard