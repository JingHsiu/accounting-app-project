import React, { useState } from 'react'
import { useQuery, useQueryClient } from 'react-query'
import { walletService } from '@/services'

interface WalletDebugPanelProps {
  userID: string
  show?: boolean
}

const WalletDebugPanel: React.FC<WalletDebugPanelProps> = ({ userID, show = false }) => {
  const [isVisible, setIsVisible] = useState(show)
  const [testResults, setTestResults] = useState<any[]>([])
  const queryClient = useQueryClient()

  const { data: walletsData, isLoading, error, refetch } = useQuery(
    ['debug-wallets', userID],
    () => walletService.getWallets(userID, 'DebugPanel'),
    {
      enabled: isVisible,
      onSuccess: (data) => {
        setTestResults(prev => [...prev, {
          timestamp: new Date().toLocaleTimeString(),
          type: 'SUCCESS',
          data,
          length: data?.length
        }])
      },
      onError: (error) => {
        setTestResults(prev => [...prev, {
          timestamp: new Date().toLocaleTimeString(),
          type: 'ERROR',
          error: error instanceof Error ? error.message : error
        }])
      }
    }
  )

  const runDirectFetch = async () => {
    try {
      const response = await fetch(`/api/v1/wallets?userID=${userID}`)
      const data = await response.json()
      setTestResults(prev => [...prev, {
        timestamp: new Date().toLocaleTimeString(),
        type: 'DIRECT_API',
        response: data,
        status: response.status
      }])
    } catch (error) {
      setTestResults(prev => [...prev, {
        timestamp: new Date().toLocaleTimeString(),
        type: 'DIRECT_API_ERROR',
        error: error instanceof Error ? error.message : error
      }])
    }
  }

  const clearCache = () => {
    queryClient.clear()
    setTestResults(prev => [...prev, {
      timestamp: new Date().toLocaleTimeString(),
      type: 'CACHE_CLEARED',
      message: 'All React Query cache cleared'
    }])
  }

  if (!isVisible) {
    return (
      <div className="fixed bottom-4 right-4 z-50">
        <button
          onClick={() => setIsVisible(true)}
          className="bg-blue-500 hover:bg-blue-600 text-white px-3 py-2 rounded-lg text-sm"
        >
          🔧 Debug Panel
        </button>
      </div>
    )
  }

  return (
    <div className="fixed bottom-4 right-4 z-50 w-96 max-h-96 bg-white border border-gray-300 rounded-lg shadow-lg overflow-hidden">
      <div className="bg-blue-500 text-white px-3 py-2 flex justify-between items-center">
        <span className="font-semibold">🔧 Wallet Debug Panel</span>
        <button 
          onClick={() => setIsVisible(false)}
          className="text-white hover:text-gray-200"
        >
          ✕
        </button>
      </div>
      
      <div className="p-3">
        <div className="flex gap-2 mb-3">
          <button
            onClick={() => refetch()}
            disabled={isLoading}
            className="bg-green-500 hover:bg-green-600 text-white px-2 py-1 rounded text-xs"
          >
            {isLoading ? 'Loading...' : 'React Query'}
          </button>
          <button
            onClick={runDirectFetch}
            className="bg-orange-500 hover:bg-orange-600 text-white px-2 py-1 rounded text-xs"
          >
            Direct API
          </button>
          <button
            onClick={clearCache}
            className="bg-red-500 hover:bg-red-600 text-white px-2 py-1 rounded text-xs"
          >
            Clear Cache
          </button>
          <button
            onClick={() => setTestResults([])}
            className="bg-gray-500 hover:bg-gray-600 text-white px-2 py-1 rounded text-xs"
          >
            Clear Logs
          </button>
        </div>

        <div className="text-xs mb-2">
          <strong>Current State:</strong>
          <div>Loading: {String(isLoading)}</div>
          <div>Error: {error ? String(error) : 'None'}</div>
          <div>Data Length: {walletsData?.length || 0}</div>
        </div>

        <div className="text-xs">
          <strong>Test Results:</strong>
          <div className="max-h-32 overflow-y-auto bg-gray-100 p-2 rounded mt-1">
            {testResults.length === 0 ? (
              <div className="text-gray-500">No tests run yet</div>
            ) : (
              testResults.slice(-10).reverse().map((result, index) => (
                <div key={index} className="mb-1 pb-1 border-b border-gray-200 last:border-b-0">
                  <div className="font-semibold">
                    [{result.timestamp}] {result.type}
                  </div>
                  {result.type === 'SUCCESS' && (
                    <div>✅ {result.length} wallets loaded</div>
                  )}
                  {result.type === 'ERROR' && (
                    <div className="text-red-600">❌ {result.error}</div>
                  )}
                  {result.type === 'DIRECT_API' && (
                    <div>📡 Status: {result.status}, Count: {result.response?.data?.count || 0}</div>
                  )}
                  {result.type === 'DIRECT_API_ERROR' && (
                    <div className="text-red-600">🚨 {result.error}</div>
                  )}
                  {result.type === 'CACHE_CLEARED' && (
                    <div>🗑️ {result.message}</div>
                  )}
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default WalletDebugPanel