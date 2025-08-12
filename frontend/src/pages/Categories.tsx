import React, { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from 'react-query'
import { 
  Tags, 
  Plus, 
  Edit, 
  Trash2,
  ArrowUpRight,
  ArrowDownRight
} from 'lucide-react'
import { Card, CardHeader, CardTitle, CardContent, Button, Modal, Input, Select } from '@/components/ui'
import { categoryService } from '@/services'
// import { getCategoryTypeDisplayName } from '@/utils/format' // TODO: Use when needed
import type { CreateCategoryRequest } from '@/services/categoryService'
import { CategoryType } from '@/types'

const Categories: React.FC = () => {
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [editingCategory, setEditingCategory] = useState<string | null>(null)
  const [selectedType, setSelectedType] = useState<CategoryType | ''>('')
  const [formData, setFormData] = useState({
    name: '',
    type: CategoryType.EXPENSE,
    color: '#a855f7',
    icon: 'Tags'
  })

  const queryClient = useQueryClient()

  // Queries
  const { data: categoriesData, isLoading } = useQuery(
    ['categories', selectedType],
    () => categoryService.getCategories(selectedType || undefined)
  )

  // Mutations
  const createCategoryMutation = useMutation(
    (category: CreateCategoryRequest) => categoryService.createCategory(category),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['categories'])
        setShowCreateModal(false)
        resetForm()
      }
    }
  )

  const deleteCategoryMutation = useMutation(
    (categoryID: string) => categoryService.deleteCategory(categoryID),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['categories'])
      }
    }
  )

  const categories = categoriesData?.data || []

  const resetForm = () => {
    setFormData({
      name: '',
      type: CategoryType.EXPENSE,
      color: '#a855f7',
      icon: 'Tags'
    })
    setEditingCategory(null)
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    if (editingCategory) {
      // TODO: Implement update category
    } else {
      createCategoryMutation.mutate(formData)
    }
  }

  // TODO: Use when needed for category display
  // const getCategoryIcon = (type: CategoryType) => {
  //   return type === 'INCOME' 
  //     ? <ArrowUpRight className="w-5 h-5 text-secondary-600" />
  //     : <ArrowDownRight className="w-5 h-5 text-accent-600" />
  // }

  const categoryTypeOptions = [
    { value: 'INCOME', label: '收入' },
    { value: 'EXPENSE', label: '支出' }
  ]

  const filterOptions = [
    { value: '', label: '所有類別' },
    { value: 'INCOME', label: '收入類別' },
    { value: 'EXPENSE', label: '支出類別' }
  ]

  // Group categories by type
  const incomeCategories = categories.filter(c => c.type === CategoryType.INCOME)
  const expenseCategories = categories.filter(c => c.type === CategoryType.EXPENSE)

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {[...Array(4)].map((_, i) => (
            <Card key={i} glass className="animate-pulse">
              <div className="h-48 bg-primary-200/20 rounded" />
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
          <h1 className="text-3xl font-bold text-gradient-primary">類別管理</h1>
          <p className="text-neutral-600 mt-1">管理您的收入和支出類別</p>
        </div>
        <Button 
          variant="primary"
          onClick={() => setShowCreateModal(true)}
        >
          <Plus className="w-4 h-4" />
          新增類別
        </Button>
      </div>

      {/* Filter */}
      <Card glass>
        <CardContent>
          <div className="flex items-center gap-4">
            <label className="text-sm font-medium text-neutral-700">篩選:</label>
            <Select
              value={selectedType}
              onChange={(e) => setSelectedType(e.target.value as CategoryType | '')}
              options={filterOptions}
            />
          </div>
        </CardContent>
      </Card>

      {/* Categories Grid */}
      {categories.length > 0 ? (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Income Categories */}
          {(!selectedType || selectedType === CategoryType.INCOME) && incomeCategories.length > 0 && (
            <Card glass>
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-secondary-600">
                  <ArrowUpRight className="w-5 h-5" />
                  收入類別
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {incomeCategories.map((category) => (
                    <div key={category.id} className="flex items-center justify-between p-3 hover:bg-secondary-50/50 rounded-lg transition-colors border border-secondary-100">
                      <div className="flex items-center gap-3">
                        <div 
                          className="w-4 h-4 rounded-full"
                          style={{ backgroundColor: category.color }}
                        />
                        <span className="font-medium text-neutral-800">{category.name}</span>
                      </div>
                      <div className="flex gap-1">
                        <Button 
                          variant="ghost" 
                          size="sm"
                          onClick={() => {
                            setEditingCategory(category.id)
                            setFormData({
                              name: category.name,
                              type: category.type,
                              color: category.color,
                              icon: category.icon
                            })
                            setShowCreateModal(true)
                          }}
                        >
                          <Edit className="w-4 h-4" />
                        </Button>
                        <Button 
                          variant="ghost" 
                          size="sm"
                          onClick={() => deleteCategoryMutation.mutate(category.id)}
                        >
                          <Trash2 className="w-4 h-4 text-accent-600" />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}

          {/* Expense Categories */}
          {(!selectedType || selectedType === CategoryType.EXPENSE) && expenseCategories.length > 0 && (
            <Card glass>
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-accent-600">
                  <ArrowDownRight className="w-5 h-5" />
                  支出類別
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {expenseCategories.map((category) => (
                    <div key={category.id} className="flex items-center justify-between p-3 hover:bg-accent-50/50 rounded-lg transition-colors border border-accent-100">
                      <div className="flex items-center gap-3">
                        <div 
                          className="w-4 h-4 rounded-full"
                          style={{ backgroundColor: category.color }}
                        />
                        <span className="font-medium text-neutral-800">{category.name}</span>
                      </div>
                      <div className="flex gap-1">
                        <Button 
                          variant="ghost" 
                          size="sm"
                          onClick={() => {
                            setEditingCategory(category.id)
                            setFormData({
                              name: category.name,
                              type: category.type,
                              color: category.color,
                              icon: category.icon
                            })
                            setShowCreateModal(true)
                          }}
                        >
                          <Edit className="w-4 h-4" />
                        </Button>
                        <Button 
                          variant="ghost" 
                          size="sm"
                          onClick={() => deleteCategoryMutation.mutate(category.id)}
                        >
                          <Trash2 className="w-4 h-4 text-accent-600" />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      ) : (
        <Card glass className="text-center py-12">
          <CardContent>
            <Tags className="w-16 h-16 text-neutral-300 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-neutral-700 mb-2">尚無類別</h3>
            <p className="text-neutral-500 mb-6">建立類別來組織您的收支記錄</p>
            <Button 
              variant="primary"
              onClick={() => setShowCreateModal(true)}
            >
              <Plus className="w-4 h-4" />
              建立類別
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
        title={editingCategory ? '編輯類別' : '新增類別'}
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="類別名稱"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="輸入類別名稱"
            required
          />
          
          <Select
            label="類別類型"
            value={formData.type}
            onChange={(e) => setFormData({ ...formData, type: e.target.value as CategoryType })}
            options={categoryTypeOptions}
          />
          
          <div>
            <label className="block text-sm font-medium text-neutral-700 mb-2">
              顏色
            </label>
            <div className="flex items-center gap-2">
              <input
                type="color"
                value={formData.color}
                onChange={(e) => setFormData({ ...formData, color: e.target.value })}
                className="w-12 h-10 rounded border border-primary-200"
              />
              <Input
                value={formData.color}
                onChange={(e) => setFormData({ ...formData, color: e.target.value })}
                placeholder="#a855f7"
                className="flex-1"
              />
            </div>
          </div>
          
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
              loading={createCategoryMutation.isLoading}
              className="flex-1"
            >
              {editingCategory ? '更新' : '建立'}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}

export default Categories