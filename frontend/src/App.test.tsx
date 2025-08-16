// Simple test component without React Query
import React from 'react'

function SimpleApp() {
  return (
    <div style={{ padding: '20px' }}>
      <h1>測試頁面</h1>
      <p>如果你看到這個，說明 React 基本功能正常</p>
      <button onClick={() => alert('按鈕有效!')}>測試按鈕</button>
    </div>
  )
}

export default SimpleApp