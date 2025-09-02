// React import not needed with jsx-runtime
// import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from 'react-query'
import Layout from '@/components/Layout'
import Dashboard from '@/pages/Dashboard'
import Wallets from '@/pages/Wallets'
import Transactions from '@/pages/Transactions'
import Categories from '@/pages/Categories'
import TestWallets from '@/TestWallets'
import DebugWallets from '@/DebugWallets'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <div className="min-h-screen bg-gradient-to-br from-primary-50 via-white to-primary-100">
          <Layout>
            <Routes>
              <Route path="/" element={<Navigate to="/dashboard" replace />} />
              <Route path="/dashboard" element={<Dashboard />} />
              <Route path="/wallets" element={<Wallets />} />
              <Route path="/transactions" element={<Transactions />} />
              <Route path="/categories" element={<Categories />} />
              <Route path="/test-wallets" element={<TestWallets />} />
              <Route path="/debug-wallets" element={<DebugWallets />} />
            </Routes>
          </Layout>
        </div>
      </Router>
    </QueryClientProvider>
  )
}

export default App