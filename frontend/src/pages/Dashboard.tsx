import React from 'react'
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
import { dashboardService, walletService } from '@/services'
import { formatMoney } from '@/utils/format'

// Mock user ID for demo purposes
const DEMO_USER_ID = "demo-user-123"

const Dashboard: React.FC = () => {
  // Queries
  const { data: dashboardData, isLoading: dashboardLoading } = useQuery(
    ['dashboard', DEMO_USER_ID],
    () => dashboardService.getDashboardData({ userID: DEMO_USER_ID }),
    { refetchInterval: 30000 }
  )

  const { data: walletsData } = useQuery(
    ['wallets', DEMO_USER_ID],
    () => walletService.getWallets(DEMO_USER_ID)
  )

  // TODO: Implement monthly stats visualization
  // const { data: monthlyStats } = useQuery(
  //   ['monthlyStats', DEMO_USER_ID],
  //   () => dashboardService.getMonthlyStats(DEMO_USER_ID, 6)
  // )

  const dashboard = dashboardData?.data
  const wallets = walletsData?.data || []

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
                {wallets.map((wallet) => (
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
                {wallets.length === 0 && (
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
                <Button variant="primary" className="w-full justify-start">
                  <Plus className="w-4 h-4" />
                  新增支出
                </Button>
                <Button variant="secondary" className="w-full justify-start">
                  <Plus className="w-4 h-4" />
                  新增收入
                </Button>
                <Button variant="outline" className="w-full justify-start">
                  <TrendingUp className="w-4 h-4" />
                  轉帳
                </Button>
                <Button variant="ghost" className="w-full justify-start">
                  <Wallet className="w-4 h-4" />
                  新增錢包
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
          {dashboard?.recentTransactions?.length ? (
            <div className="space-y-3">
              {dashboard.recentTransactions.slice(0, 5).map((transaction: any, index) => (
                <div key={index} className="flex items-center justify-between p-3 hover:bg-primary-50/50 rounded-lg transition-colors">
                  <div className="flex items-center gap-3">
                    <div className={`p-2 rounded-lg ${
                      'amount' in transaction && transaction.amount.amount > 0 
                        ? 'bg-secondary-100 text-secondary-600' 
                        : 'bg-accent-100 text-accent-600'
                    }`}>
                      {'amount' in transaction && transaction.amount.amount > 0 
                        ? <ArrowUpRight className="w-4 h-4" />
                        : <ArrowDownRight className="w-4 h-4" />
                    }
                    </div>
                    <div>
                      <p className="font-medium text-neutral-800">
                        {transaction.description || '交易記錄'}
                      </p>
                      <p className="text-sm text-neutral-500">
                        {new Date(transaction.date || transaction.createdAt).toLocaleDateString()}
                      </p>
                    </div>
                  </div>
                  <div className={`font-semibold ${
                    'amount' in transaction && transaction.amount.amount > 0 
                      ? 'text-secondary-600' 
                      : 'text-accent-600'
                  }`}>
                    {'amount' in transaction 
                      ? formatMoney(transaction.amount)
                      : 'NT$ 0'
                    }
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8">
              <TrendingUp className="w-12 h-12 text-neutral-300 mx-auto mb-4" />
              <p className="text-neutral-500">尚無交易記錄</p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

export default Dashboard