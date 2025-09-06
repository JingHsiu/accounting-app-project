import React, { useState } from 'react'
import { TrendingUp, Filter, Search, Plus, Calendar } from 'lucide-react'
import { Card, CardContent, Button, Input, Select, Modal } from '@/components/ui'
import { useIncomes, useDeleteIncome } from '@/hooks'
import { useQuery } from 'react-query'
import { walletService, categoryService } from '@/services'
import TransactionCard from './TransactionCard'
import IncomeForm from '@/components/forms/IncomeForm'
import type { IncomeRecord, IncomeExpenseFilters } from '@/types'
import { CategoryType } from '@/types'
import { formatMoney } from '@/utils/format'

interface IncomeListProps {
  userID: string
  initialFilters?: IncomeExpenseFilters
}

const IncomeList: React.FC<IncomeListProps> = ({
  userID,
  initialFilters = {}
}) => {
  const [filters, setFilters] = useState<IncomeExpenseFilters>(initialFilters)
  const [showFilters, setShowFilters] = useState(false)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [editingIncome, setEditingIncome] = useState<IncomeRecord | null>(null)
  const [searchTerm, setSearchTerm] = useState('')

  // Fetch data
  const { 
    data: incomes = [], 
    isLoading, 
    error,
    refetch 
  } = useIncomes(filters, {
    onError: (error) => {
      console.error('❌ Failed to load incomes:', error)
    }
  })

  const { data: wallets = [] } = useQuery(
    ['wallets', userID],
    () => walletService.getWallets(userID, 'IncomeList')
  )

  const { data: categoriesResponse } = useQuery(
    'categories',
    () => categoryService.getCategories()
  )

  // Extract categories from API response and filter income categories
  const categories = categoriesResponse?.success ? categoriesResponse.data || [] : []
  const incomeCategories = categories.filter((cat: any) => cat.type === CategoryType.INCOME)

  // Delete mutation
  const deleteIncome = useDeleteIncome({
    onSuccess: () => {
      console.log('✅ Income deleted successfully')
    },
    onError: (error) => {
      console.error('❌ Failed to delete income:', error)
    }
  })

  // Filter incomes by search term
  const filteredIncomes = incomes.filter(income =>
    income.description.toLowerCase().includes(searchTerm.toLowerCase())
  )

  // Calculate total
  const totalAmount = filteredIncomes.reduce((sum, income) => {
    // Convert to base currency (TWD) for calculation - simplified approach
    const amount = income.amount.currency === 'TWD' ? income.amount.amount : income.amount.amount * 30 // Rough conversion
    return sum + amount
  }, 0)

  const handleEditIncome = (income: IncomeRecord) => {
    setEditingIncome(income)
    setShowCreateModal(true)
  }

  const handleDeleteIncome = (income: IncomeRecord) => {
    if (window.confirm('確定要刪除這筆收入記錄嗎？')) {
      deleteIncome.mutate(income.id)
    }
  }

  const handleFormSuccess = () => {
    setShowCreateModal(false)
    setEditingIncome(null)
    refetch()
  }

  const handleFormCancel = () => {
    setShowCreateModal(false)
    setEditingIncome(null)
  }

  const clearFilters = () => {
    setFilters({})
    setSearchTerm('')
  }

  // Create options for selects
  const walletOptions = [
    { value: '', label: '全部錢包' },
    ...wallets.map(wallet => ({
      value: wallet.id,
      label: `${wallet.name} (${wallet.currency})`
    }))
  ]

  const categoryOptions = [
    { value: '', label: '全部分類' },
    ...incomeCategories.map((category: any) => ({
      value: category.id,
      label: category.name
    }))
  ]

  if (error) {
    return (
      <Card glass className="text-center py-12">
        <CardContent>
          <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <span className="text-2xl">❌</span>
          </div>
          <h3 className="text-lg font-medium text-red-700 mb-2">載入收入記錄失敗</h3>
          <p className="text-red-600 mb-4">
            {error instanceof Error ? error.message : '發生未知錯誤'}
          </p>
          <Button variant="secondary" onClick={() => refetch()}>
            重新載入
          </Button>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 className="text-2xl font-bold text-gradient-primary flex items-center gap-2">
            <TrendingUp className="w-6 h-6 text-green-600" />
            收入記錄
          </h2>
          <p className="text-neutral-600 mt-1">
            總計 {filteredIncomes.length} 筆記錄
            {totalAmount > 0 && (
              <span className="ml-2 text-green-600 font-medium">
                • {formatMoney({ amount: totalAmount, currency: 'TWD' })}
              </span>
            )}
          </p>
        </div>
        <Button
          variant="primary"
          onClick={() => setShowCreateModal(true)}
        >
          <Plus className="w-4 h-4" />
          新增收入
        </Button>
      </div>

      {/* Search and Filters */}
      <Card glass>
        <CardContent>
          <div className="flex flex-col gap-4">
            <div className="flex flex-col sm:flex-row gap-3">
              <Input
                placeholder="搜尋收入記錄..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                icon={<Search className="w-4 h-4" />}
                className="flex-1"
              />
              <Button
                variant={showFilters ? "primary" : "outline"}
                onClick={() => setShowFilters(!showFilters)}
              >
                <Filter className="w-4 h-4" />
                篩選
              </Button>
            </div>

            {showFilters && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 pt-4 border-t border-neutral-200">
                <Select
                  label="錢包"
                  value={filters.walletID || ''}
                  onChange={(e) => setFilters({ ...filters, walletID: e.target.value || undefined })}
                  options={walletOptions}
                />
                
                <Select
                  label="分類"
                  value={filters.categoryID || ''}
                  onChange={(e) => setFilters({ ...filters, categoryID: e.target.value || undefined })}
                  options={categoryOptions}
                />
                
                <Input
                  label="開始日期"
                  type="date"
                  value={filters.startDate || ''}
                  onChange={(e) => setFilters({ ...filters, startDate: e.target.value || undefined })}
                  icon={<Calendar className="w-4 h-4" />}
                />
                
                <Input
                  label="結束日期"
                  type="date"
                  value={filters.endDate || ''}
                  onChange={(e) => setFilters({ ...filters, endDate: e.target.value || undefined })}
                  icon={<Calendar className="w-4 h-4" />}
                />
                
                <div className="md:col-span-2 lg:col-span-4 flex gap-2">
                  <Button variant="outline" onClick={clearFilters}>
                    清除篩選
                  </Button>
                </div>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Income List */}
      {isLoading ? (
        <div className="space-y-4">
          {[...Array(5)].map((_, i) => (
            <Card key={i} glass className="animate-pulse">
              <div className="h-32 bg-primary-200/20 rounded" />
            </Card>
          ))}
        </div>
      ) : filteredIncomes.length > 0 ? (
        <div className="space-y-4">
          {filteredIncomes
            .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
            .map((income) => {
              const wallet = wallets.find(w => w.id === income.walletID)
              const category = categories.find((c: any) => c.id === income.categoryID)
              
              return (
                <TransactionCard
                  key={income.id}
                  transaction={income}
                  type="income"
                  wallet={wallet}
                  category={category}
                  onEdit={handleEditIncome}
                  onDelete={handleDeleteIncome}
                  isLoading={deleteIncome.isLoading}
                />
              )
            })}
        </div>
      ) : (
        <Card glass className="text-center py-12">
          <CardContent>
            <TrendingUp className="w-16 h-16 text-neutral-300 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-neutral-700 mb-2">尚無收入記錄</h3>
            <p className="text-neutral-500 mb-6">建立您的第一筆收入記錄</p>
            <Button
              variant="primary"
              onClick={() => setShowCreateModal(true)}
            >
              <Plus className="w-4 h-4" />
              新增收入
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Create/Edit Modal */}
      <Modal
        isOpen={showCreateModal}
        onClose={handleFormCancel}
        title={editingIncome ? '編輯收入記錄' : '新增收入記錄'}
        size="lg"
      >
        <IncomeForm
          userID={userID}
          initialData={editingIncome || undefined}
          onSuccess={handleFormSuccess}
          onCancel={handleFormCancel}
        />
      </Modal>
    </div>
  )
}

export default IncomeList