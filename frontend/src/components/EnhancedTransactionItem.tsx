import React from 'react'
import { cn } from '@/lib/utils'
import { Card, CardContent } from '@/components/ui'
import { Button } from '@/components/ui'
import { 
  ArrowUpRight, 
  ArrowDownLeft, 
  Calendar, 
  X, 
  Wallet
} from 'lucide-react'
import { formatDate } from '@/utils/format'
import type { ExpenseRecord, IncomeRecord } from '@/types'

// Transaction Item Props
interface EnhancedTransactionItemProps {
  transaction: ExpenseRecord | IncomeRecord;
  walletName?: string;
  categoryName?: string;
  onDelete?: (id: string) => void;
  className?: string;
}

export const EnhancedTransactionItem: React.FC<EnhancedTransactionItemProps> = ({ 
  transaction, 
  walletName = 'Unknown Wallet',
  categoryName = 'Unknown Category',
  onDelete,
  className 
}) => {
  // Check if transaction has type property (from combined view)
  // Both ExpenseRecord and IncomeRecord have identical structures, so we rely on the type property
  const isExpense = (transaction as any).type === 'expense';
  const amount = Math.abs(transaction.amount.amount);
  const symbol = transaction.amount.currency === 'TWD' ? 'NT$' : '$';

  // Format amount based on currency
  const formattedAmount = transaction.amount.currency === 'TWD' 
    ? amount.toFixed(0) 
    : amount.toFixed(2);

  return (
    <Card className={cn("hover:shadow-md transition-shadow group", className)}>
      <CardContent className="p-4">
        <div className="flex items-start justify-between">
          <div className="flex items-start gap-3 flex-1 min-w-0">
            {/* Transaction Type Icon */}
            <div className={cn(
              "w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0",
              isExpense ? "bg-red-100 text-red-600" : "bg-green-100 text-green-600"
            )}>
              {isExpense ? <ArrowDownLeft className="w-5 h-5" /> : <ArrowUpRight className="w-5 h-5" />}
            </div>
            
            <div className="flex-1 min-w-0">
              {/* Main Info: Wallet • Category • Amount */}
              <div className="flex items-center justify-between mb-1">
                <div className="flex items-center gap-2 min-w-0 flex-1">
                  <div className="flex items-center gap-1 text-sm font-medium truncate">
                    <Wallet className="w-3 h-3 text-muted-foreground flex-shrink-0" />
                    <span className="truncate">{walletName}</span>
                  </div>
                  <span className="text-muted-foreground flex-shrink-0">•</span>
                  <span className="text-sm text-muted-foreground truncate">{categoryName}</span>
                </div>
                <span className={cn(
                  "font-semibold text-sm ml-3 flex-shrink-0",
                  isExpense ? "text-red-600" : "text-green-600"
                )}>
                  {isExpense ? "-" : "+"}{symbol}{formattedAmount}
                </span>
              </div>
              
              {/* Sub Info: Description • Date • Currency */}
              <div className="flex items-center justify-between text-xs text-muted-foreground">
                <div className="flex items-center gap-2 min-w-0 flex-1">
                  {transaction.description && (
                    <>
                      <span className="truncate">{transaction.description}</span>
                      <span className="flex-shrink-0">•</span>
                    </>
                  )}
                  <span className="flex items-center gap-1 flex-shrink-0">
                    <Calendar className="w-3 h-3" />
                    {formatDate(transaction.date, 'short')}
                  </span>
                </div>
                <div className="flex items-center gap-2 ml-2 flex-shrink-0">
                  <span className={cn(
                    "inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium border",
                    "bg-background border-border"
                  )}>
                    {transaction.amount.currency}
                  </span>
                </div>
              </div>
            </div>
          </div>
          
          {/* Delete Button */}
          {onDelete && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onDelete(transaction.id)}
              className="ml-2 flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity h-8 w-8 p-0 hover:bg-red-100 hover:text-red-600"
            >
              <X className="w-4 h-4" />
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
};

// Transaction List Component
interface EnhancedTransactionListProps {
  transactions: (ExpenseRecord | IncomeRecord)[];
  wallets?: Array<{ id: string; name: string; }>;
  categories?: Array<{ id: string; name: string; }>;
  onDeleteTransaction?: (id: string) => void;
  className?: string;
}

export const EnhancedTransactionList: React.FC<EnhancedTransactionListProps> = ({
  transactions,
  wallets = [],
  categories = [],
  onDeleteTransaction,
  className
}) => {
  const getWalletName = (walletId: string) => {
    return wallets.find(w => w.id === walletId)?.name || 'Unknown Wallet';
  };

  const getCategoryName = (categoryId: string) => {
    return categories.find(c => c.id === categoryId)?.name || 'Unknown Category';
  };

  if (transactions.length === 0) {
    return (
      <Card className={className}>
        <CardContent className="p-8 text-center">
          <div className="w-12 h-12 bg-muted rounded-lg flex items-center justify-center mx-auto mb-4">
            <Wallet className="w-6 h-6 text-muted-foreground" />
          </div>
          <h3 className="font-medium mb-2">No transactions found</h3>
          <p className="text-muted-foreground">
            Add your first transaction to get started
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className={cn("space-y-4", className)}>
      {transactions.map((transaction) => (
        <EnhancedTransactionItem
          key={transaction.id}
          transaction={transaction}
          walletName={getWalletName(transaction.walletID)}
          categoryName={getCategoryName(transaction.categoryID)}
          onDelete={onDeleteTransaction}
        />
      ))}
    </div>
  );
};

export default EnhancedTransactionItem;