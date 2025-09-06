import React from 'react'
import { 
  TrendingUp, 
  TrendingDown, 
  Edit, 
  Trash2, 
  Calendar,
  Wallet as WalletIcon,
  Tag
} from 'lucide-react'
import { Card, CardContent, Button } from '@/components/ui'
import { formatMoney, formatDate } from '@/utils/format'
import type { IncomeRecord, ExpenseRecord, Wallet, Category } from '@/types'

type TransactionRecord = IncomeRecord | ExpenseRecord
type TransactionType = 'income' | 'expense'

interface TransactionCardProps {
  transaction: TransactionRecord
  type: TransactionType
  wallet?: Wallet
  category?: Category
  onEdit?: (transaction: TransactionRecord) => void
  onDelete?: (transaction: TransactionRecord) => void
  isLoading?: boolean
}

const TransactionCard: React.FC<TransactionCardProps> = ({
  transaction,
  type,
  wallet,
  category,
  onEdit,
  onDelete,
  isLoading = false
}) => {
  const isIncome = type === 'income'
  const iconColor = isIncome ? 'text-green-600' : 'text-red-600'
  const bgColor = isIncome ? 'bg-green-100' : 'bg-red-100'
  const amountColor = isIncome ? 'text-green-700' : 'text-red-700'
  
  const Icon = isIncome ? TrendingUp : TrendingDown
  
  return (
    <Card className="card-hover border-l-4" style={{ borderLeftColor: isIncome ? '#059669' : '#dc2626' }}>
      <CardContent>
        <div className="flex items-start justify-between">
          <div className="flex items-start gap-4 flex-1">
            {/* Transaction Icon */}
            <div className={`p-3 rounded-xl ${bgColor} ${iconColor} flex-shrink-0`}>
              <Icon className="w-5 h-5" />
            </div>
            
            <div className="flex-1 min-w-0">
              {/* Description and Amount */}
              <div className="flex items-start justify-between gap-4 mb-3">
                <div>
                  <h3 className="font-medium text-neutral-800 text-base">
                    {transaction.description}
                  </h3>
                  <div className={`text-2xl font-bold ${amountColor} mt-1`}>
                    {isIncome ? '+' : '-'}{formatMoney(transaction.amount)}
                  </div>
                </div>
              </div>

              {/* Details Grid */}
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 text-sm text-neutral-600">
                <div className="flex items-center gap-2">
                  <Calendar className="w-4 h-4" />
                  <span>{formatDate(transaction.date)}</span>
                </div>
                
                {wallet && (
                  <div className="flex items-center gap-2">
                    <WalletIcon className="w-4 h-4" />
                    <span className="truncate">{wallet.name}</span>
                  </div>
                )}
                
                {category && (
                  <div className="flex items-center gap-2">
                    <Tag className="w-4 h-4" />
                    <span className="truncate">{category.name}</span>
                  </div>
                )}
              </div>

              {/* Timestamps */}
              <div className="flex items-center gap-4 mt-3 pt-3 border-t border-neutral-100 text-xs text-neutral-500">
                <span>建立：{formatDate(transaction.createdAt)}</span>
                {transaction.updatedAt !== transaction.createdAt && (
                  <span>更新：{formatDate(transaction.updatedAt)}</span>
                )}
              </div>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex flex-col gap-1 ml-2">
            {onEdit && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => onEdit(transaction)}
                disabled={isLoading}
                className="p-2 hover:bg-amber-50 hover:text-amber-600"
                title="編輯"
              >
                <Edit className="w-4 h-4" />
              </Button>
            )}
            {onDelete && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => onDelete(transaction)}
                disabled={isLoading}
                className="p-2 hover:bg-red-50 hover:text-red-600"
                title="刪除"
              >
                <Trash2 className="w-4 h-4" />
              </Button>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

export default TransactionCard