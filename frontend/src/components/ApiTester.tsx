import React, { useState } from 'react'
import { Button } from '@/components/ui'
import { walletService, categoryService } from '@/services'
import { apiDebugger } from '@/utils/apiDebug'

const DEMO_USER_ID = "demo-user-123"

const ApiTester: React.FC = () => {
  const [testing, setTesting] = useState(false)
  const [results, setResults] = useState<string[]>([])

  const addResult = (message: string) => {
    setResults(prev => [...prev, `[${new Date().toLocaleTimeString()}] ${message}`])
  }

  const runTests = async () => {
    setTesting(true)
    setResults([])
    addResult('🚀 Starting API tests...')

    try {
      // Test wallets endpoint
      addResult('📋 Testing wallets endpoint...')
      const wallets = await walletService.getWallets(DEMO_USER_ID, 'ApiTester')
      addResult(`✅ Wallets loaded: ${wallets.length} found`)

      // Test categories endpoint
      addResult('📁 Testing categories endpoint...')
      try {
        const categories = await categoryService.getCategories()
        addResult(`✅ Categories loaded: ${categories.length} found`)
      } catch (error: any) {
        addResult(`❌ Categories error: ${error.message}`)
      }

      addResult('🎯 Check console for detailed API debugging logs')
      addResult('🔍 Run apiDebugger.detectPortIssues() in console to check ports')
      addResult('📊 Run apiDebugger.compareComponents("Dashboard", "WalletsPage") to compare')

    } catch (error: any) {
      addResult(`❌ Test failed: ${error.message}`)
    } finally {
      setTesting(false)
    }
  }

  const clearLogs = () => {
    apiDebugger.clearLogs()
    addResult('🧹 API debug logs cleared')
  }

  const showDebugInfo = () => {
    console.group('🔍 API Debug Information')
    console.log('📊 Current logs:', apiDebugger.getLogs())
    console.log('🎯 Port detection:')
    apiDebugger.detectPortIssues()
    console.log('📱 Component comparison:')
    apiDebugger.compareComponents('Dashboard', 'WalletsPage')
    console.groupEnd()
    
    addResult('📊 Debug info printed to console')
  }

  return (
    <div className="space-y-4 p-6 bg-white/50 rounded-lg border border-neutral-200">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-neutral-800">API Tester & Debugger</h3>
        <div className="flex gap-2">
          <Button 
            size="sm" 
            variant="secondary" 
            onClick={clearLogs}
            disabled={testing}
          >
            Clear Logs
          </Button>
          <Button 
            size="sm" 
            variant="outline" 
            onClick={showDebugInfo}
            disabled={testing}
          >
            Debug Info
          </Button>
          <Button 
            size="sm" 
            variant="primary" 
            onClick={runTests}
            disabled={testing}
            loading={testing}
          >
            Run Tests
          </Button>
        </div>
      </div>

      {results.length > 0 && (
        <div className="bg-neutral-50 rounded-lg p-4">
          <h4 className="text-sm font-medium text-neutral-700 mb-2">Test Results:</h4>
          <div className="space-y-1 text-sm font-mono">
            {results.map((result, index) => (
              <div 
                key={index} 
                className={`${
                  result.includes('✅') ? 'text-green-600' :
                  result.includes('❌') ? 'text-red-600' :
                  result.includes('⚠️') ? 'text-yellow-600' :
                  'text-neutral-600'
                }`}
              >
                {result}
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="text-sm text-neutral-500 space-y-1">
        <p><strong>Usage:</strong></p>
        <p>• Click "Run Tests" to test API endpoints</p>
        <p>• Check browser console for detailed logs</p>
        <p>• Use <code>apiDebugger</code> in console for advanced debugging</p>
        <p>• <code>apiDebugger.detectPortIssues()</code> - Check port problems</p>
        <p>• <code>apiDebugger.compareComponents()</code> - Compare API calls</p>
      </div>
    </div>
  )
}

export default ApiTester