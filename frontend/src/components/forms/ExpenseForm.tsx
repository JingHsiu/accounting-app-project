import React, { useState, useEffect } from 'react'
import { TrendingDown, DollarSign, Calendar, FileText } from 'lucide-react'
import { Button, Input, Select } from '@/components/ui'
import { useQuery } from 'react-query'
import { walletService, categoryService } from '@/services'
import { useCreateExpense, useUpdateExpense } from '@/hooks'
import type { ExpenseRecord, CreateExpenseRequest, UpdateExpenseRequest } from '@/types'
import { CategoryType } from '@/types'
import { formatDate } from '@/utils/format'

interface ExpenseFormProps {
  userID: string
  initialData?: ExpenseRecord
  onSuccess?: (expense: ExpenseRecord | {id: string}) => void
  onCancel?: () => void
}

const ExpenseForm: React.FC<ExpenseFormProps> = ({
  userID,
  initialData,
  onSuccess,
  onCancel
}) => {
  const isEditing = !!initialData
  
  const [formData, setFormData] = useState({
    walletID: initialData?.walletID || '',
    categoryID: initialData?.categoryID || '',
    amount: initialData?.amount?.amount || 0,
    currency: initialData?.amount?.currency || 'TWD',
    description: initialData?.description || '',
    date: initialData?.date ? new Date(initialData.date).toISOString().split('T')[0] : new Date().toISOString().split('T')[0]
  })
  
  const [errors, setErrors] = useState<Record<string, string>>({})

  // Fetch wallets and categories
  const { data: wallets = [], isLoading: walletsLoading } = useQuery(
    ['wallets', userID],
    () => walletService.getWallets(userID, 'ExpenseForm')
  )

  const { data: categoriesResponse, isLoading: categoriesLoading } = useQuery(
    'categories',
    () => categoryService.getCategories()
  )

  // Extract categories from API response and filter expense categories
  const categories = categoriesResponse?.success ? categoriesResponse.data || [] : []
  const expenseCategories = categories.filter((cat: any) => cat.type === CategoryType.EXPENSE)

  // Mutations
  const createExpense = useCreateExpense({
    onSuccess: (data) => {
      console.log('✅ Expense created successfully:', data)
      onSuccess?.(data)
    },
    onError: (error) => {
      console.error('❌ Failed to create expense:', error)
      setErrors({ submit: error.message })
    }
  })

  const updateExpense = useUpdateExpense({
    onSuccess: (data) => {
      console.log('✅ Expense updated successfully:', data)
      onSuccess?.(data)
    },
    onError: (error) => {
      console.error('❌ Failed to update expense:', error)
      setErrors({ submit: error.message })
    }
  })

  // Update currency when wallet changes
  useEffect(() => {
    if (formData.walletID && !isEditing) {
      const selectedWallet = wallets.find(w => w.id === formData.walletID)
      if (selectedWallet && selectedWallet.currency !== formData.currency) {
        setFormData(prev => ({ ...prev, currency: selectedWallet.currency }))
      }
    }
  }, [formData.walletID, wallets, isEditing])

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {}

    if (!formData.walletID.trim()) {
      newErrors.walletID = '請選擇錢包'
    }
    if (!formData.categoryID.trim()) {
      newErrors.categoryID = '請選擇分類'
    }
    if (formData.amount <= 0) {
      newErrors.amount = '金額必須大於 0'
    }
    if (!formData.description.trim()) {
      newErrors.description = '請輸入描述'
    }
    if (!formData.date) {
      newErrors.date = '請選擇日期'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setErrors({})
    
    if (!validateForm()) {
      return
    }

    if (isEditing && initialData) {
      const updates: UpdateExpenseRequest = {
        walletID: formData.walletID,
        categoryID: formData.categoryID,
        amount: formData.amount,
        currency: formData.currency,
        description: formData.description.trim(),
        date: new Date(formData.date).toISOString()
      }
      
      updateExpense.mutate({ expenseID: initialData.id, updates })
    } else {
      const expenseData: CreateExpenseRequest = {
        walletID: formData.walletID,
        categoryID: formData.categoryID,
        amount: formData.amount,
        currency: formData.currency,
        description: formData.description.trim(),
        date: new Date(formData.date).toISOString()
      }
      
      createExpense.mutate(expenseData)
    }
  }

  const isLoading = createExpense.isLoading || updateExpense.isLoading || walletsLoading || categoriesLoading

  const walletOptions = wallets.map(wallet => ({
    value: wallet.id,
    label: `${wallet.name} (${wallet.currency})`
  }))

  const categoryOptions = expenseCategories.map((category: any) => ({
    value: category.id,
    label: category.name
  }))

  const currencyOptions = [
    { value: 'TWD', label: '台幣 (TWD)' },
    { value: 'USD', label: '美元 (USD)' },
    { value: 'JPY', label: '日圓 (JPY)' },
    { value: 'EUR', label: '歐元 (EUR)' },
    { value: 'CNY', label: '人民幣 (CNY)' }
  ]

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3 pb-2 border-b border-primary-100">
        <div className="p-2 bg-red-100 text-red-600 rounded-lg">
          <TrendingDown className="w-5 h-5" />
        </div>
        <div>
          <h3 className="text-lg font-semibold text-neutral-800">
            {isEditing ? '編輯支出記錄' : '新增支出記錄'}
          </h3>
          <p className="text-sm text-neutral-500">
            {isEditing ? '修改支出記錄資訊' : '記錄新的支出項目'}
          </p>
        </div>
      </div>

      {/* Form Fields */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Select
          label="錢包"
          value={formData.walletID}
          onChange={(e) => setFormData(prev => ({ ...prev, walletID: e.target.value }))}
          options={[
            { value: '', label: '選擇錢包...' },
            ...walletOptions
          ]}
          error={errors.walletID}
          disabled={isLoading}
        />

        <Select
          label="分類"
          value={formData.categoryID}
          onChange={(e) => setFormData(prev => ({ ...prev, categoryID: e.target.value }))}
          options={[
            { value: '', label: '選擇分類...' },
            ...categoryOptions
          ]}
          error={errors.categoryID}
          disabled={isLoading || categoriesLoading}
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Input
          label="金額"
          type="number"
          value={formData.amount}
          onChange={(e) => setFormData(prev => ({ ...prev, amount: Number(e.target.value) }))}
          placeholder="0.00"
          step="0.01"
          min="0"
          error={errors.amount}
          disabled={isLoading}
          icon={<DollarSign className="w-4 h-4" />}
        />

        <Select
          label="貨幣"
          value={formData.currency}
          onChange={(e) => setFormData(prev => ({ ...prev, currency: e.target.value }))}
          options={currencyOptions}
          disabled={isEditing || isLoading} // Currency cannot be changed when editing
        />
      </div>

      <Input
        label="描述"
        value={formData.description}
        onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
        placeholder="描述這筆支出..."
        error={errors.description}
        disabled={isLoading}
        icon={<FileText className="w-4 h-4" />}
      />

      <Input
        label="日期"
        type="date"
        value={formData.date}
        onChange={(e) => setFormData(prev => ({ ...prev, date: e.target.value }))}
        error={errors.date}
        disabled={isLoading}
        icon={<Calendar className="w-4 h-4" />}
      />

      {/* Error Message */}
      {errors.submit && (
        <div className="p-3 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-sm text-red-600">{errors.submit}</p>
        </div>
      )}

      {/* Preview (when editing) */}
      {isEditing && initialData && (
        <div className="p-4 bg-neutral-50 border border-neutral-200 rounded-lg">
          <h4 className="text-sm font-medium text-neutral-700 mb-2">原始記錄</h4>
          <div className="grid grid-cols-2 gap-2 text-sm">
            <div>
              <span className="text-neutral-500">建立時間：</span>
              <span className="font-medium">{formatDate(initialData.createdAt)}</span>
            </div>
            <div>
              <span className="text-neutral-500">最後更新：</span>
              <span className="font-medium">{formatDate(initialData.updatedAt)}</span>
            </div>
          </div>
        </div>
      )}

      {/* Action Buttons */}
      <div className="flex flex-col-reverse sm:flex-row gap-3 pt-4">
        <Button
          type="button"
          variant="secondary"
          onClick={onCancel}
          disabled={isLoading}
          className="flex-1 sm:flex-none"
        >
          取消
        </Button>
        <Button
          type="submit"
          variant="primary"
          loading={isLoading}
          className="flex-1 sm:flex-none"
          disabled={!formData.walletID || !formData.categoryID}
        >
          <TrendingDown className="w-4 h-4" />
          {isEditing ? '更新支出' : '新增支出'}
        </Button>
      </div>
    </form>
  )
}

export default ExpenseForm