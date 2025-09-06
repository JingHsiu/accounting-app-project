import { useQuery, useMutation, useQueryClient } from 'react-query'
import { incomeService } from '@/services'
import type { IncomeRecord, CreateIncomeRequest, UpdateIncomeRequest, IncomeExpenseFilters } from '@/types'

// Query keys for React Query
export const incomeKeys = {
  all: ['incomes'] as const,
  lists: () => [...incomeKeys.all, 'list'] as const,
  list: (filters?: IncomeExpenseFilters) => [...incomeKeys.lists(), { filters }] as const,
  details: () => [...incomeKeys.all, 'detail'] as const,
  detail: (id: string) => [...incomeKeys.details(), id] as const,
}

// Hook for fetching income records
export const useIncomes = (filters?: IncomeExpenseFilters, options?: {
  enabled?: boolean
  onSuccess?: (data: IncomeRecord[]) => void
  onError?: (error: Error) => void
}) => {
  return useQuery(
    incomeKeys.list(filters),
    () => incomeService.getIncomes(filters, 'useIncomes'),
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

// Hook for fetching a single income record
export const useIncome = (incomeID: string, options?: {
  enabled?: boolean
  onSuccess?: (data: IncomeRecord) => void
  onError?: (error: Error) => void
}) => {
  return useQuery(
    incomeKeys.detail(incomeID),
    async () => {
      const response = await incomeService.getIncome(incomeID)
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
      enabled: !!incomeID && (options?.enabled !== false),
    }
  )
}

// Hook for creating income records
export const useCreateIncome = (options?: {
  onSuccess?: (data: {id: string}, variables: CreateIncomeRequest) => void
  onError?: (error: Error, variables: CreateIncomeRequest) => void
  onMutate?: (variables: CreateIncomeRequest) => void
}) => {
  const queryClient = useQueryClient()
  
  return useMutation(
    (income: CreateIncomeRequest) => incomeService.createIncome(income),
    {
      onMutate: async (newIncome) => {
        // Cancel any outgoing refetches
        await queryClient.cancelQueries(incomeKeys.lists())
        
        options?.onMutate?.(newIncome)
      },
      onSuccess: (data, variables) => {
        // Invalidate and refetch income queries
        queryClient.invalidateQueries(incomeKeys.lists())
        
        // Also invalidate wallet queries since balance might change
        queryClient.invalidateQueries(['wallets'])
        
        options?.onSuccess?.(data, variables)
      },
      onError: (error, variables) => {
        console.error('Failed to create income record:', error)
        options?.onError?.(error as Error, variables)
      },
    }
  )
}

// Hook for updating income records
export const useUpdateIncome = (options?: {
  onSuccess?: (data: IncomeRecord, variables: {incomeID: string; updates: UpdateIncomeRequest}) => void
  onError?: (error: Error, variables: {incomeID: string; updates: UpdateIncomeRequest}) => void
  onMutate?: (variables: {incomeID: string; updates: UpdateIncomeRequest}) => void
}) => {
  const queryClient = useQueryClient()
  
  return useMutation(
    ({ incomeID, updates }: { incomeID: string; updates: UpdateIncomeRequest }) => 
      incomeService.updateIncome(incomeID, updates),
    {
      onMutate: async ({ incomeID, updates }) => {
        // Cancel any outgoing refetches
        await queryClient.cancelQueries(incomeKeys.detail(incomeID))
        await queryClient.cancelQueries(incomeKeys.lists())
        
        options?.onMutate?.({ incomeID, updates })
      },
      onSuccess: async (response, variables) => {
        if (response.success && response.data) {
          // Update the specific income in cache
          queryClient.setQueryData(incomeKeys.detail(variables.incomeID), response.data)
          
          // Invalidate income lists to refresh
          queryClient.invalidateQueries(incomeKeys.lists())
          
          // Also invalidate wallet queries since balance might change
          queryClient.invalidateQueries(['wallets'])
          
          options?.onSuccess?.(response.data, variables)
        } else {
          throw new Error(response.error || 'Failed to update income record')
        }
      },
      onError: (error, variables) => {
        console.error('Failed to update income record:', error)
        options?.onError?.(error as Error, variables)
      },
    }
  )
}

// Hook for deleting income records
export const useDeleteIncome = (options?: {
  onSuccess?: (data: void, variables: string) => void
  onError?: (error: Error, variables: string) => void
  onMutate?: (variables: string) => void
}) => {
  const queryClient = useQueryClient()
  
  return useMutation(
    (incomeID: string) => incomeService.deleteIncome(incomeID),
    {
      onMutate: async (incomeID) => {
        // Cancel any outgoing refetches
        await queryClient.cancelQueries(incomeKeys.detail(incomeID))
        await queryClient.cancelQueries(incomeKeys.lists())
        
        options?.onMutate?.(incomeID)
      },
      onSuccess: async (response, variables) => {
        if (response.success) {
          // Remove the specific income from cache
          queryClient.removeQueries(incomeKeys.detail(variables))
          
          // Invalidate income lists to refresh
          queryClient.invalidateQueries(incomeKeys.lists())
          
          // Also invalidate wallet queries since balance might change
          queryClient.invalidateQueries(['wallets'])
          
          options?.onSuccess?.(undefined, variables)
        } else {
          throw new Error(response.error || 'Failed to delete income record')
        }
      },
      onError: (error, variables) => {
        console.error('Failed to delete income record:', error)
        options?.onError?.(error as Error, variables)
      },
    }
  )
}