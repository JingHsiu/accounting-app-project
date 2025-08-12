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

const DEMO_USER_ID = "demo-user-123"

const Wallets: React.FC = () => {
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [editingWallet, setEditingWallet] = useState<string | null>(null)
  const [formData, setFormData] = useState({
    name: '',
    type: WalletType.CASH
  })

  const queryClient = useQueryClient()

  // Queries
  const { data: walletsData, isLoading } = useQuery(
    ['wallets', DEMO_USER_ID],
    () => walletService.getWallets(DEMO_USER_ID)
  )

  // Mutations
  const createWalletMutation = useMutation(
    (wallet: CreateWalletRequest) => walletService.createWallet(wallet),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['wallets'])
        setShowCreateModal(false)
        resetForm()
      }
    }
  )

  const deleteWalletMutation = useMutation(
    (walletID: string) => walletService.deleteWallet(walletID),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['wallets'])
      }
    }
  )

  const wallets = walletsData?.data || []

  const resetForm = () => {
    setFormData({ name: '', type: WalletType.CASH })
    setEditingWallet(null)
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    if (editingWallet) {
      // TODO: Implement update wallet
    } else {
      createWalletMutation.mutate({
        ...formData,
        userID: DEMO_USER_ID
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

  if (isLoading) {
    return (
      <div className="space-y-6">
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
          onClick={() => setShowCreateModal(true)}
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
                      onClick={() => {
                        setEditingWallet(wallet.id)
                        setFormData({ name: wallet.name, type: wallet.type })
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
            <Button 
              variant="primary"
              onClick={() => setShowCreateModal(true)}
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
        title={editingWallet ? '編輯錢包' : '新增錢包'}
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
              loading={createWalletMutation.isLoading}
              className="flex-1"
            >
              {editingWallet ? '更新' : '建立'}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}

export default Wallets