import React, { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from 'react-query'
import { walletService } from '@/services'
import type { CreateWalletRequest } from '@/services/walletService'
import { WalletType } from '@/types'

const DEMO_USER_ID = "demo-user-123"

const DebugWallets: React.FC = () => {
  const [logs, setLogs] = useState<string[]>([])
  const queryClient = useQueryClient()

  const addLog = (message: string) => {
    const timestamp = new Date().toLocaleTimeString()
    setLogs(prev => [...prev, `[${timestamp}] ${message}`])
  }

  // Queries
  const { data: walletsData, isLoading, error, refetch } = useQuery(
    ['wallets', DEMO_USER_ID],
    () => walletService.getWallets(DEMO_USER_ID),
    {
      onSuccess: (data) => {
        addLog(`✅ Query SUCCESS: ${data?.length} wallets loaded`)
        console.log('✅ Wallets query success:', data)
      },
      onError: (error) => {
        addLog(`❌ Query ERROR: ${error}`)
        console.error('❌ Wallets query error:', error)
      },
      refetchOnWindowFocus: false,
      staleTime: 0
    }
  )

  // Mutation
  const createWalletMutation = useMutation(
    (wallet: CreateWalletRequest) => walletService.createWallet(wallet),
    {
      onSuccess: (data) => {
        addLog(`✅ Create SUCCESS: ${data.id}`)
        queryClient.invalidateQueries(['wallets', DEMO_USER_ID])
        addLog(`🔄 Cache invalidated`)
      },
      onError: (error) => {
        addLog(`❌ Create ERROR: ${error}`)
      }
    }
  )

  const handleCreateWallet = () => {
    const walletData = {
      name: `Debug Test ${Date.now()}`,
      type: WalletType.CASH,
      currency: 'TWD',
      user_id: DEMO_USER_ID,
      initialBalance: 100
    }
    
    addLog(`🚀 Creating wallet: ${walletData.name}`)
    createWalletMutation.mutate(walletData)
  }

  const handleDirectApiTest = async () => {
    try {
      addLog(`🔍 Direct API test starting...`)
      const response = await fetch('/api/v1/wallets?userID=demo-user-123')
      const data = await response.json()
      addLog(`📋 Direct API result: ${data.data?.count} wallets`)
      console.log('Direct API result:', data)
    } catch (error) {
      addLog(`❌ Direct API error: ${error}`)
    }
  }

  const wallets = walletsData || []

  useEffect(() => {
    addLog(`🔄 Component mounted`)
  }, [])

  useEffect(() => {
    addLog(`📊 Wallets data updated: ${wallets.length} wallets`)
  }, [wallets.length])

  return (
    <div style={{ padding: '20px', fontFamily: 'monospace' }}>
      <h1>Wallet Debug Console</h1>
      
      <div style={{ marginBottom: '20px' }}>
        <button onClick={handleCreateWallet} disabled={createWalletMutation.isLoading}>
          {createWalletMutation.isLoading ? 'Creating...' : 'Create Test Wallet'}
        </button>
        
        <button onClick={() => refetch()} disabled={isLoading} style={{ marginLeft: '10px' }}>
          {isLoading ? 'Refetching...' : 'Refetch Wallets'}
        </button>
        
        <button onClick={handleDirectApiTest} style={{ marginLeft: '10px' }}>
          Direct API Test
        </button>
        
        <button onClick={() => setLogs([])} style={{ marginLeft: '10px' }}>
          Clear Logs
        </button>
      </div>

      <div style={{ display: 'flex', gap: '20px' }}>
        <div style={{ flex: 1 }}>
          <h2>Current State</h2>
          <pre style={{ background: '#f5f5f5', padding: '10px', fontSize: '12px' }}>
{JSON.stringify({
  isLoading,
  error: error?.toString(),
  walletsCount: wallets.length,
  queryState: queryClient.getQueryState(['wallets', DEMO_USER_ID]),
  mutationState: createWalletMutation
}, null, 2)}
          </pre>
          
          <h3>Wallets Data:</h3>
          <pre style={{ background: '#f0f8ff', padding: '10px', fontSize: '12px', maxHeight: '300px', overflow: 'auto' }}>
            {JSON.stringify(wallets, null, 2)}
          </pre>
        </div>
        
        <div style={{ flex: 1 }}>
          <h2>Debug Logs</h2>
          <div style={{ background: '#f5f5f5', padding: '10px', height: '500px', overflow: 'auto', fontSize: '12px' }}>
            {logs.map((log, index) => (
              <div key={index}>{log}</div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

export default DebugWallets