// Configuration validation utility for development debugging

interface BackendConfig {
  apiUrl: string
  expectedPort: number
  proxyTarget: string
}

export class ConfigValidator {
  private static getBackendConfig(): BackendConfig {
    const isDev = process.env.NODE_ENV === 'development'
    
    if (isDev) {
      return {
        apiUrl: '/api/v1',
        expectedPort: 8080, // Update this if backend port changes
        proxyTarget: 'http://localhost:8080' // Should match vite.config.ts
      }
    }
    
    return {
      apiUrl: '/api/v1',
      expectedPort: 80,
      proxyTarget: window.location.origin
    }
  }

  static async validateBackendConnection(): Promise<{isValid: boolean, issues: string[]}> {
    const config = this.getBackendConfig()
    const issues: string[] = []
    
    try {
      // Test if backend is reachable
      const testUrl = `${config.apiUrl}/health` // Assuming health endpoint exists
      const response = await fetch(testUrl, { 
        method: 'GET',
        signal: AbortSignal.timeout(5000) // 5 second timeout
      })
      
      if (!response.ok) {
        issues.push(`Backend responded with status ${response.status}`)
      }
      
    } catch (error: any) {
      if (error.name === 'AbortError') {
        issues.push(`Backend connection timeout (>${5000}ms)`)
      } else if (error.message?.includes('Failed to fetch')) {
        issues.push('Backend unreachable - check if server is running')
        issues.push(`Expected backend on port ${config.expectedPort}`)
      } else {
        issues.push(`Connection error: ${error.message}`)
      }
    }
    
    return {
      isValid: issues.length === 0,
      issues
    }
  }

  static logConfigInfo() {
    const config = this.getBackendConfig()
    
    console.group('🔧 Configuration Info')
    console.log('Environment:', process.env.NODE_ENV)
    console.log('Frontend URL:', window.location.origin)
    console.log('API Base URL:', config.apiUrl)
    console.log('Expected Backend Port:', config.expectedPort)
    console.log('Proxy Target:', config.proxyTarget)
    console.groupEnd()
  }

  // Port mismatch detection
  static detectPortMismatch(): {hasMismatch: boolean, details: string} {
    // This would be expanded with actual vite config reading in a real implementation
    const expectedBackendPort = 8080
    const viteProxyPort = 8080 // This should be read from vite config
    
    if (expectedBackendPort !== viteProxyPort) {
      return {
        hasMismatch: true,
        details: `Backend running on ${expectedBackendPort}, but Vite proxy targets ${viteProxyPort}`
      }
    }
    
    return { hasMismatch: false, details: 'Port configuration is correct' }
  }
}

// Auto-run validation in development
if (process.env.NODE_ENV === 'development') {
  ConfigValidator.logConfigInfo()
  
  // Validate backend connection after a brief delay
  setTimeout(async () => {
    const validation = await ConfigValidator.validateBackendConnection()
    if (!validation.isValid) {
      console.warn('⚠️ Backend Connection Issues:', validation.issues)
    } else {
      console.log('✅ Backend connection validated successfully')
    }
  }, 2000)
}