/**
 * 統一的 API 字段轉換工具
 * 處理前端 camelCase 與後端 snake_case 的轉換
 */

// 基礎轉換函數
export const toSnakeCase = (str: string): string => {
  return str.replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`)
}

export const toCamelCase = (str: string): string => {
  return str.replace(/_([a-z])/g, (match, letter) => letter.toUpperCase())
}

// 深度轉換對象的鍵名
export const transformKeys = (obj: any, keyTransform: (key: string) => string): any => {
  if (obj === null || obj === undefined || typeof obj !== 'object') {
    return obj
  }
  
  if (Array.isArray(obj)) {
    return obj.map(item => transformKeys(item, keyTransform))
  }
  
  const result: any = {}
  Object.keys(obj).forEach(key => {
    const transformedKey = keyTransform(key)
    const value = obj[key]
    
    if (value && typeof value === 'object' && !Array.isArray(value)) {
      result[transformedKey] = transformKeys(value, keyTransform)
    } else if (Array.isArray(value)) {
      result[transformedKey] = value.map(item => 
        typeof item === 'object' ? transformKeys(item, keyTransform) : item
      )
    } else {
      result[transformedKey] = value
    }
  })
  
  return result
}

// 專用轉換函數
export const toBackendFormat = <T>(frontendData: T): any => {
  return transformKeys(frontendData, toSnakeCase)
}

export const toFrontendFormat = <T>(backendData: any): T => {
  return transformKeys(backendData, toCamelCase)
}

// 特殊字段映射（處理不規則轉換）
const FIELD_MAPPINGS = {
  frontend_to_backend: {
    'categoryID': 'subcategory_id',
    'walletID': 'wallet_id',
    'createdAt': 'created_at',
    'updatedAt': 'updated_at'
  },
  backend_to_frontend: {
    'subcategory_id': 'categoryID',
    'wallet_id': 'walletID', 
    'created_at': 'createdAt',
    'updated_at': 'updatedAt'
  }
} as const

// 帶特殊映射的轉換
export const toBackendFormatWithMapping = <T>(frontendData: T): any => {
  const basicTransform = transformKeys(frontendData, toSnakeCase)
  
  // 應用特殊映射
  const result = { ...basicTransform }
  Object.entries(FIELD_MAPPINGS.frontend_to_backend).forEach(([frontend, backend]) => {
    if (result[frontend] !== undefined) {
      result[backend] = result[frontend]
      delete result[frontend]
    }
  })
  
  return result
}

export const toFrontendFormatWithMapping = <T>(backendData: any): T => {
  // 先應用特殊映射
  const mappedData = { ...backendData }
  Object.entries(FIELD_MAPPINGS.backend_to_frontend).forEach(([backend, frontend]) => {
    if (mappedData[backend] !== undefined) {
      mappedData[frontend] = mappedData[backend]
      delete mappedData[backend]
    }
  })
  
  // 再應用基礎轉換
  return transformKeys(mappedData, toCamelCase)
}

// 泛型 API 響應處理
export const handleApiResponse = <T>(response: any): T[] | T => {
  if (!response.success || !response.data) {
    throw new Error(response.error || 'API request failed')
  }
  
  const data = response.data
  
  // 處理數組格式：{data: [...], count: number} 或直接數組
  if (Array.isArray(data)) {
    return data.map(item => toFrontendFormatWithMapping<T>(item))
  } else if (data && typeof data === 'object' && Array.isArray(data.data)) {
    return data.data.map((item: any) => toFrontendFormatWithMapping<T>(item))
  } else if (data && typeof data === 'object') {
    // 單個對象：{data: object} 或直接對象
    const actualData = data.data || data
    return toFrontendFormatWithMapping<T>(actualData)
  }
  
  throw new Error('Invalid API response format')
}

// 類型安全的轉換器生成器
export const createApiTransformer = <TFrontend, TBackend>() => ({
  toBackend: (data: TFrontend): TBackend => toBackendFormatWithMapping(data),
  fromBackend: (data: any): TFrontend => toFrontendFormatWithMapping<TFrontend>(data),
  handleResponse: (response: any): TFrontend[] | TFrontend => handleApiResponse<TFrontend>(response)
})

// 預定義的轉換器
export const expenseTransformer = createApiTransformer<
  import('@/types').ExpenseRecord,
  any
>()

export const incomeTransformer = createApiTransformer<
  import('@/types').IncomeRecord, 
  any
>()

export const categoryTransformer = createApiTransformer<
  import('@/types').Category,
  any
>()

export const walletTransformer = createApiTransformer<
  import('@/types').Wallet,
  any
>()