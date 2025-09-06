import { useState, useCallback } from 'react'

export interface ValidationRule {
  required?: boolean
  min?: number
  max?: number
  minLength?: number
  maxLength?: number
  pattern?: RegExp
  custom?: (value: any) => string | null
  message?: string
}

export interface ValidationSchema {
  [key: string]: ValidationRule | ValidationRule[]
}

export interface ValidationErrors {
  [key: string]: string
}

export const useFormValidation = (schema: ValidationSchema) => {
  const [errors, setErrors] = useState<ValidationErrors>({})

  const validateField = useCallback((name: string, value: any): string | null => {
    const rules = schema[name]
    if (!rules) return null

    const ruleArray = Array.isArray(rules) ? rules : [rules]

    for (const rule of ruleArray) {
      // Required validation
      if (rule.required && (!value || (typeof value === 'string' && !value.trim()))) {
        return rule.message || `${name} 是必填欄位`
      }

      // Skip other validations if value is empty and not required
      if (!value && !rule.required) continue

      // Min value validation
      if (rule.min !== undefined && Number(value) < rule.min) {
        return rule.message || `${name} 不能小於 ${rule.min}`
      }

      // Max value validation
      if (rule.max !== undefined && Number(value) > rule.max) {
        return rule.message || `${name} 不能大於 ${rule.max}`
      }

      // Min length validation
      if (rule.minLength !== undefined && String(value).length < rule.minLength) {
        return rule.message || `${name} 長度不能小於 ${rule.minLength} 個字元`
      }

      // Max length validation
      if (rule.maxLength !== undefined && String(value).length > rule.maxLength) {
        return rule.message || `${name} 長度不能大於 ${rule.maxLength} 個字元`
      }

      // Pattern validation
      if (rule.pattern && !rule.pattern.test(String(value))) {
        return rule.message || `${name} 格式不正確`
      }

      // Custom validation
      if (rule.custom) {
        const customError = rule.custom(value)
        if (customError) return customError
      }
    }

    return null
  }, [schema])

  const validateForm = useCallback((data: Record<string, any>): boolean => {
    const newErrors: ValidationErrors = {}
    let hasErrors = false

    // Validate all fields in schema
    Object.keys(schema).forEach(fieldName => {
      const error = validateField(fieldName, data[fieldName])
      if (error) {
        newErrors[fieldName] = error
        hasErrors = true
      }
    })

    setErrors(newErrors)
    return !hasErrors
  }, [schema, validateField])

  const validateSingleField = useCallback((name: string, value: any) => {
    const error = validateField(name, value)
    setErrors(prev => ({
      ...prev,
      [name]: error || ''
    }))
    return !error
  }, [validateField])

  const clearErrors = useCallback(() => {
    setErrors({})
  }, [])

  const clearFieldError = useCallback((fieldName: string) => {
    setErrors(prev => {
      const newErrors = { ...prev }
      delete newErrors[fieldName]
      return newErrors
    })
  }, [])

  const setFieldError = useCallback((fieldName: string, error: string) => {
    setErrors(prev => ({
      ...prev,
      [fieldName]: error
    }))
  }, [])

  return {
    errors,
    validateForm,
    validateSingleField,
    clearErrors,
    clearFieldError,
    setFieldError,
    hasErrors: Object.keys(errors).length > 0
  }
}

// Common validation rules
export const validationRules = {
  required: (message?: string): ValidationRule => ({
    required: true,
    message: message || '此欄位為必填'
  }),

  minLength: (length: number, message?: string): ValidationRule => ({
    minLength: length,
    message: message || `至少需要 ${length} 個字元`
  }),

  maxLength: (length: number, message?: string): ValidationRule => ({
    maxLength: length,
    message: message || `不能超過 ${length} 個字元`
  }),

  min: (value: number, message?: string): ValidationRule => ({
    min: value,
    message: message || `不能小於 ${value}`
  }),

  max: (value: number, message?: string): ValidationRule => ({
    max: value,
    message: message || `不能大於 ${value}`
  }),

  email: (message?: string): ValidationRule => ({
    pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
    message: message || '請輸入有效的電子郵件地址'
  }),

  positiveNumber: (message?: string): ValidationRule => ({
    custom: (value) => {
      const num = Number(value)
      return isNaN(num) || num <= 0 ? (message || '必須是正數') : null
    }
  }),

  currency: (message?: string): ValidationRule => ({
    pattern: /^[A-Z]{3}$/,
    message: message || '貨幣代碼必須是3位大寫字母 (如: TWD)'
  }),

  date: (message?: string): ValidationRule => ({
    custom: (value) => {
      if (!value) return null
      const date = new Date(value)
      return isNaN(date.getTime()) ? (message || '請輸入有效的日期') : null
    }
  })
}