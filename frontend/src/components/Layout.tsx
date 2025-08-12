import React, { useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { 
  LayoutDashboard, 
  Wallet, 
  Receipt, 
  Tags, 
  Menu, 
  X,
  Sparkles
} from 'lucide-react'

interface LayoutProps {
  children: React.ReactNode
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const location = useLocation()

  const navigation = [
    { name: '儀表板', href: '/dashboard', icon: LayoutDashboard },
    { name: '錢包管理', href: '/wallets', icon: Wallet },
    { name: '交易記錄', href: '/transactions', icon: Receipt },
    { name: '類別管理', href: '/categories', icon: Tags },
  ]

  const isActive = (path: string) => location.pathname === path

  return (
    <div className="min-h-screen">
      {/* Mobile sidebar backdrop */}
      {sidebarOpen && (
        <div 
          className="fixed inset-0 z-40 bg-black/50 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <div className={`
        fixed inset-y-0 left-0 z-50 w-64 transform transition-transform duration-300 ease-in-out lg:translate-x-0
        ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
      `}>
        <div className="flex flex-col h-full glass-card rounded-none lg:rounded-r-2xl">
          {/* Logo */}
          <div className="flex items-center h-16 px-6 border-b border-primary-200/50">
            <div className="flex items-center gap-2">
              <div className="p-2 bg-gradient-primary rounded-lg">
                <Sparkles className="w-6 h-6 text-white" />
              </div>
              <span className="text-xl font-bold text-gradient-primary">理財管理</span>
            </div>
          </div>

          {/* Navigation */}
          <nav className="flex-1 px-6 py-6 space-y-2">
            {navigation.map((item) => {
              const Icon = item.icon
              return (
                <Link
                  key={item.name}
                  to={item.href}
                  className={`
                    flex items-center gap-3 px-4 py-3 rounded-lg transition-all duration-200 group
                    ${isActive(item.href) 
                      ? 'bg-gradient-primary text-white shadow-lg' 
                      : 'text-neutral-700 hover:bg-primary-100 hover:text-primary-700'
                    }
                  `}
                  onClick={() => setSidebarOpen(false)}
                >
                  <Icon className={`w-5 h-5 transition-colors ${
                    isActive(item.href) ? 'text-white' : 'text-neutral-500 group-hover:text-primary-600'
                  }`} />
                  <span className="font-medium">{item.name}</span>
                </Link>
              )
            })}
          </nav>

          {/* User info */}
          <div className="px-6 py-4 border-t border-primary-200/50">
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 bg-gradient-primary rounded-full flex items-center justify-center">
                <span className="text-sm font-medium text-white">U</span>
              </div>
              <div>
                <p className="text-sm font-medium text-neutral-800">使用者</p>
                <p className="text-xs text-neutral-500">user@example.com</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Main content */}
      <div className="lg:pl-64">
        {/* Top bar */}
        <div className="sticky top-0 z-30 glass-card rounded-none lg:rounded-b-2xl lg:mx-6 lg:mt-6">
          <div className="flex items-center justify-between h-16 px-6">
            <button
              className="lg:hidden text-neutral-600 hover:text-neutral-800"
              onClick={() => setSidebarOpen(true)}
            >
              <Menu className="w-6 h-6" />
            </button>
            
            {/* Breadcrumb or title can go here */}
            <div className="flex-1 lg:block hidden">
              <h1 className="text-2xl font-bold text-gradient-primary">
                {navigation.find(item => isActive(item.href))?.name || '理財管理系統'}
              </h1>
            </div>
          </div>
        </div>

        {/* Page content */}
        <main className="p-6">
          <div className="max-w-7xl mx-auto">
            {children}
          </div>
        </main>
      </div>

      {/* Mobile sidebar close button */}
      {sidebarOpen && (
        <button
          className="fixed top-4 right-4 z-50 p-2 bg-white rounded-full shadow-lg lg:hidden"
          onClick={() => setSidebarOpen(false)}
        >
          <X className="w-5 h-5" />
        </button>
      )}
    </div>
  )
}

export default Layout