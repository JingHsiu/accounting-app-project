import React from 'react'
import { TrendingUp, TrendingDown, Loader2 } from 'lucide-react'
import { Card, CardContent, Button } from '@/components/ui'
import { useIncomes, useExpenses, useCreateIncome, useCreateExpense } from '@/hooks'
import { formatMoney, formatDate } from '@/utils/format'

const IncomeExpenseTestPanel: React.FC = () => {
  // Test data fetching
  const { 
    data: incomes = [], 
    isLoading: incomesLoading, 
    error: incomesError,
    refetch: refetchIncomes 
  } = useIncomes({}, {
    onSuccess: (data) => {
      console.log('✅ Incomes loaded successfully:', data)
    },
    onError: (error) => {
      console.error('❌ Failed to load incomes:', error)
    }
  })

  const { 
    data: expenses = [], 
    isLoading: expensesLoading, 
    error: expensesError,
    refetch: refetchExpenses 
  } = useExpenses({}, {
    onSuccess: (data) => {
      console.log('✅ Expenses loaded successfully:', data)
    },
    onError: (error) => {
      console.error('❌ Failed to load expenses:', error)
    }
  })

  // Test mutations
  const createIncome = useCreateIncome({
    onSuccess: (data) => {
      console.log('✅ Test income created:', data)
    }
  })

  const createExpense = useCreateExpense({
    onSuccess: (data) => {
      console.log('✅ Test expense created:', data)
    }
  })

  const handleCreateTestIncome = () => {
    createIncome.mutate({
      wallet_id: 'test-wallet-id',
      subcategory_id: 'test-category-id',
      amount: 1000,
      currency: 'TWD',
      description: '測試收入記錄',
      date: new Date().toISOString()
    })
  }

  const handleCreateTestExpense = () => {
    createExpense.mutate({
      wallet_id: 'test-wallet-id',
      subcategory_id: 'test-category-id',
      amount: 500,
      currency: 'TWD',
      description: '測試支出記錄',
      date: new Date().toISOString()
    })
  }

  return (
    <div className="space-y-6 p-6 bg-neutral-50 border-l-4 border-blue-500">
      <div className="flex items-center gap-3">
        <div className="p-2 bg-blue-100 text-blue-600 rounded-lg">
          <TrendingUp className="w-5 h-5" />
        </div>
        <div>
          <h3 className="text-lg font-semibold text-neutral-800">
            收入支出 API 測試面板
          </h3>
          <p className="text-sm text-neutral-500">
            測試新的 API 整合和 React Query hooks
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Income Section */}
        <Card>
          <CardContent>
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-2">
                <TrendingUp className="w-5 h-5 text-green-600" />
                <h4 className="font-medium text-neutral-800">收入記錄</h4>
              </div>
              <div className="flex gap-2">
                <Button 
                  size="sm" 
                  variant="outline"
                  onClick={() => refetchIncomes()}
                  disabled={incomesLoading}
                >
                  {incomesLoading ? <Loader2 className="w-4 h-4 animate-spin" /> : '刷新'}
                </Button>
                <Button 
                  size="sm" 
                  variant="primary"
                  onClick={handleCreateTestIncome}
                  disabled={createIncome.isLoading}
                >
                  {createIncome.isLoading ? <Loader2 className="w-4 h-4 animate-spin" /> : '建立測試'}
                </Button>
              </div>
            </div>

            {incomesError ? (
              <div className="text-red-600 text-sm bg-red-50 p-3 rounded">
                <p className="font-medium">載入失敗</p>
                <p>{incomesError instanceof Error ? incomesError.message : String(incomesError || '未知錯誤')}</p>
              </div>
            ) : incomesLoading ? (
              <div className="flex items-center gap-2 text-neutral-500">
                <Loader2 className="w-4 h-4 animate-spin" />
                <span>載入中...</span>
              </div>
            ) : incomes.length > 0 ? (
              <div className="space-y-3">
                <p className="text-sm text-green-600 font-medium">
                  ✅ 載入成功 ({incomes.length} 筆記錄)
                </p>
                <div className="space-y-2 max-h-32 overflow-y-auto">
                  {incomes.slice(0, 3).map((income) => (
                    <div key={income.id} className="text-xs bg-green-50 p-2 rounded border">
                      <div className="flex justify-between items-center">
                        <span className="font-medium">{income.description}</span>
                        <span className="text-green-600">{formatMoney(income.amount)}</span>
                      </div>
                      <div className="text-neutral-500 mt-1">
                        {formatDate(income.date)}
                      </div>
                    </div>
                  ))}
                  {incomes.length > 3 && (
                    <p className="text-xs text-neutral-500 text-center">
                      還有 {incomes.length - 3} 筆記錄...
                    </p>
                  )}
                </div>
              </div>
            ) : (
              <p className="text-sm text-neutral-500">尚無收入記錄</p>
            )}

            {createIncome.error && (
              <div className="mt-3 text-red-600 text-sm bg-red-50 p-2 rounded">
                建立失敗: {createIncome.error instanceof Error ? createIncome.error.message : String(createIncome.error || '未知錯誤')}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Expense Section */}
        <Card>
          <CardContent>
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-2">
                <TrendingDown className="w-5 h-5 text-red-600" />
                <h4 className="font-medium text-neutral-800">支出記錄</h4>
              </div>
              <div className="flex gap-2">
                <Button 
                  size="sm" 
                  variant="outline"
                  onClick={() => refetchExpenses()}
                  disabled={expensesLoading}
                >
                  {expensesLoading ? <Loader2 className="w-4 h-4 animate-spin" /> : '刷新'}
                </Button>
                <Button 
                  size="sm" 
                  variant="primary"
                  onClick={handleCreateTestExpense}
                  disabled={createExpense.isLoading}
                >
                  {createExpense.isLoading ? <Loader2 className="w-4 h-4 animate-spin" /> : '建立測試'}
                </Button>
              </div>
            </div>

            {expensesError ? (
              <div className="text-red-600 text-sm bg-red-50 p-3 rounded">
                <p className="font-medium">載入失敗</p>
                <p>{expensesError instanceof Error ? expensesError.message : String(expensesError || '未知錯誤')}</p>
              </div>
            ) : expensesLoading ? (
              <div className="flex items-center gap-2 text-neutral-500">
                <Loader2 className="w-4 h-4 animate-spin" />
                <span>載入中...</span>
              </div>
            ) : expenses.length > 0 ? (
              <div className="space-y-3">
                <p className="text-sm text-red-600 font-medium">
                  ✅ 載入成功 ({expenses.length} 筆記錄)
                </p>
                <div className="space-y-2 max-h-32 overflow-y-auto">
                  {expenses.slice(0, 3).map((expense) => (
                    <div key={expense.id} className="text-xs bg-red-50 p-2 rounded border">
                      <div className="flex justify-between items-center">
                        <span className="font-medium">{expense.description}</span>
                        <span className="text-red-600">{formatMoney(expense.amount)}</span>
                      </div>
                      <div className="text-neutral-500 mt-1">
                        {formatDate(expense.date)}
                      </div>
                    </div>
                  ))}
                  {expenses.length > 3 && (
                    <p className="text-xs text-neutral-500 text-center">
                      還有 {expenses.length - 3} 筆記錄...
                    </p>
                  )}
                </div>
              </div>
            ) : (
              <p className="text-sm text-neutral-500">尚無支出記錄</p>
            )}

            {createExpense.error && (
              <div className="mt-3 text-red-600 text-sm bg-red-50 p-2 rounded">
                建立失敗: {createExpense.error instanceof Error ? createExpense.error.message : String(createExpense.error || '未知錯誤')}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Summary */}
      <Card>
        <CardContent>
          <h4 className="font-medium text-neutral-800 mb-3">API 整合狀態</h4>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span className="text-neutral-600">收入 API:</span>
              <span className={`ml-2 font-medium ${incomesError ? 'text-red-600' : 'text-green-600'}`}>
                {incomesError ? '❌ 錯誤' : '✅ 正常'}
              </span>
            </div>
            <div>
              <span className="text-neutral-600">支出 API:</span>
              <span className={`ml-2 font-medium ${expensesError ? 'text-red-600' : 'text-green-600'}`}>
                {expensesError ? '❌ 錯誤' : '✅ 正常'}
              </span>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

export default IncomeExpenseTestPanel