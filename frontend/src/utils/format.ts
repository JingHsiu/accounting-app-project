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
 * Format date in various formats with proper error handling
 */
export const formatDate = (date: string | Date | undefined | null, format: 'short' | 'long' | 'time' = 'short'): string => {
  // Handle null/undefined input
  if (!date) {
    return '無日期'
  }
  
  const d = typeof date === 'string' ? new Date(date) : date
  
  // Validate date object
  if (!d || isNaN(d.getTime())) {
    console.warn('⚠️ Invalid date passed to formatDate:', date)
    return '無效日期'
  }
  
  try {
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
  } catch (error) {
    console.error('❌ formatDate error:', error, 'Input:', date)
    return '格式錯誤'
  }
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
 * Get currency subdivision (how many smaller units make up 1 major unit)
 */
export const getCurrencySubdivision = (currency: string): number => {
  const currencySubdivisions: Record<string, number> = {
    // Whole number currencies (no subdivision) - TWD as default
    TWD: 1,     // 1 台幣 = 1 台幣 (no cents, base unit)
    JPY: 1,     // 1 yen = 1 yen (no sen in practice)
    KRW: 1,     // 1 won = 1 won (no subdivision)
    VND: 1,     // 1 dong = 1 dong (no subdivision)
    
    // Decimal currencies (1 unit = 100 smaller units)
    USD: 100,   // 1 dollar = 100 cents
    EUR: 100,   // 1 euro = 100 cents
    GBP: 100,   // 1 pound = 100 pence
    CNY: 100,   // 1 yuan = 100 fen
  }
  
  return currencySubdivisions[currency] || 1 // Default to 1 for unknown currencies (like TWD)
}

/**
 * Convert display amount to backend storage format
 * Handles currency-specific subdivisions correctly
 */
export const convertToBackendAmount = (displayAmount: number, currency: string): number => {
  const subdivision = getCurrencySubdivision(currency)
  return Math.round(displayAmount * subdivision)
}

/**
 * Convert backend storage amount to display format
 * Handles currency-specific subdivisions correctly  
 */
export const convertFromBackendAmount = (backendAmount: number, currency: string): number => {
  const subdivision = getCurrencySubdivision(currency)
  return backendAmount / subdivision
}

/**
 * Validate and parse money input - defaults to TWD
 */
export const parseMoney = (input: string, currency: string = 'TWD'): Money | null => {
  const cleanInput = input.replace(/[^\d.-]/g, '')
  const amount = parseFloat(cleanInput)
  
  if (isNaN(amount)) return null
  
  // Convert display amount to storage format using TWD as base unit
  const storageAmount = convertToBackendAmount(amount, currency)
  
  return { amount: storageAmount, currency }
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