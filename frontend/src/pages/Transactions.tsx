import React, { useState } from 'react'
import { 
  Receipt, 
  TrendingUp,
  TrendingDown,
  Plus
} from 'lucide-react'
import { Card, CardHeader, CardTitle, CardContent, Button } from '@/components/ui'
import { IncomeList, ExpenseList } from '@/components/lists'
import { IncomeForm, ExpenseForm } from '@/components/forms'
import { EnhancedTransactionList } from '@/components/EnhancedTransactionItem'
import { useIncomes, useExpenses } from '@/hooks'
import { useQuery } from 'react-query'
import { walletService, categoryService } from '@/services'
import type { IncomeExpenseFilters } from '@/types'
import { CategoryType } from '@/types'

const DEMO_USER_ID = "demo-user-123"

const Transactions: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'combined' | 'income' | 'expense'>('combined')
  const [globalFilters] = useState<IncomeExpenseFilters>({})
  const [showIncomeForm, setShowIncomeForm] = useState(false)
  const [showExpenseForm, setShowExpenseForm] = useState(false)

  // Fetch data for combined view with global filters
  const { 
    data: incomes = [], 
    isLoading: incomesLoading 
  } = useIncomes(globalFilters, {
    enabled: activeTab === 'combined' || activeTab === 'income'
  })

  const { 
    data: expenses = [], 
    isLoading: expensesLoading 
  } = useExpenses(globalFilters, {
    enabled: activeTab === 'combined' || activeTab === 'expense'
  })

  // Fetch wallets and categories for enhanced display
  const { data: walletsData } = useQuery(
    ['wallets', DEMO_USER_ID],
    () => walletService.getWallets(DEMO_USER_ID, 'Transactions'),
    {
      onError: (error) => {
        console.error('❌ Failed to load wallets:', error)
      }
    }
  )

  const { data: expenseCategories = [] } = useQuery(
    ['categories', 'expense'],
    () => categoryService.getCategories(CategoryType.EXPENSE)
  )

  const { data: incomeCategories = [] } = useQuery(
    ['categories', 'income'],
    () => categoryService.getCategories(CategoryType.INCOME)
  )

  const wallets = Array.isArray(walletsData) ? walletsData : []
  const categories = [...expenseCategories, ...incomeCategories]

  const isLoading = incomesLoading || expensesLoading

  // For combined view, merge and sort by date
  const allTransactions = [
    ...incomes.map(income => ({ ...income, type: 'income' as const })),
    ...expenses.map(expense => ({ ...expense, type: 'expense' as const }))
  ].sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold text-gradient-primary">交易記錄</h1>
          <p className="text-neutral-600 mt-1">查看和管理您的收支記錄</p>
        </div>
      </div>

      {/* Tab Navigation */}
      <Card glass>
        <CardContent>
          <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
            <div className="flex gap-1 bg-neutral-100 p-1 rounded-lg">
              <Button
                variant={activeTab === 'combined' ? 'primary' : 'ghost'}
                size="sm"
                onClick={() => setActiveTab('combined')}
                className="flex items-center gap-2"
              >
                <Receipt className="w-4 h-4" />
                全部交易
              </Button>
              <Button
                variant={activeTab === 'income' ? 'primary' : 'ghost'}
                size="sm"
                onClick={() => setActiveTab('income')}
                className="flex items-center gap-2"
              >
                <TrendingUp className="w-4 h-4" />
                收入
              </Button>
              <Button
                variant={activeTab === 'expense' ? 'primary' : 'ghost'}
                size="sm"
                onClick={() => setActiveTab('expense')}
                className="flex items-center gap-2"
              >
                <TrendingDown className="w-4 h-4" />
                支出
              </Button>
            </div>
            
            {/* Action Buttons */}
            <div className="flex gap-2">
              {(activeTab === 'income' || activeTab === 'combined') && (
                <Button
                  variant="primary"
                  size="sm"
                  onClick={() => setShowIncomeForm(true)}
                  className="flex items-center gap-2"
                >
                  <Plus className="w-4 h-4" />
                  新增收入
                </Button>
              )}
              {(activeTab === 'expense' || activeTab === 'combined') && (
                <Button
                  variant="primary"
                  size="sm"
                  onClick={() => setShowExpenseForm(true)}
                  className="flex items-center gap-2"
                >
                  <Plus className="w-4 h-4" />
                  新增支出
                </Button>
              )}
            </div>
            
            {/* Summary for combined view */}
            {activeTab === 'combined' && (
              <div className="flex gap-4 text-sm">
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                  <span className="text-neutral-600">
                    收入: <span className="font-medium text-green-600">{incomes.length} 筆</span>
                  </span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 bg-red-500 rounded-full"></div>
                  <span className="text-neutral-600">
                    支出: <span className="font-medium text-red-600">{expenses.length} 筆</span>
                  </span>
                </div>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Content based on active tab */}
      {activeTab === 'combined' && (
        <Card glass>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Receipt className="w-5 h-5 text-primary-600" />
              所有交易記錄
            </CardTitle>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="space-y-4">
                {[...Array(5)].map((_, i) => (
                  <Card key={i} glass className="animate-pulse">
                    <div className="h-24 bg-primary-200/20 rounded" />
                  </Card>
                ))}
              </div>
            ) : allTransactions.length > 0 ? (
              <EnhancedTransactionList
                transactions={allTransactions}
                wallets={wallets}
                categories={categories}
                className="space-y-4"
              />
            ) : (
              <div className="text-center py-12">
                <Receipt className="w-16 h-16 text-neutral-300 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-neutral-700 mb-2">尚無交易記錄</h3>
                <p className="text-neutral-500 mb-6">開始記錄您的第一筆交易</p>
              </div>
            )}
          </CardContent>
        </Card>
      )}
      
      {activeTab === 'income' && (
        <IncomeList 
          userID={DEMO_USER_ID} 
          initialFilters={globalFilters}
        />
      )}
      
      {activeTab === 'expense' && (
        <ExpenseList 
          userID={DEMO_USER_ID} 
          initialFilters={globalFilters}
        />
      )}
      
      {/* Form Modals */}
      {showIncomeForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg max-w-md w-full max-h-[90vh] overflow-y-auto">
            <div className="p-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-xl font-semibold text-neutral-800">新增收入</h2>
                <Button
                  variant="primary"
                  size="sm"
                  onClick={() => setShowIncomeForm(false)}
                  className="text-neutral-500 hover:text-neutral-700"
                >
                  ✕
                </Button>
              </div>
              <IncomeForm
                userID={DEMO_USER_ID}
                onSuccess={() => {
                  setShowIncomeForm(false)
                }}
                onCancel={() => setShowIncomeForm(false)}
              />
            </div>
          </div>
        </div>
      )}

      {showExpenseForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg max-w-md w-full max-h-[90vh] overflow-y-auto">
            <div className="p-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-xl font-semibold text-neutral-800">新增支出</h2>
                <Button
                  variant="primary"
                  size="sm"
                  onClick={() => setShowExpenseForm(false)}
                  className="text-neutral-500 hover:text-neutral-700"
                >
                  ✕
                </Button>
              </div>
              <ExpenseForm
                userID={DEMO_USER_ID}
                onSuccess={() => {
                  setShowExpenseForm(false)
                }}
                onCancel={() => setShowExpenseForm(false)}
              />
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default Transactions