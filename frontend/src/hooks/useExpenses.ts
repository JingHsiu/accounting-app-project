import { useQuery, useMutation, useQueryClient } from 'react-query'
import { expenseService } from '@/services'
import type { ExpenseRecord, CreateExpenseRequest, UpdateExpenseRequest, IncomeExpenseFilters } from '@/types'

// Query keys for React Query
export const expenseKeys = {
  all: ['expenses'] as const,
  lists: () => [...expenseKeys.all, 'list'] as const,
  list: (filters?: IncomeExpenseFilters) => [...expenseKeys.lists(), { filters }] as const,
  details: () => [...expenseKeys.all, 'detail'] as const,
  detail: (id: string) => [...expenseKeys.details(), id] as const,
}

// Hook for fetching expense records
export const useExpenses = (filters?: IncomeExpenseFilters, options?: {
  enabled?: boolean
  onSuccess?: (data: ExpenseRecord[]) => void
  onError?: (error: Error) => void
}) => {
  return useQuery(
    expenseKeys.list(filters),
    () => expenseService.getExpenses(filters, 'useExpenses'),
    {
      retry: 1,
      refetchOnWindowFocus: false,
      staleTime: 5 * 60 * 1000, // 5 minutes
      cacheTime: 10 * 60 * 1000, // 10 minutes
      onSuccess: options?.onSuccess,
      onError: options?.onError,
      enabled: options?.enabled,
    }
  )
}

// Hook for fetching a single expense record
export const useExpense = (expenseID: string, options?: {
  enabled?: boolean
  onSuccess?: (data: ExpenseRecord) => void
  onError?: (error: Error) => void
}) => {
  return useQuery(
    expenseKeys.detail(expenseID),
    async () => {
      const response = await expenseService.getExpense(expenseID)
      if (!response.success) {
        throw new Error(response.error)
      }
      return response.data!
    },
    {
      retry: 1,
      refetchOnWindowFocus: false,
      staleTime: 5 * 60 * 1000, // 5 minutes
      cacheTime: 10 * 60 * 1000, // 10 minutes
      onSuccess: options?.onSuccess,
      onError: options?.onError,
      enabled: !!expenseID && (options?.enabled !== false),
    }
  )
}

// Hook for creating expense records
export const useCreateExpense = (options?: {
  onSuccess?: (data: {id: string}, variables: CreateExpenseRequest) => void
  onError?: (error: Error, variables: CreateExpenseRequest) => void
  onMutate?: (variables: CreateExpenseRequest) => void
}) => {
  const queryClient = useQueryClient()
  
  return useMutation(
    (expense: CreateExpenseRequest) => expenseService.createExpense(expense),
    {
      onMutate: async (newExpense) => {
        // Cancel any outgoing refetches
        await queryClient.cancelQueries(expenseKeys.lists())
        
        options?.onMutate?.(newExpense)
      },
      onSuccess: (data, variables) => {
        // Invalidate and refetch expense queries
        queryClient.invalidateQueries(expenseKeys.lists())
        
        // Also invalidate wallet queries since balance might change
        queryClient.invalidateQueries(['wallets'])
        
        options?.onSuccess?.(data, variables)
      },
      onError: (error, variables) => {
        console.error('Failed to create expense record:', error)
        options?.onError?.(error as Error, variables)
      },
    }
  )
}

// Hook for updating expense records
export const useUpdateExpense = (options?: {
  onSuccess?: (data: ExpenseRecord, variables: {expenseID: string; updates: UpdateExpenseRequest}) => void
  onError?: (error: Error, variables: {expenseID: string; updates: UpdateExpenseRequest}) => void
  onMutate?: (variables: {expenseID: string; updates: UpdateExpenseRequest}) => void
}) => {
  const queryClient = useQueryClient()
  
  return useMutation(
    ({ expenseID, updates }: { expenseID: string; updates: UpdateExpenseRequest }) => 
      expenseService.updateExpense(expenseID, updates),
    {
      onMutate: async ({ expenseID, updates }) => {
        // Cancel any outgoing refetches
        await queryClient.cancelQueries(expenseKeys.detail(expenseID))
        await queryClient.cancelQueries(expenseKeys.lists())
        
        options?.onMutate?.({ expenseID, updates })
      },
      onSuccess: async (response, variables) => {
        if (response.success && response.data) {
          // Update the specific expense in cache
          queryClient.setQueryData(expenseKeys.detail(variables.expenseID), response.data)
          
          // Invalidate expense lists to refresh
          queryClient.invalidateQueries(expenseKeys.lists())
          
          // Also invalidate wallet queries since balance might change
          queryClient.invalidateQueries(['wallets'])
          
          options?.onSuccess?.(response.data, variables)
        } else {
          throw new Error(response.error || 'Failed to update expense record')
        }
      },
      onError: (error, variables) => {
        console.error('Failed to update expense record:', error)
        options?.onError?.(error as Error, variables)
      },
    }
  )
}

// Hook for deleting expense records
export const useDeleteExpense = (options?: {
  onSuccess?: (data: void, variables: string) => void
  onError?: (error: Error, variables: string) => void
  onMutate?: (variables: string) => void
}) => {
  const queryClient = useQueryClient()
  
  return useMutation(
    (expenseID: string) => expenseService.deleteExpense(expenseID),
    {
      onMutate: async (expenseID) => {
        // Cancel any outgoing refetches
        await queryClient.cancelQueries(expenseKeys.detail(expenseID))
        await queryClient.cancelQueries(expenseKeys.lists())
        
        options?.onMutate?.(expenseID)
      },
      onSuccess: async (response, variables) => {
        if (response.success) {
          // Remove the specific expense from cache
          queryClient.removeQueries(expenseKeys.detail(variables))
          
          // Invalidate expense lists to refresh
          queryClient.invalidateQueries(expenseKeys.lists())
          
          // Also invalidate wallet queries since balance might change
          queryClient.invalidateQueries(['wallets'])
          
          options?.onSuccess?.(undefined, variables)
        } else {
          throw new Error(response.error || 'Failed to delete expense record')
        }
      },
      onError: (error, variables) => {
        console.error('Failed to delete expense record:', error)
        options?.onError?.(error as Error, variables)
      },
    }
  )
}