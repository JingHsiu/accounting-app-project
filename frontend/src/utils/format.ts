import type { Money } from '@/types'

/**
 * Format money amount with currency
 */
export const formatMoney = (money: Money): string => {
  if (!money) return 'NT$ 0'
  
  const { amount, currency } = money
  const formattedAmount = new Intl.NumberFormat('zh-TW', {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(Math.abs(amount))

  // Handle different currencies
  const currencySymbols: Record<string, string> = {
    TWD: 'NT$',
    USD: '$',
    JPY: '¥',
    EUR: '€',
    CNY: '¥',
  }

  const symbol = currencySymbols[currency] || currency
  const sign = amount < 0 ? '-' : ''
  
  return `${sign}${symbol} ${formattedAmount}`
}

/**
 * Format date in various formats
 */
export const formatDate = (date: string | Date, format: 'short' | 'long' | 'time' = 'short'): string => {
  const d = typeof date === 'string' ? new Date(date) : date
  
  if (format === 'time') {
    return d.toLocaleTimeString('zh-TW', {
      hour: '2-digit',
      minute: '2-digit'
    })
  }
  
  if (format === 'long') {
    return d.toLocaleDateString('zh-TW', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      weekday: 'long'
    })
  }
  
  // short format
  return d.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

/**
 * Format percentage
 */
export const formatPercentage = (value: number, decimals: number = 1): string => {
  return `${value.toFixed(decimals)}%`
}

/**
 * Format large numbers with K, M suffixes
 */
export const formatCompactNumber = (num: number): string => {
  const formatter = new Intl.NumberFormat('zh-TW', {
    notation: 'compact',
    compactDisplay: 'short'
  })
  
  return formatter.format(num)
}

/**
 * Get relative time (e.g., "2 hours ago")
 */
export const formatRelativeTime = (date: string | Date): string => {
  const d = typeof date === 'string' ? new Date(date) : date
  const now = new Date()
  const diffMs = now.getTime() - d.getTime()
  const diffSecs = Math.floor(diffMs / 1000)
  const diffMins = Math.floor(diffSecs / 60)
  const diffHours = Math.floor(diffMins / 60)
  const diffDays = Math.floor(diffHours / 24)

  if (diffSecs < 60) return '剛剛'
  if (diffMins < 60) return `${diffMins} 分鐘前`
  if (diffHours < 24) return `${diffHours} 小時前`
  if (diffDays < 7) return `${diffDays} 天前`
  
  return formatDate(d)
}

/**
 * Validate and parse money input
 */
export const parseMoney = (input: string, currency: string = 'TWD'): Money | null => {
  const cleanInput = input.replace(/[^\d.-]/g, '')
  const amount = parseFloat(cleanInput)
  
  if (isNaN(amount)) return null
  
  return { amount, currency }
}

/**
 * Get wallet type display name
 */
export const getWalletTypeDisplayName = (type: string): string => {
  const typeNames: Record<string, string> = {
    CASH: '現金',
    BANK: '銀行帳戶',
    CREDIT: '信用卡',
    INVESTMENT: '投資帳戶'
  }
  
  return typeNames[type] || type
}

/**
 * Get category type display name
 */
export const getCategoryTypeDisplayName = (type: string): string => {
  const typeNames: Record<string, string> = {
    INCOME: '收入',
    EXPENSE: '支出'
  }
  
  return typeNames[type] || type
}