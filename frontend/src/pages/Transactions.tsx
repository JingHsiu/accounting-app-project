import React, { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from 'react-query'
import { 
  Receipt, 
  Plus, 
  ArrowUpRight, 
  ArrowDownRight
} from 'lucide-react'
import { Card, CardHeader, CardTitle, CardContent, Button, Modal, Input, Select } from '@/components/ui'
import { transactionService, walletService, categoryService } from '@/services'
import { formatMoney, formatDate } from '@/utils/format'
import type { AddExpenseRequest, AddIncomeRequest } from '@/services/transactionService'
import { CategoryType } from '@/types'

const DEMO_USER_ID = "demo-user-123"

const Transactions: React.FC = () => {
  const [showAddModal, setShowAddModal] = useState(false)
  const [transactionType, setTransactionType] = useState<'income' | 'expense'>('expense')
  const [filters, setFilters] = useState({
    walletID: '',
    categoryID: '',
    startDate: '',
    endDate: ''
  })
  const [formData, setFormData] = useState({
    walletID: '',
    categoryID: '',
    amount: '',
    description: '',
    date: new Date().toISOString().split('T')[0]
  })

  const queryClient = useQueryClient()

  // Queries
  const { data: walletsData } = useQuery(
    ['wallets', DEMO_USER_ID],
    () => walletService.getWallets(DEMO_USER_ID)
  )

  const { data: categoriesData } = useQuery(
    ['categories'],
    () => categoryService.getCategories()
  )

  const { data: expensesData, isLoading: expensesLoading } = useQuery(
    ['expenses', filters],
    () => transactionService.getExpenses(filters),
    { enabled: Object.values(filters).some(v => v) || Object.keys(filters).length === 4 }
  )

  const { data: incomesData, isLoading: incomesLoading } = useQuery(
    ['incomes', filters],
    () => transactionService.getIncomes(filters),
    { enabled: Object.values(filters).some(v => v) || Object.keys(filters).length === 4 }
  )

  // Mutations
  const addExpenseMutation = useMutation(
    (expense: AddExpenseRequest) => transactionService.addExpense(expense),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['expenses'])
        queryClient.invalidateQueries(['wallets'])
        setShowAddModal(false)
        resetForm()
      }
    }
  )

  const addIncomeMutation = useMutation(
    (income: AddIncomeRequest) => transactionService.addIncome(income),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['incomes'])
        queryClient.invalidateQueries(['wallets'])
        setShowAddModal(false)
        resetForm()
      }
    }
  )

  const wallets = walletsData?.data || []
  const categories = categoriesData?.data || []
  const expenses = expensesData?.data || []
  const incomes = incomesData?.data || []

  // Combine and sort transactions by date
  const allTransactions = [
    ...expenses.map(e => ({ ...e, type: 'expense' as const })),
    ...incomes.map(i => ({ ...i, type: 'income' as const }))
  ].sort((a, b) => new Date(b.date || b.createdAt).getTime() - new Date(a.date || a.createdAt).getTime())

  const resetForm = () => {
    setFormData({
      walletID: '',
      categoryID: '',
      amount: '',
      description: '',
      date: new Date().toISOString().split('T')[0]
    })
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    const transactionData = {
      ...formData,
      amount: {
        amount: parseFloat(formData.amount),
        currency: 'TWD'
      }
    }

    if (transactionType === 'expense') {
      addExpenseMutation.mutate(transactionData)
    } else {
      addIncomeMutation.mutate(transactionData)
    }
  }

  const walletOptions = [
    { value: '', label: '所有錢包' },
    ...wallets.map(w => ({ value: w.id, label: w.name }))
  ]

  const categoryOptions = [
    { value: '', label: '所有類別' },
    ...categories.map(c => ({ value: c.id, label: c.name }))
  ]

  const getTransactionIcon = (type: 'income' | 'expense') => {
    return type === 'income' 
      ? <ArrowUpRight className="w-4 h-4 text-secondary-600" />
      : <ArrowDownRight className="w-4 h-4 text-accent-600" />
  }

  const isLoading = expensesLoading || incomesLoading

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold text-gradient-primary">交易記錄</h1>
          <p className="text-neutral-600 mt-1">查看和管理您的收支記錄</p>
        </div>
        <div className="flex gap-2">
          <Button 
            variant="secondary"
            onClick={() => {
              setTransactionType('income')
              setShowAddModal(true)
            }}
          >
            <Plus className="w-4 h-4" />
            新增收入
          </Button>
          <Button 
            variant="primary"
            onClick={() => {
              setTransactionType('expense')
              setShowAddModal(true)
            }}
          >
            <Plus className="w-4 h-4" />
            新增支出
          </Button>
        </div>
      </div>

      {/* Filters */}
      <Card glass>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <Select
              label="錢包"
              value={filters.walletID}
              onChange={(e) => setFilters({ ...filters, walletID: e.target.value })}
              options={walletOptions}
            />
            <Select
              label="類別"
              value={filters.categoryID}
              onChange={(e) => setFilters({ ...filters, categoryID: e.target.value })}
              options={categoryOptions}
            />
            <Input
              label="開始日期"
              type="date"
              value={filters.startDate}
              onChange={(e) => setFilters({ ...filters, startDate: e.target.value })}
            />
            <Input
              label="結束日期"
              type="date"
              value={filters.endDate}
              onChange={(e) => setFilters({ ...filters, endDate: e.target.value })}
            />
          </div>
        </CardContent>
      </Card>

      {/* Transactions List */}
      <Card glass>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Receipt className="w-5 h-5 text-primary-600" />
            交易列表
          </CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-3">
              {[...Array(5)].map((_, i) => (
                <div key={i} className="animate-pulse flex items-center justify-between p-4 bg-primary-50/50 rounded-lg">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-primary-200 rounded-lg" />
                    <div>
                      <div className="w-24 h-4 bg-primary-200 rounded mb-2" />
                      <div className="w-16 h-3 bg-primary-100 rounded" />
                    </div>
                  </div>
                  <div className="w-20 h-4 bg-primary-200 rounded" />
                </div>
              ))}
            </div>
          ) : allTransactions.length > 0 ? (
            <div className="space-y-3">
              {allTransactions.map((transaction, index) => {
                const wallet = wallets.find(w => w.id === transaction.walletID)
                const category = categories.find(c => c.id === transaction.categoryID)
                
                return (
                  <div key={`${transaction.type}-${transaction.id}-${index}`} 
                       className="flex items-center justify-between p-4 hover:bg-primary-50/50 rounded-lg transition-colors border border-primary-100">
                    <div className="flex items-center gap-4">
                      <div className={`p-2 rounded-lg ${
                        transaction.type === 'income' 
                          ? 'bg-secondary-100' 
                          : 'bg-accent-100'
                      }`}>
                        {getTransactionIcon(transaction.type)}
                      </div>
                      <div>
                        <p className="font-medium text-neutral-800">
                          {transaction.description}
                        </p>
                        <div className="flex items-center gap-2 text-sm text-neutral-500">
                          <span>{wallet?.name}</span>
                          {category && (
                            <>
                              <span>•</span>
                              <span>{category.name}</span>
                            </>
                          )}
                          <span>•</span>
                          <span>{formatDate(transaction.date || transaction.createdAt)}</span>
                        </div>
                      </div>
                    </div>
                    <div className={`font-semibold text-lg ${
                      transaction.type === 'income' 
                        ? 'text-secondary-600' 
                        : 'text-accent-600'
                    }`}>
                      {transaction.type === 'income' ? '+' : '-'}{formatMoney(transaction.amount)}
                    </div>
                  </div>
                )
              })}
            </div>
          ) : (
            <div className="text-center py-12">
              <Receipt className="w-16 h-16 text-neutral-300 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-neutral-700 mb-2">尚無交易記錄</h3>
              <p className="text-neutral-500 mb-6">開始記錄您的第一筆交易</p>
              <Button 
                variant="primary"
                onClick={() => {
                  setTransactionType('expense')
                  setShowAddModal(true)
                }}
              >
                <Plus className="w-4 h-4" />
                新增交易
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Add Transaction Modal */}
      <Modal
        isOpen={showAddModal}
        onClose={() => {
          setShowAddModal(false)
          resetForm()
        }}
        title={transactionType === 'income' ? '新增收入' : '新增支出'}
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <Select
            label="錢包"
            value={formData.walletID}
            onChange={(e) => setFormData({ ...formData, walletID: e.target.value })}
            options={[{ value: '', label: '選擇錢包' }, ...wallets.map(w => ({ value: w.id, label: w.name }))]}
            required
          />
          
          <Select
            label="類別"
            value={formData.categoryID}
            onChange={(e) => setFormData({ ...formData, categoryID: e.target.value })}
            options={[
              { value: '', label: '選擇類別' }, 
              ...categories
                .filter(c => c.type === (transactionType === 'income' ? CategoryType.INCOME : CategoryType.EXPENSE))
                .map(c => ({ value: c.id, label: c.name }))
            ]}
            required
          />
          
          <Input
            label="金額"
            type="number"
            step="0.01"
            value={formData.amount}
            onChange={(e) => setFormData({ ...formData, amount: e.target.value })}
            placeholder="輸入金額"
            required
          />
          
          <Input
            label="描述"
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            placeholder="輸入描述"
            required
          />
          
          <Input
            label="日期"
            type="date"
            value={formData.date}
            onChange={(e) => setFormData({ ...formData, date: e.target.value })}
            required
          />
          
          <div className="flex gap-2 pt-4">
            <Button
              type="button"
              variant="secondary"
              onClick={() => {
                setShowAddModal(false)
                resetForm()
              }}
              className="flex-1"
            >
              取消
            </Button>
            <Button
              type="submit"
              variant="primary"
              loading={addExpenseMutation.isLoading || addIncomeMutation.isLoading}
              className="flex-1"
            >
              新增
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}

export default Transactions