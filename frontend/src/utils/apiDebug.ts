// Enhanced API debugging utility for tracking requests and responses

interface DebugLogEntry {
  timestamp: string
  url: string
  method: string
  fullUrl: string
  component: string
  success: boolean
  data?: any
  error?: any
  responseTime?: number
}

class ApiDebugger {
  private logs: DebugLogEntry[] = []
  private enabled: boolean = true

  log(entry: Omit<DebugLogEntry, 'timestamp'>) {
    if (!this.enabled) return
    
    const logEntry: DebugLogEntry = {
      ...entry,
      timestamp: new Date().toISOString()
    }
    
    this.logs.push(logEntry)
    
    // Console logging with colors
    const colorStyle = entry.success 
      ? 'color: green; font-weight: bold'
      : 'color: red; font-weight: bold'
    
    console.group(`%c[API Debug] ${entry.method} ${entry.url}`, colorStyle)
    console.log('🕒 Timestamp:', logEntry.timestamp)
    console.log('🌐 Full URL:', entry.fullUrl)
    console.log('📱 Component:', entry.component)
    console.log('⚡ Response Time:', entry.responseTime ? `${entry.responseTime}ms` : 'N/A')
    
    if (entry.success) {
      console.log('✅ Success Data:', entry.data)
    } else {
      console.log('❌ Error:', entry.error)
    }
    console.groupEnd()
    
    // Keep only last 50 logs
    if (this.logs.length > 50) {
      this.logs = this.logs.slice(-50)
    }
  }

  getLogs(): DebugLogEntry[] {
    return [...this.logs]
  }

  getLogsByComponent(component: string): DebugLogEntry[] {
    return this.logs.filter(log => log.component === component)
  }

  clearLogs() {
    this.logs = []
    console.log('🧹 API Debug logs cleared')
  }

  enable() {
    this.enabled = true
    console.log('🔍 API Debugging enabled')
  }

  disable() {
    this.enabled = false
    console.log('🚫 API Debugging disabled')
  }

  // Port detection utility
  detectPortIssues() {
    console.group('🔍 Port Detection Analysis')
    console.log('🎯 Expected: Frontend should use port 3000 (with proxy to 8080)')
    console.log('🏠 Current location:', window.location.href)
    console.log('🌐 Current origin:', window.location.origin)
    
    const wrongPortLogs = this.logs.filter(log => 
      log.fullUrl.includes(':3001') || log.fullUrl.includes(':8080')
    )
    
    if (wrongPortLogs.length > 0) {
      console.log('⚠️ Found requests to wrong ports:')
      wrongPortLogs.forEach(log => {
        console.log(`  - ${log.method} ${log.fullUrl} (from ${log.component})`)
      })
    } else {
      console.log('✅ No wrong port requests detected')
    }
    
    console.groupEnd()
  }

  // Component comparison utility
  compareComponents(component1: string, component2: string) {
    const logs1 = this.getLogsByComponent(component1)
    const logs2 = this.getLogsByComponent(component2)
    
    console.group(`🔍 Comparing ${component1} vs ${component2}`)
    console.log(`📊 ${component1} API calls:`, logs1.length)
    console.log(`📊 ${component2} API calls:`, logs2.length)
    
    // Compare endpoints
    const endpoints1 = new Set(logs1.map(log => log.url))
    const endpoints2 = new Set(logs2.map(log => log.url))
    
    console.log(`🎯 ${component1} endpoints:`, Array.from(endpoints1))
    console.log(`🎯 ${component2} endpoints:`, Array.from(endpoints2))
    
    // Find differences
    const onlyIn1 = Array.from(endpoints1).filter(ep => !endpoints2.has(ep))
    const onlyIn2 = Array.from(endpoints2).filter(ep => !endpoints1.has(ep))
    
    if (onlyIn1.length > 0) {
      console.log(`⚠️ Only in ${component1}:`, onlyIn1)
    }
    if (onlyIn2.length > 0) {
      console.log(`⚠️ Only in ${component2}:`, onlyIn2)
    }
    
    console.groupEnd()
  }
}

// Global instance
export const apiDebugger = new ApiDebugger()

// Make it available globally for browser console debugging
if (typeof window !== 'undefined') {
  (window as any).apiDebugger = apiDebugger
}

export default apiDebugger